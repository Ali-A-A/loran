package cmd

import (
	"github.com/ali-a-a/loran/config"
	"github.com/ali-a-a/loran/internal/app/loran/cmd/abacus"
	"github.com/ali-a-a/loran/internal/app/loran/cmd/cranmer"
	"github.com/ali-a-a/loran/internal/app/loran/cmd/scheduler"
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

	abacus.Register(root, cfg)
	cranmer.Register(root, cfg)
	scheduler.Register(root, cfg)

	return root
}
