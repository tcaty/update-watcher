package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tcaty/update-watcher/pkg/utils"
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
		utils.HandleFatal("rotCmd execute occured error", err)
	}
	return rootFlags
}

func init() {
	rootCmd.PersistentFlags().StringVarP(
		&rootFlags.CfgFile, "config", "c", "", "config file path (default is $HOME/config.yaml)",
	)
}
