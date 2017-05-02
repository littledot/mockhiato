package test_test

import (
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerate(t *testing.T) {
	a := assert.New(t)

	files, err := ioutil.ReadDir(".")
	if err != nil {
		a.FailNow(err.Error())
	}

	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		t.Run(file.Name(), func(t *testing.T) {
			actualPath := filepath.Join(file.Name(), "actual")
			cmd := exec.Command("mockhiato", "generate", "-p="+actualPath)
			if stdout, err := cmd.CombinedOutput(); err != nil {
				a.Error(err)
				a.Fail(string(stdout))
			}

			expectPath := filepath.Join(file.Name(), "expect")
			cmd = exec.Command("diff", "-r", expectPath, actualPath)
			if stdout, err := cmd.CombinedOutput(); err != nil {
				a.Error(err)
				a.Fail(string(stdout))
			}
		})
	}
}
