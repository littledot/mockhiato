package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.com/littledot/mockhiato/lib/generate"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate mocks",
	Long:  `Generate mocks.`,
	Run: func(cmd *cobra.Command, args []string) {
		config := getConfig(cmd)
		generate.Run(config)
	},
}

func init() {
	RootCmd.AddCommand(generateCmd)

	generateCmd.Flags().StringP("ProjectPath", "p", ".", "Project root path")
	generateCmd.Flags().StringSliceP("IgnorePaths", "i", []string{"vendor"}, "Ignore paths")
	generateCmd.Flags().StringP("OutputPath", "o", "mocks", "Output path")
}
