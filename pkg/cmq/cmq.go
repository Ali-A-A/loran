package cmq

import (
	"time"

	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

const (
	wait = 10
)

// NatsOptions represents nats connection options
type NatsOptions struct {
	URL             string
	Subject string
}

// NewNatsConn creates new nats connection
func NewNatsConn(natsConfig NatsOptions) *nats.Conn {
	// Connect Options
	opts := []nats.Option{nats.Name("NATS Subscriber")}
	opts = setupConnOptions(opts)

	// Connect to NATS
	nc, err := nats.Connect(natsConfig.URL, opts...)
	if err != nil {
		logrus.Fatal(err)
	}

	return nc
}

func setupConnOptions(opts []nats.Option) []nats.Option {
	totalWait := wait * time.Minute
	reconnectDelay := time.Second

	opts = append(opts, nats.ReconnectWait(reconnectDelay))
	opts = append(opts, nats.MaxReconnects(int(totalWait/reconnectDelay)))
	opts = append(opts, nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
		logrus.Printf("Disconnected due to:%s, will attempt reconnects for %.0fm", err, totalWait.Minutes())
	}))
	opts = append(opts, nats.ErrorHandler(func(_ *nats.Conn, _ *nats.Subscription, err error) {
		logrus.Errorf("nats error handler: %s", err)
	}))
	opts = append(opts, nats.ReconnectHandler(func(nc *nats.Conn) {
		logrus.Printf("Reconnected [%s]", nc.ConnectedUrl())
	}))
	opts = append(opts, nats.ClosedHandler(func(nc *nats.Conn) {
		logrus.Fatalf("Exiting: %v", nc.LastError())
	}))

	return opts
}
