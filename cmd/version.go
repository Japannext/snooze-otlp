package cmd

import (
	"github.com/spf13/cobra"

	"github.com/japannext/snooze-otlp/server"
)

var versionCmd = &cobra.Command{
  Use: "version",
  Short: "Display snooze-otlp version",
  Long: "Display snooze-otlp version",
  Run: func(cmd *cobra.Command, args []string) {
    server.PrintVersion()
  },
}
