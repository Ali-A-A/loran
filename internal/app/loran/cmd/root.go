package cmd

import (
	"github.com/ali-a-a/loran/config"
	"github.com/ali-a-a/loran/internal/app/loran/cmd/consumer"
	"github.com/ali-a-a/loran/pkg/log"

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
