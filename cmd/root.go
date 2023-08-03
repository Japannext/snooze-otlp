package cmd

import (
  "fmt"
  "os"

	"github.com/spf13/cobra"

  "github.com/japannext/snooze-otlp/server"
)

var rootCmd = &cobra.Command{
  Use: "snooze-otlp",
  Short: "",
  Long: ``,
}

func Execute() {
  if err := rootCmd.Execute(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}

func init() {
  cobra.OnInitialize(server.InitServer)

  rootCmd.AddCommand(runCmd)
  rootCmd.AddCommand(versionCmd)
}
