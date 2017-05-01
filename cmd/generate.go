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

	generateCmd.Flags().StringP("ProjectPath", "p", ".",
		"Configures path of the project to generate mocks for. Default is current working directory.")
	generateCmd.Flags().StringP("MockFileName", "n", "mockhiato_mocks.go",
		"Configures the file name of generated mocks.")
}
