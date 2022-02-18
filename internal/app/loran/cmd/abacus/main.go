package abacus

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/ali-a-a/loran/config"
	"github.com/ali-a-a/loran/internal/app/loran/abacus"
	"github.com/ali-a-a/loran/internal/app/loran/abacus/model"
	"github.com/ali-a-a/loran/pkg/cmq"
	"github.com/ali-a-a/loran/pkg/redis"
	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

func main(cfg config.Config) {
	conn, err := cmq.CreateJetStreamConnection(cfg.NATS)
	if err != nil {
		logrus.Fatalf("failed to create nats connection: %s", err.Error())
	}

	defer func() {
		conn.NC.Close()
	}()

	rc, err := redis.NewRedisClient(cfg.Redis)
	if err != nil {
		logrus.Fatalf("failed to create redis connection: %s", err.Error())
	}

	defer func() {
		if err := rc.Close(); err != nil {
			logrus.Errorf("redis close error: %s", err.Error())
		}
	}()

	cr := model.NewInMemoryCalculator(rc)

	handler, err := abacus.NewHandler(cr, conn, cfg.NATS)
	if err != nil {
		logrus.Fatalf("failed to create new abacus handler: %s", err.Error())
	}

	go func() {
		if err = handler.Run(); err != nil {
			logrus.Fatalf("failed to create run abacus handler: %s", err.Error())
		}
	}()

	logrus.Info("abacus is ready!")

	signals := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-signals
		logrus.Info("signal received: ", sig)
		done <- true
	}()

	<-done
}

// Register registers abacus command for loran binary.
// Abacus is consumer module.
func Register(root *cobra.Command, cfg config.Config) {
	root.AddCommand(
		&cobra.Command{
			Use:   "abacus",
			Short: "Run loran abacus component",
			Run: func(cmd *cobra.Command, args []string) {
				main(cfg)
			},
		},
	)
}
