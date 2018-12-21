package generate

import (
	"go/parser"
	"path/filepath"
	"sort"
	"strings"

	"github.com/littledot/mockhiato/lib"
	"github.com/littledot/mockhiato/lib/plugin/github.com/stretchr/testify"
	log "github.com/sirupsen/logrus"
	"golang.org/x/tools/go/loader"
)

// Oracle parses Go projects, looking for interfaces to mock.
type Oracle struct {
	formatter lib.Formatter

	config lib.Config
}

// NewOracle creates a new oracle.
func NewOracle(config lib.Config) *Oracle {
	return &Oracle{
		testify.NewTestifyFormatter(config),
		config,
	}
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
	project.PackagePath = projectPath[srcPos+len(src):]
	project.VendorPath = filepath.Join(project.PackagePath, "vendor")
	project.DependentMocksPath = filepath.Join(project.PackagePath, r.config.DependentMocksPath)

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

	finder := NewMockFinder()
	// Find interfaces defined by project
	for _, pkgInfo := range program.AllPackages {
		pkgPath := pkgInfo.Pkg.Path()
		if lib.IsExternalDependency(project, pkgPath) { // External dependency? Skip
			continue
		}

		finder.findMockable(pkgInfo.Info.Defs)
		finder.findMockable(pkgInfo.Info.Uses)
	}

	project.Packages = finder.packages
	for _, pkg := range project.Packages {
		sort.Sort(byInterfaceName(pkg.Interfaces))
		sort.Sort(bySignatureName(pkg.Signatures))
	}

	logTypeCheckProjectResults(project)
}

// GenerateMocks generate mocks for the project
func (r *Oracle) GenerateMocks(project *lib.Project) {
	log.Debugf("Generating mocks...")
	r.formatter.GenerateMocks(project)
	logGenerateMocksResults(project)
}

func logScanProjectResults(project *lib.Project) {
	log.Infof("Scan project complete")
	log.Infof("Project path is %s", project.ProjectAbsPath)
	log.Infof("GOPATH is %s", project.GoAbsPath)
	log.Infof("Project package is %s", project.PackagePath)
	log.Infof("Project vendor path is %s", project.VendorPath)
	log.Infof("Dependent mocks path is %s", project.DependentMocksPath)
}

func logTypeCheckProjectResults(project *lib.Project) {
	log.Infof("Type check complete")
	for _, dep := range project.Packages {
		log.Infof("Type checker found %d interface(s) in %s:", len(dep.Interfaces), dep.TPackage.Path())
		for _, iface := range dep.Interfaces {
			log.Infof("\t%s", iface.TObject.Name())
		}
	}
}

func logGenerateMocksResults(project *lib.Project) {
	log.Infof("Generate mocks complete")
	for _, path := range project.GenAbsPaths {
		log.Infof("Generated %s", path)
	}
}
