package generate_test

import (
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/blang/semver"
	"github.com/littledot/mockhiato/lib"
	"github.com/littledot/mockhiato/lib/generate"
	"github.com/stretchr/testify/assert"
)

func TestGenerate(t *testing.T) {
	a := assert.New(t)

	testsPath := "testdata"
	files, err := ioutil.ReadDir(testsPath)
	if err != nil {
		a.FailNow(err.Error())
	}

	for _, verFile := range files {
		if !verFile.IsDir() {
			continue
		}
		minGoVer, err := semver.ParseTolerant(verFile.Name())
		a.NoError(err)

		goVer, err := semver.ParseTolerant(strings.TrimPrefix(runtime.Version(), "go"))
		a.NoError(err)

		if goVer.GTE(minGoVer) {
			verPath := filepath.Join(testsPath, verFile.Name())
			testFiles, err := ioutil.ReadDir(verPath)
			a.NoError(err)

			for _, testFile := range testFiles {
				if !testFile.IsDir() {
					continue
				}

				t.Run(testFile.Name(), func(t *testing.T) {
					actualPath := filepath.Join(verPath, testFile.Name(), "actual")

					config := lib.Config{}
					config.ProjectPath = actualPath
					config.MockFileName = "mockhiato_mocks.go"
					config.DependentMocksPath = "mocks"
					config.StructNameFormat = "{interface}Mock"
					config.DependentPackageNameFormat = "m{package}"
					config.MockFunctionMethodNameFormat = "{func}Run"

					generate.Run(config)

					expectPath := filepath.Join(verPath, testFile.Name(), "expect")
					cmd := exec.Command("diff", "-r", expectPath, actualPath)
					if stdout, err := cmd.CombinedOutput(); err != nil {
						a.Error(err)
						a.Fail(string(stdout))
					}
				})
			}
		}
	}

}
