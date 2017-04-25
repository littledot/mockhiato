package generate

import (
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"sort"
	"strings"

	log "github.com/sirupsen/logrus"
	"gitlab.com/littledot/mockhiato/lib"
	"gitlab.com/littledot/mockhiato/lib/plugin/github.com/stretchr/testify"
)

// Oracle parses Go projects, looking for interfaces to mock.
type Oracle struct {
	lib.Formatter

	config lib.Config
}

// NewOracle creates a new oracle.
func NewOracle(config lib.Config) *Oracle {
	oracle := &Oracle{
		Formatter: testify.NewTestifyFormatter(config),
		config:    config,
	}
	return oracle
}

// ScanProject walks the project directory, indexing valid Go source code
func (r *Oracle) ScanProject() *lib.Project {
	projectPath, err := filepath.Abs(r.config.ProjectPath)
	if err != nil {
		panic(err)
	}
	log.Infof(`Project path is %s`, projectPath)

	project := lib.NewProject()
	project.PackagePath = lib.GetPackagePath(projectPath)
	log.Infof(`Package path is %s`, project.PackagePath)
	project.VendorPath = filepath.Join(project.PackagePath, "vendor")
	log.Infof(`Vendor path is %s`, project.VendorPath)

	err = filepath.Walk(projectPath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil { // Something wrong? Skip
			return nil
		}
		if !strings.HasSuffix(filePath, ".go") { // Not Go source? Skip
			return nil
		}
		if strings.HasSuffix(filePath, "_test.go") { // Go test? Skip
			return nil
		}

		relPath, err := filepath.Rel(projectPath, filePath)
		if err != nil {
			return nil
		}
		for _, ignorePath := range r.config.IgnorePaths {
			if strings.HasPrefix(relPath, ignorePath) { // Part of ignored paths? Skip
				return filepath.SkipDir
			}
		}

		file, err := os.Open(filePath)
		if err != nil {
			return nil
		}
		defer file.Close()
		if r.IsMockFile(file) { // Formatter says its a mock? Skip
			return nil
		}

		packageAbsPath := filepath.Dir(filePath)
		pack := project.PathToPackage[packageAbsPath]
		if pack == nil {
			pack = &lib.Package{}
			pack.AbsPath = packageAbsPath
			pack.PackagePath = lib.GetPackagePath(packageAbsPath)
			project.PathToPackage[packageAbsPath] = pack
		}
		pack.SourceAbsPaths = append(pack.SourceAbsPaths, filePath)
		return nil
	})
	if err != nil {
		panic(err)
	}
	return project
}

// TypeCheckProject type-checks Go source code.
func (r *Oracle) TypeCheckProject(project *lib.Project) {
	for _, pack := range project.PathToPackage {
		r.typeCheckSources(project, pack)
	}
}

func (r *Oracle) typeCheckSources(project *lib.Project, pack *lib.Package) {
	// Ignore Pain - Ignore sources which cannot be type-checked
	defer recoverTypeCheckErrors(project, pack)()

	// Build package AST from sources
	fset := token.NewFileSet()
	astFiles := []*ast.File{}
	for _, sourcePath := range pack.SourceAbsPaths {
		astFile, err := parser.ParseFile(fset, sourcePath, nil, parser.AllErrors)
		if err != nil {
			panic(err)
		}
		astFiles = append(astFiles, astFile)
	}

	// Type-check package AST
	typeChecker := &types.Config{}
	typeChecker.Importer = importer.Default()
	info := &types.Info{}
	info.Defs = map[*ast.Ident]types.Object{}
	tPackage, err := typeChecker.Check(pack.PackagePath, fset, astFiles, info)
	if err != nil {
		panic(err)
	}
	pack.Context = tPackage

	// Index imports used by the package
	pack.Imports = tPackage.Imports()

	// Index interfaces defined in the package
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
		pack.Interfaces = append(pack.Interfaces, iface)
	}
	sort.Slice(pack.Interfaces, func(i, j int) bool { return pack.Interfaces[i].TObject.Name() < pack.Interfaces[j].TObject.Name() })
}

func recoverTypeCheckErrors(project *lib.Project, pack *lib.Package) func() {
	return func() {
		if err := lib.Err(recover()); err != nil {
			log.Error(err.ErrorStack())
			log.Errorf(`ignoring %s due to type check errors; see log for details`, pack.PackagePath)
			project.AllErrors = append(project.AllErrors, err)
			pack.Error = err
		}
	}
}
