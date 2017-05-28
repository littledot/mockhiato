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
		if !strings.HasPrefix(packagePath, project.PackagePath) { // External dependency? Skip
			continue
		}
		if strings.HasPrefix(packagePath, project.VendorPath) { // Vendor dependency? Skip
			continue
		}
		interfaces := getInterfaces(allPackage.Info.Defs)
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

		for _, iface := range interfaces {
			r.allMockedInterfaces[iface.TObject] = iface
		}
		log.Debugf("Found package %s with %d interfaces", packagePath, num)
	}

	// Find interfaces used by project
	genPackage := &lib.GeneratedPackage{}
	genPackage.ContextPath = project.DependentMocksPath
	genPackage.ContextName = filepath.Base(genPackage.ContextPath)
	genPackage.Interfaces = []*lib.Interface{}
	project.DependentPackage = genPackage

	for _, pack := range project.Packages { // Only need to inspect packages that needs to be mocked
		interfaces := getInterfaces(pack.PackageInfo.Uses)
		for _, iface := range interfaces {
			if _, exists := r.allMockedInterfaces[iface.TObject]; !exists { // Interface not mocked before? Mock it
				genPackage.Interfaces = append(genPackage.Interfaces, iface)
				r.allMockedInterfaces[iface.TObject] = iface
			}
		}
	}

	logTypeCheckProjectResults(project)
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
	sort.Sort(byInterfaceName(interfaces))
	return interfaces
}

// filterInterfaces returns interfaces.
func filterInterfaces(objs map[*ast.Ident]types.Object) map[*ast.Ident]types.Object {
	interfaces := map[*ast.Ident]types.Object{}
	for ident, obj := range objs {
		if obj == nil {
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

	// Report stats for interfaces defined by project
	for _, pack := range project.Packages {
		ifaces := []string{}
		for _, iface := range pack.Interfaces {
			ifaces = append(ifaces, iface.TObject.Name())
		}
		log.Infof("Type checker found %d interface(s) defined in package %s: %s",
			len(ifaces), pack.Context.Path(), strings.Join(ifaces, ", "))
	}

	// Report stats for interfaces used by project
	dep := project.DependentPackage
	ifaces := []string{}
	for _, iface := range dep.Interfaces {
		ifaces = append(ifaces, iface.TObject.Name())
	}
	log.Infof("Type checker found %d interface(s) used by project: %s",
		len(dep.Interfaces), strings.Join(ifaces, ", "))
}

func logGenerateMocksResults(project *lib.Project) {
	log.Infof("Generate mocks complete")
	for _, path := range project.GenAbsPaths {
		log.Infof("Generated %s", path)
	}
}
