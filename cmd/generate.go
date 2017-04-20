// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
