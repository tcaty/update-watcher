package cmd

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

type Flags struct {
	CfgFile string
}

var (
	rootFlags = &Flags{}
	rootCmd   = &cobra.Command{
		Use: "watch",
	}
)

func Execute() *Flags {
	if err := rootCmd.Execute(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	return rootFlags
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&rootFlags.CfgFile, "config", "c", "", "config file path (default is $HOME/config.yaml)")
}
