package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.com/littledot/mockhiato/lib/clean"
)

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Delete mocks",
	Long:  `Delete mocks.`,
	Run: func(cmd *cobra.Command, args []string) {
		config := getConfig(cmd)
		clean.Run(config)
	},
}

func init() {
	RootCmd.AddCommand(cleanCmd)
}
