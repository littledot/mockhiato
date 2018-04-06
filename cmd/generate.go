package cmd

import (
	"github.com/littledot/mockhiato/lib/generate"
	"github.com/spf13/cobra"
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
	generateCmd.Flags().StringP("DependentMocksPath", "d", "mocks",
		"Configures where mocks for dependent interfaces (referenced but not defined by the project) will be created.")
	generateCmd.Flags().StringP("StructNameFormat", "s", "{interface}Mock",
		"Configures the name format for generated structs. Use {interface} as a placeholder for the name of the interface being mocked. For example, '{interface}Mock' means suffixing generated structs with 'Mock' (XMock, YMock, PipeReaderMock).")
	generateCmd.Flags().StringP("DependentPackageNameFormat", "f", "m{package}",
		"Configures the name format for generated dependent package names. Use {package} as a placeholder for the name of the package being generated. For example, 'm{package}' means prefixing generated package names with 'm' (mio, mbytes, mhttp).")
}
