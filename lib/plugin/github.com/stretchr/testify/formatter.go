package testify

import (
	"bufio"
	"bytes"
	"fmt"
	"go/format"
	"go/types"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"gitlab.com/littledot/mockhiato/lib"
)

const (
	magic = "// Auto-generated by Mockhiato"
)

// NewTestifyFormatter creates a new testifyFormatter
func NewTestifyFormatter(config lib.Config) lib.Formatter {
	return &testifyFormatter{
		config: config,
	}
}

type testifyFormatter struct {
	config lib.Config
}

func (r *testifyFormatter) IsMockFile(file *os.File) bool {
	fileScanner := bufio.NewScanner(file)
	for i := 0; i < 6 && fileScanner.Scan(); i++ { // Check first 5 lines or EOF
		line := fileScanner.Text()
		if strings.HasPrefix(line, "//") { // Check 1st comment, which should be the 3rd line
			if line == magic { // Found magic string? File is a mock
				return true
			}
			break
		}
	}
	return false
}

func (r *testifyFormatter) GenerateMocks(project *lib.Project) {
	for _, pack := range project.Packages {
		r.generateMock(project, pack)
	}
}

func (r *testifyFormatter) generateMock(project *lib.Project, pack *lib.Package) {
	if len(pack.Interfaces) == 0 { // Nothing to mock? Return early
		return
	}
	log.Debugf("Generating mocks for %s", pack.Context.Path())
	mockPath := filepath.Join(project.GoSrcAbsPath, pack.Context.Path(), r.config.MockFileName)
	log.Debugf("Creating file: %s", mockPath)
	mockFile, err := os.Create(mockPath)
	if err != nil {
		panic(err)
	}
	defer mockFile.Close()

	buf := &bytes.Buffer{}
	pf := lib.NewPackageFormatter(pack.Context)
	pf.RecordDependency(types.NewPackage("github.com/stretchr/testify/mock", "mock"))
	pf.IndexImports(pack.Context)

	// Write package
	packageName := pack.Context.Name()
	log.Debugf("Writing package: %s", packageName)
	buf.WriteString(fmt.Sprintf("package %s\n", packageName))

	// Write magic string
	log.Debugf("Writing magic: %s", magic)
	buf.WriteString(magic + "\n")

	// Write imports
	imports := []string{}
	for depPath, depAlias := range pf.PathToAlias {
		imports = append(imports, fmt.Sprintf(`%s "%s"`, depAlias, strings.TrimPrefix(depPath, project.VendorPath+"/")))
	}
	log.Debugf("Writing imports: %s", strings.Join(imports, ", "))
	buf.WriteString("import(\n")
	buf.WriteString(strings.Join(imports, "\n"))
	buf.WriteString("\n)\n")

	for _, iface := range pack.Interfaces {
		// Write struct
		interfaceName := iface.TObject.Name()
		mockName := interfaceName + "Mock"
		log.Debugf("Writing struct: %s", mockName)
		buf.WriteString(fmt.Sprintf("// %s implements %s.%s\n", mockName, packageName, interfaceName))
		buf.WriteString(fmt.Sprintf("type %s struct { mock.Mock }\n", mockName))

		for i := 0; i < iface.TInterface.NumMethods(); i++ {
			// Write method
			method := iface.TInterface.Method(i)
			signature := method.Type().(*types.Signature)

			paramNames := []string{}
			paramExprs := []string{}
			for j := 0; j < signature.Params().Len(); j++ {
				param := signature.Params().At(j)
				varName := "p" + strconv.Itoa(j)
				paramNames = append(paramNames, varName)
				paramExprs = append(paramExprs, varName+" "+pf.ObjectTypeString(param))
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
				returnTypes = append(returnTypes, pf.ObjectTypeString(result))
				verifyReturnLines = append(verifyReturnLines, fmt.Sprintf("%s := ret.Get(%d).(%s)\n", varName, j, pf.ObjectTypeString(result)))
			}

			log.Debugf("Writing method: %s()", method.Name())
			commentLine := fmt.Sprintf("// %s implements (%s.%s).%s\n", method.Name(), packageName, interfaceName, method.Name())
			signatureLine := fmt.Sprintf("func (r *%s) %s(%s) (%s) {\n", mockName, method.Name(), strings.Join(paramExprs, ", "), strings.Join(returnTypes, ","))
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

	// Format generated code
	sourceCode, err := format.Source(buf.Bytes())
	if err != nil {
		panic(err)
	}

	mockFile.Write(sourceCode)
	mockFile.Close()

	project.GenAbsPaths = append(project.GenAbsPaths, mockPath)
}
