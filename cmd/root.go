package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type Flags struct {
	CfgFile string
}

var (
	rootFlags = &Flags{}
	rootCmd   = &cobra.Command{
		Use:   "watch",
		Short: "Update watcher is tool for automatic updates check",
		Long:  "[Long] Update watcher is tool for automatic updates check",
	}
)

func Execute() *Flags {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return rootFlags
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&rootFlags.CfgFile, "config", "c", "", "config file path (default is $HOME/config.yaml)")
}
