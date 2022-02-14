package consumer

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Ali-A-A/loran/config"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// nolint:funlen
func main(cfg config.Config) {
	signals := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-signals
		logrus.Info("Signal received: ", sig)
		done <- true
	}()

	<-done
}

// Register registers server command for loran binary.
func Register(root *cobra.Command, cfg config.Config) {
	root.AddCommand(
		&cobra.Command{
			Use:   "server",
			Short: "Run Loran consumer component",
			Run: func(cmd *cobra.Command, args []string) {
				main(cfg)
			},
		},
	)
}
