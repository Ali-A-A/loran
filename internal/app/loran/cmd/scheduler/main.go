package scheduler

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/ali-a-a/loran/config"
	"github.com/ali-a-a/loran/pkg/cmq"
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

	logrus.Info("scheduler is ready!")

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

// Register registers scheduler command for loran binary.
func Register(root *cobra.Command, cfg config.Config) {
	root.AddCommand(
		&cobra.Command{
			Use:   "scheduler",
			Short: "Run loran scheduler component",
			Run: func(cmd *cobra.Command, args []string) {
				main(cfg)
			},
		},
	)
}
