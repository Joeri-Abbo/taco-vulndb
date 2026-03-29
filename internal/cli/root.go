package cli

import (
	"log/slog"

	"github.com/spf13/cobra"
)

var (
	debug bool
	quiet bool
)

func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "taco-vulndb",
		Short: "TACO VulnDB — vulnerability database builder and distributor",
		Long: `Build, update, and distribute the TACO vulnerability database.

This tool fetches vulnerability data from multiple sources (NVD, OSV, GHSA,
Alpine, Debian, Ubuntu, Red Hat, ALAS, CISA KEV), merges them with precedence
rules, and distributes the result via OCI registries, HTTP, or file export.`,
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if debug {
				slog.SetLogLoggerLevel(slog.LevelDebug)
			}
		},
	}

	root.PersistentFlags().BoolVar(&debug, "debug", false, "enable debug logging")
	root.PersistentFlags().BoolVar(&quiet, "quiet", false, "suppress non-essential output")

	root.AddCommand(newUpdateCmd())
	root.AddCommand(newDownloadCmd())
	root.AddCommand(newLoadCmd())
	root.AddCommand(newBuildCmd())
	root.AddCommand(newExportCmd())
	root.AddCommand(newServeCmd())
	root.AddCommand(newStatusCmd())
	root.AddCommand(newPushCmd())
	root.AddCommand(newPullCmd())

	return root
}
