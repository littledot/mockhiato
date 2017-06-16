package generate_test

import (
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/littledot/mockhiato/lib"
	"github.com/littledot/mockhiato/lib/generate"
	"github.com/stretchr/testify/assert"
)

func TestGenerate(t *testing.T) {
	a := assert.New(t)

	testPath := "test"
	files, err := ioutil.ReadDir(testPath)
	if err != nil {
		a.FailNow(err.Error())
	}

	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		t.Run(file.Name(), func(t *testing.T) {
			actualPath := filepath.Join(testPath, file.Name(), "actual")

			config := lib.Config{}
			config.ProjectPath = actualPath
			config.MockFileName = "mockhiato_mocks.go"
			config.DependentMocksPath = "mocks"
			config.StructNameFormat = "{interface}Mock"

			generate.Run(config)

			expectPath := filepath.Join(testPath, file.Name(), "expect")
			cmd := exec.Command("diff", "-r", expectPath, actualPath)
			if stdout, err := cmd.CombinedOutput(); err != nil {
				a.Error(err)
				a.Fail(string(stdout))
			}
		})
	}
}
