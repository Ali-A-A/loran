package cranmer

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ali-a-a/loran/config"
	"github.com/ali-a-a/loran/internal/app/loran/cranmer"
	"github.com/ali-a-a/loran/pkg/cmq"
	"github.com/ali-a-a/loran/pkg/router"
	"github.com/labstack/echo/v4"
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

	handler := cranmer.NewHandler(conn, cfg.NATS)

	e := router.New()

	e.GET("/ready", func(c echo.Context) error { return c.NoContent(http.StatusNoContent) })

	api := e.Group("/api")

	api.POST("/add", handler.Add)

	go func() {
		if err := e.Start(fmt.Sprintf(":%d", cfg.Cranmer.Port)); !errors.Is(err, http.ErrServerClosed) && err != nil {
			e.Logger.Fatal(err.Error())
		}
	}()

	logrus.Info("cranmer is ready!")

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

// Register registers cranmer command for loran binary.
// Cranmer is producer module.
func Register(root *cobra.Command, cfg config.Config) {
	root.AddCommand(
		&cobra.Command{
			Use:   "cranmer",
			Short: "Run loran cranmer component",
			Run: func(cmd *cobra.Command, args []string) {
				main(cfg)
			},
		},
	)
}
