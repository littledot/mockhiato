package generate

import (
	"go/ast"
	"go/parser"
	"go/types"
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
	// allMockedInterfaces indexes interfaces that will be mocked
	allMockedInterfaces map[types.Object]*lib.Interface
}

// NewOracle creates a new oracle.
func NewOracle(config lib.Config) *Oracle {
	return &Oracle{
		testify.NewTestifyFormatter(config),
		config,
		map[types.Object]*lib.Interface{},
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

	// Find interfaces defined by project
	for _, allPackage := range program.AllPackages {
		packagePath := allPackage.Pkg.Path()
		if lib.IsExternalDependency(project, packagePath) { // External dependency? Skip
			continue
		}

		defInterfaces := getInterfaces(allPackage.Info.Defs)
		num := len(defInterfaces)
		if num == 0 { // 0 interfaces defined? Skip
			log.Debugf("Ignore package %s because it has 0 interfaces", packagePath)
			continue
		}
		log.Debugf("Check %d interface(s) defined in %s", num, packagePath)
		r.recordInterfaces(project, defInterfaces)

		useInterfaces := getInterfaces(allPackage.Info.Uses)
		log.Debugf("Check %d interface(s) used in %s", len(useInterfaces), packagePath)
		r.recordInterfaces(project, useInterfaces)
	}

	for _, pack := range project.Packages {
		sort.Sort(byInterfaceName(pack.Interfaces))
	}

	logTypeCheckProjectResults(project)
}

func (r *Oracle) recordInterfaces(project *lib.Project, interfaces []*lib.Interface) {
	for _, iface := range interfaces {
		if _, exists := r.allMockedInterfaces[iface.TObject]; exists { // Interface already indexed? Skip
			log.Debugf("Ignore interface %s because it is already indexed", iface.TObject.Name())
			continue
		}

		context := iface.TObject.Pkg()
		depPackage := project.Packages[context]
		if depPackage == nil {
			depPackage = &lib.Package{}
			depPackage.Context = context
			depPackage.Interfaces = []*lib.Interface{}
			project.Packages[context] = depPackage
			log.Debugf("Record package: %s (%s)", context.Name(), context.Path())

		}
		depPackage.Interfaces = append(depPackage.Interfaces, iface)
		log.Debugf("Record interface: %s", iface.TObject.Name())

		r.allMockedInterfaces[iface.TObject] = iface
	}
}

// GenerateMocks generate mocks for the project
func (r *Oracle) GenerateMocks(project *lib.Project) {
	log.Debugf("Generating mocks...")
	r.formatter.GenerateMocks(project)
	logGenerateMocksResults(project)
}

// getInterfaces returns interfaces in the package.
func getInterfaces(objs map[*ast.Ident]types.Object) []*lib.Interface {
	interfaces := []*lib.Interface{}
	for _, def := range filterInterfaces(objs) {
		interfaceDef := def.Type().Underlying().(*types.Interface).Complete()
		iface := &lib.Interface{
			TObject:    def,
			TInterface: interfaceDef,
		}
		interfaces = append(interfaces, iface)
	}
	return interfaces
}

// filterInterfaces returns interfaces.
func filterInterfaces(objs map[*ast.Ident]types.Object) map[*ast.Ident]types.Object {
	interfaces := map[*ast.Ident]types.Object{}
	for ident, obj := range objs {
		if obj == nil {
			continue
		}
		if obj.Pkg() == nil {
			continue
		}
		if _, ok := obj.(*types.TypeName); !ok {
			continue
		}
		if !types.IsInterface(obj.Type()) {
			continue
		}
		interfaces[ident] = obj
	}
	return interfaces
}

type byInterfaceName []*lib.Interface

func (r byInterfaceName) Len() int           { return len(r) }
func (r byInterfaceName) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r byInterfaceName) Less(i, j int) bool { return r[i].TObject.Name() < r[j].TObject.Name() }

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
		log.Infof("Type checker found %d interface(s) in %s:", len(dep.Interfaces), dep.Context.Path())
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
