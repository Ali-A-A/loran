package cmd

import (
	"github.com/Ali-A-A/loran/config"
	"github.com/Ali-A-A/loran/internal/app/loran/cmd/consumer"
	"github.com/Ali-A-A/loran/pkg/log"

	"github.com/spf13/cobra"
)

// NewRootCommand creates a new loran root command.
func NewRootCommand() *cobra.Command {
	var root = &cobra.Command{
		Use: "loran",
	}

	cfg := config.Init()

	log.SetupLogger(log.AppLogger{
		Level:  cfg.Logger.Level,
		StdOut: true,
	})

	consumer.Register(root, cfg)

	return root
}
