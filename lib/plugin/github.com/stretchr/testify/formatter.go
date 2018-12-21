package testify

import (
	"bufio"
	"fmt"
	"go/format"
	"go/types"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/littledot/mockhiato/lib"
	log "github.com/sirupsen/logrus"
)

// NewTestifyFormatter creates a new testifyFormatter
func NewTestifyFormatter(config lib.Config) lib.Formatter {
	return &testifyFormatter{config}
}

type testifyFormatter struct {
	config lib.Config
}

func (r *testifyFormatter) IsMockFile(file *os.File) bool {
	fileScanner := bufio.NewScanner(file)
	for i := 0; i < 6 && fileScanner.Scan(); i++ { // Check first 5 lines or EOF
		line := fileScanner.Text()
		if strings.HasPrefix(line, "//") { // Check 1st comment, which should be the 3rd line
			if line == lib.Magic { // Found magic string? File is a mock
				return true
			}
			break
		}
	}
	return false
}

func (r *testifyFormatter) GenerateMocks(project *lib.Project) {
	for _, pack := range project.Packages {
		r.generateMocksForPackage(project, pack)
	}
}

func (r *testifyFormatter) generateMocksForPackage(project *lib.Project, pack *lib.Package) {
	mockPkg := r.mockPackage(project, pack.TPackage)
	r.generateMock(project, mockPkg, pack)
}

func (r *testifyFormatter) mockPackage(project *lib.Project, srcPkg *types.Package) *types.Package {
	// Internal package? Put mocks directly in original package
	pkgName := srcPkg.Name()
	pkgPath := srcPkg.Path()
	if lib.IsExternalDependency(project, pkgPath) { // External package? Put mocks under dependent mocks path
		pkgName = strings.Replace(r.config.DependentPackageNameFormat, lib.PackageNameToken, pkgName, -1)
		pkgPath = strings.TrimPrefix(pkgPath, project.VendorPath)
		pkgPath = filepath.Join(project.DependentMocksPath, pkgPath, pkgName)
	}
	return types.NewPackage(pkgPath, pkgName)
}

func (r *testifyFormatter) generateMock(project *lib.Project, pkg *types.Package, pack *lib.Package) {
	if len(pack.Interfaces) == 0 && len(pack.Signatures) == 0 {
		return
	}
	log.Debugf("Generating mocks for %s", pkg.Path())

	buf := &lib.Bufferw{}
	pm := lib.NewPackageMapper(pkg.Path())
	pm.RecordDependency(types.NewPackage("github.com/stretchr/testify/mock", "mock"))
	for _, iface := range pack.Interfaces {
		pm.IndexInterface(iface.TInterface)
	}
	for _, signature := range pack.Signatures {
		pm.IndexSignature(signature.TSignature)
	}

	// Write package
	pkgName := pkg.Name()
	log.Debugf("Writing package: %s", pkgName)
	buf.WriteString(fmt.Sprintf("package %s\n", pkgName))

	// Write magic string
	log.Debugf("Writing magic: %s", lib.Magic)
	buf.WriteString(lib.Magic + "\n")

	// Write imports
	imports := []string{}
	for depPath, depAlias := range pm.PathToAlias {
		imports = append(imports, fmt.Sprintf(`%s "%s"`, depAlias, strings.TrimPrefix(depPath, project.VendorPath+"/")))
	}
	log.Debugf("Writing imports: %s", strings.Join(imports, ", "))
	buf.WriteString("import(\n")
	buf.WriteString(strings.Join(imports, "\n"))
	buf.WriteString("\n)\n")

	// Write interface mocks
	for _, iface := range pack.Interfaces {
		writeStruct(buf, r.config, iface.TObject, pkgName)
		for i := 0; i < iface.TInterface.NumMethods(); i++ {
			method := iface.TInterface.Method(i)
			signature := method.Type().(*types.Signature)
			writeMethod(buf, r.config, iface.TObject, pkgName, method.Name(), pm, signature)
		}
	}

	// Write signature mocks
	for _, signature := range pack.Signatures {
		methodName := strings.Replace(r.config.MockFunctionMethodNameFormat, lib.FuncNameToken, signature.TObject.Name(), -1)
		writeStruct(buf, r.config, signature.TObject, pkgName)
		writeMethod(buf, r.config, signature.TObject, pkgName, methodName, pm, signature.TSignature)
	}

	// Format generated code
	sourceCode, err := format.Source(buf.Bytes())
	if err != nil {
		log.Errorf("Unformattable code: %s", buf.Bytes())
		panic(err)
	}

	// Flush code to disk
	mockPath := filepath.Join(project.GoSrcAbsPath, pkg.Path(), r.config.MockFileName)
	log.Debugf("Creating file: %s", mockPath)
	if err := os.MkdirAll(filepath.Dir(mockPath), 0755); err != nil {
		panic(err)
	}
	mockFile, err := os.Create(mockPath)
	if err != nil {
		panic(err)
	}
	if _, err := mockFile.Write(sourceCode); err != nil {
		panic(err)
	}
	if err := mockFile.Close(); err != nil {
		panic(err)
	}

	project.GenAbsPaths = append(project.GenAbsPaths, mockPath)
}

func writeStruct(buf *lib.Bufferw, config lib.Config, original types.Object, pkgName string) {
	interfaceName := original.Name()
	mockName := strings.Replace(config.StructNameFormat, lib.InterfaceNameToken, interfaceName, -1)

	// Write constructor
	buf.WriteString(fmt.Sprintf("// New%s creates a new %s\n", mockName, mockName))
	buf.WriteString(fmt.Sprintf("func New%s() *%s {return &%s{}}\n", mockName, mockName, mockName))

	// Write struct
	log.Debugf("Writing struct: %s", mockName)
	buf.WriteString(fmt.Sprintf("// %s implements %s.%s\n", mockName, pkgName, interfaceName))
	buf.WriteString(fmt.Sprintf("type %s struct { mock.Mock }\n", mockName))
}

func writeMethod(buf *lib.Bufferw, config lib.Config, original types.Object, pkgName, methodName string, pm *lib.PackageMapper, signature *types.Signature) {
	interfaceName := original.Name()
	mockName := strings.Replace(config.StructNameFormat, lib.InterfaceNameToken, interfaceName, -1)

	paramNames := []string{}
	paramExprs := []string{}
	for j := 0; j < signature.Params().Len(); j++ {
		param := signature.Params().At(j)
		varName := "p" + strconv.Itoa(j)
		paramNames = append(paramNames, varName)
		paramExprs = append(paramExprs, varName+" "+pm.ObjectTypeString(param))
	}

	if signature.Variadic() { // Variadic method? Replace last parameter's [] with ... ("p1 []int" -> "p1 ...int")
		pos := signature.Params().Len() - 1
		paramExprs[pos] = strings.Replace(paramExprs[pos], "[]", "...", 1)
	}

	returnNames := []string{}
	returnTypes := []string{}
	verifyReturnLines := []string{}
	for j := 0; j < signature.Results().Len(); j++ {
		result := signature.Results().At(j)
		resultTypeString := pm.ObjectTypeString(result)
		varName := "ret" + strconv.Itoa(j)
		returnNames = append(returnNames, varName)
		returnTypes = append(returnTypes, resultTypeString)
		verifyReturnLine := fmt.Sprintf("var %s %s; if a := ret.Get(%d); a != nil { %s = a.(%s) }\n", varName, resultTypeString, j, varName, resultTypeString)
		verifyReturnLines = append(verifyReturnLines, verifyReturnLine)
	}

	log.Debugf("Writing method: %s()", methodName)
	commentLine := fmt.Sprintf("// %s implements (%s.%s).%s\n", methodName, pkgName, interfaceName, methodName)
	signatureLine := fmt.Sprintf("func (r *%s) %s(%s) (%s) {\n", mockName, methodName, strings.Join(paramExprs, ", "), strings.Join(returnTypes, ","))
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
