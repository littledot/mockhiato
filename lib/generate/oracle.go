package generate

import (
	"go/parser"
	"go/types"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/tools/go/loader"

	log "github.com/sirupsen/logrus"
	"gitlab.com/littledot/mockhiato/lib"
	"gitlab.com/littledot/mockhiato/lib/plugin/github.com/stretchr/testify"
)

// Oracle parses Go projects, looking for interfaces to mock.
type Oracle struct {
	formatter lib.Formatter

	config lib.Config
}

// NewOracle creates a new oracle.
func NewOracle(config lib.Config) *Oracle {
	oracle := &Oracle{
		formatter: testify.NewTestifyFormatter(config),
		config:    config,
	}
	return oracle
}

// ScanProject walks the project directory, indexing valid Go source code
func (r *Oracle) ScanProject(project *lib.Project) {
	log.Debugf("Scanning project...")
	projectPath, err := filepath.Abs(r.config.ProjectPath)
	if err != nil {
		panic(err)
	}
	project.ProjectAbsPath = projectPath

	const src = "src/"
	srcPos := strings.Index(projectPath, src)
	project.GoAbsPath = projectPath[0 : srcPos-1]
	project.GoSrcAbsPath = projectPath[0 : srcPos+len(src)-1]
	project.PackagePath = projectPath[srcPos+len(src) : len(projectPath)]
	project.VendorPath = filepath.Join(project.PackagePath, "vendor")

	logScanProjectResults(project)
}

// TypeCheckProject type-checks Go source code.
func (r *Oracle) TypeCheckProject(project *lib.Project) {
	log.Debugf("Type checking project...")
	goloader := &loader.Config{}
	goloader.AllowErrors = true
	goloader.ParserMode = parser.AllErrors
	goloader.Import(project.PackagePath)
	program, err := goloader.Load()
	if err != nil {
		panic(err)
	}

	for _, allPackage := range program.AllPackages {
		packagePath := allPackage.Pkg.Path()
		if !strings.HasPrefix(packagePath, project.PackagePath) { // External dependency? Skip
			continue
		}
		if strings.HasPrefix(packagePath, project.VendorPath) { // Vendor dependency? Skip
			continue
		}
		interfaces := getDefinedInterfaces(&allPackage.Info)
		num := len(interfaces)
		if num == 0 { // 0 interfaces defined? Skip
			log.Debugf("Ignore package %s because it has 0 interfaces", packagePath)
			continue
		}

		pack := &lib.Package{}
		pack.PackageInfo = allPackage
		pack.Context = allPackage.Pkg
		pack.Interfaces = interfaces

		project.Packages = append(project.Packages, pack)
		log.Debugf("Found package %s with %d interfaces", packagePath, num)
	}

	logTypeCheckProjectResults(project)
}

// GenerateMocks generate mocks for the project
func (r *Oracle) GenerateMocks(project *lib.Project) {
	log.Debugf("Generating mocks...")
	r.formatter.GenerateMocks(project)
	logGenerateMocksResults(project)
}

// getDefinedInterfaces returns defined interfaces in the package.
func getDefinedInterfaces(info *types.Info) []*lib.Interface {
	interfaces := []*lib.Interface{}
	for _, def := range info.Defs {
		if def == nil {
			continue
		}
		if _, ok := def.(*types.TypeName); !ok {
			continue
		}
		if !types.IsInterface(def.Type()) {
			continue
		}

		interfaceDef := def.Type().Underlying().(*types.Interface).Complete()
		iface := &lib.Interface{
			TObject:    def,
			TInterface: interfaceDef,
		}
		interfaces = append(interfaces, iface)
	}

	sort.Slice(interfaces, func(i, j int) bool { return interfaces[i].TObject.Name() < interfaces[j].TObject.Name() })
	return interfaces
}

func logScanProjectResults(project *lib.Project) {
	log.Infof("Scan project complete")
	log.Infof("Project path is %s", project.ProjectAbsPath)
	log.Infof("GOPATH is %s", project.GoAbsPath)
	log.Infof("Project package is %s", project.PackagePath)
	log.Infof("Project vendor path is %s", project.VendorPath)
}

func logTypeCheckProjectResults(project *lib.Project) {
	log.Infof("Type check complete")
	for _, pack := range project.Packages {
		ifaces := []string{}
		for _, iface := range pack.Interfaces {
			ifaces = append(ifaces, iface.TObject.Name())
		}
		log.Infof("Type checker found %d interface(s) in package %s: %s",
			len(ifaces), pack.Context.Path(), strings.Join(ifaces, ", "))
	}
}

func logGenerateMocksResults(project *lib.Project) {
	log.Infof("Generate mocks complete")
	for _, path := range project.GenAbsPaths {
		log.Infof("Generated %s", path)
	}
}
