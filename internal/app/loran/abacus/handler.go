package abacus

import (
	"fmt"

	"github.com/ali-a-a/loran/config"
	"github.com/ali-a-a/loran/pkg/cmq"
	"github.com/go-redis/redis/v8"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"

	"github.com/panjf2000/ants/v2"
)

const (
	// workerSize represents the size of the pool.
	// In this module, we have a worker pool that each
	// of them is responsible to subscribe on event subject
	// and consume messages for processing and calculating
	// distinct counts.
	workerSize = 20
)

// Handler represents abacus handler which is responsible
// to calculate distinct counts.
type Handler struct {
	rc      *redis.Client
	nc      *cmq.Conn
	subject string
	durable string

	workerPool *ants.Pool
}

// NewHandler creates new handler with *redis.Client and *cmq.Conn fields.
func NewHandler(rc *redis.Client, conn *cmq.Conn, cfg config.NATS) (*Handler, error) {
	// In this project, we use panjf2000/ants package for creating modules worker pool.
	workerPool, err := ants.NewPool(workerSize)
	if err != nil {
		return nil, err
	}

	return &Handler{
		rc:         rc,
		nc:         conn,
		subject:    cfg.JetStream.Consumer.Subject,
		durable:    cfg.JetStream.Consumer.Durable,
		workerPool: workerPool,
	}, nil
}

// Run starts pulling messages from jet stream server.
// First, it subscribe on subject.
// Then, it starts to fetch messages from server and submits them into the pool.
// It get config.NATS to figure out the subject and durable of the nats.
func (h *Handler) Run() error {
	sub, err := h.nc.JS.PullSubscribe(h.subject, h.durable)
	if err != nil {
		return fmt.Errorf("failed to run handler: %w", err)
	}

	h.fetch(sub)

	return nil
}

// fetch starts pulling messages based on subscription.
func (h *Handler) fetch(sub *nats.Subscription) {
	for {
		messages, err := sub.Fetch(workerSize)
		if err != nil {
			continue
		}

		for _, message := range messages {
			if message != nil {
				if err = h.workerPool.Submit(h.newTask(message)); err != nil {
					logrus.Errorf("failed to submit new task: %s", err.Error())
				}
			}
		}
	}
}

// newTask creates new pool task.
func (h *Handler) newTask(message *nats.Msg) func() {
	return func() {
		err := message.Ack()
		if err != nil {
			logrus.Errorf("failed to ack message: %s", err.Error())
		}

		logrus.Infof("new message received: %s", string(message.Data))
	}
}
