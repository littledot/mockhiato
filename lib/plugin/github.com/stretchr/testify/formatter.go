package testify

import (
	"bytes"
	"fmt"
	"go/types"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"gitlab.com/littledot/mockhiato/lib"
)

// NewTestifyFormatter creates a new testifyFormatter
func NewTestifyFormatter() lib.Formatter {
	return &testifyFormatter{}
}

type testifyFormatter struct {
}

func (r *testifyFormatter) GenerateMocks(spec *lib.Spec) {
	for _, pack := range spec.Packages {
		r.generateMock(pack)
	}
}

func (r *testifyFormatter) generateMock(pack *lib.Package) {
	if len(pack.Interfaces) == 0 { // Nothing to mock? Return early
		return
	}

	// Create file {package_name}_mocks.go
	if err := os.MkdirAll("mocks", os.ModeDir); err != nil {
		panic(err)
	}
	pkgName := pack.Interfaces[0].Object.Pkg().Name()
	mockPath := fmt.Sprintf("mocks/%s_mocks.go", strings.Replace(pkgName, "/", "_", -1))
	mockFile, err := os.Create(mockPath)
	if err != nil {
		panic(err)
	}
	defer mockFile.Close()

	buf := &bytes.Buffer{}

	// Write package
	buf.WriteString("package mocks\n")

	// Write magic string
	buf.WriteString("// Auto-generated by mockhiato\n")

	// Write imports
	imports := []string{`"github.com/stretchr/testify/mock"`}
	for _, imported := range pack.Imports {
		imports = append(imports, strconv.Quote(imported.Path()))
	}
	buf.WriteString("import(\n")
	buf.WriteString(strings.Join(imports, "\n"))
	buf.WriteString("\n)\n")

	for _, iface := range pack.Interfaces {
		// Write struct
		interfaceName := iface.Object.Name()
		buf.WriteString(fmt.Sprintf("// %s implements %s.%s\n", interfaceName, pkgName, interfaceName))
		buf.WriteString(fmt.Sprintf("type %s struct { mock.Mock }\n", interfaceName))

		for i := 0; i < iface.Interface.NumMethods(); i++ {
			// Write method
			method := iface.Interface.Method(i)
			signature := method.Type().(*types.Signature)

			paramNames := []string{}
			paramExprs := []string{}
			for j := 0; j < signature.Params().Len(); j++ {
				param := signature.Params().At(j)
				varName := "p" + strconv.Itoa(j)
				paramNames = append(paramNames, varName)
				paramExprs = append(paramExprs, varName+" "+lib.GetTypeString(param))
			}

			if signature.Variadic() { // Variadic method? Replace last parameter's [] with ... ("p1 []int" -> "p1 ...int")
				pos := signature.Params().Len() - 1
				paramExprs[pos] = strings.Replace(paramExprs[pos], "[]", "...", -1)
			}

			returnNames := []string{}
			returnTypes := []string{}
			verifyReturnLines := []string{}
			for j := 0; j < signature.Results().Len(); j++ {
				result := signature.Results().At(j)
				varName := "ret" + strconv.Itoa(j)
				returnNames = append(returnNames, varName)
				returnTypes = append(returnTypes, lib.GetTypeString(result))
				verifyReturnLines = append(verifyReturnLines, fmt.Sprintf("%s := ret.Get(%d).(%s)\n", varName, j, lib.GetTypeString(result)))
			}

			commentLine := fmt.Sprintf("// %s implements (%s.%s).%s\n", method.Name(), pkgName, interfaceName, method.Name())
			signatureLine := fmt.Sprintf("func (r *%s) %s(%s) (%s) {\n", interfaceName, method.Name(), strings.Join(paramExprs, ", "), strings.Join(returnTypes, ","))
			verifyInvokedLine := fmt.Sprintf("r.Called(%s)\n", strings.Join(paramNames, ", "))
			if signature.Results().Len() > 0 {
				verifyInvokedLine = "ret := " + verifyInvokedLine
			}
			returnLine := fmt.Sprintf("return %s\n}\n", strings.Join(returnNames, ", "))

			buf.WriteString(commentLine)
			buf.WriteString(signatureLine)
			buf.WriteString(verifyInvokedLine)
			for _, verifyReturnLine := range verifyReturnLines {
				buf.WriteString(verifyReturnLine)
			}
			buf.WriteString(returnLine)
		}

	}

	mockFile.Write(buf.Bytes())
	mockFile.Close()

	// Format generated code
	cmd := exec.Command("goimports", "-w", mockPath)
	if stdout, e := cmd.CombinedOutput(); e != nil {
		panic(stdout)
	}
}