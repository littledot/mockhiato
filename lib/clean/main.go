package clean

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/littledot/mockhiato/lib"
	"github.com/littledot/mockhiato/lib/plugin/github.com/stretchr/testify"
)

// Run executes the command.
func Run(config lib.Config) {
	projectPath, err := filepath.Abs(config.ProjectPath)
	if err != nil {
		panic(err)
	}

	formatter := testify.NewTestifyFormatter(config)
	err = filepath.Walk(projectPath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil { // Something wrong? Skip
			return nil
		}
		if !strings.HasSuffix(filePath, ".go") { // Not Go source? Skip
			return nil
		}

		file, err := os.Open(filePath)
		if err != nil {
			return nil
		}
		defer file.Close()
		if formatter.IsMockFile(file) { // Formatter says its a mock? Remove
			fmt.Println("Removing", filePath)
			if err := os.Remove(filePath); err != nil {
				panic(err)
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}
