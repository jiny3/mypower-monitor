package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd cobra.Command = cobra.Command{
	Use:   "ucasnj-smi",
	Short: "ucasnj-smi 是一个宿舍电量监控工具",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() error {
	rootCmd.AddCommand(checkCmd())
	rootCmd.AddCommand(serverCmd())
	return rootCmd.Execute()
}
