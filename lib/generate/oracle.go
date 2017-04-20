package generate

import (
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"strings"

	"gitlab.com/littledot/mockhiato/lib"
	"gitlab.com/littledot/mockhiato/lib/plugin/github.com/stretchr/testify"
)

// Oracle parses Go projects, looking for interfaces to mock.
type Oracle struct {
	lib.Formatter

	config lib.Config
}

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
	r.Formatter.SetProjectPackage(strings.Split(projectPath, "src/")[1])

	project := lib.NewProject()
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

		packagePath := strings.Split(filepath.Dir(filePath), "src/")[1]
		project.SourceSet[packagePath] = append(project.SourceSet[packagePath], filePath)
		return nil
	})
	if err != nil {
		panic(err)
	}
	return project
}

func (r *Oracle) TypeCheckProject(project *lib.Project) *lib.Spec {
	spec := &lib.Spec{}
	for packagePath, sourcePaths := range project.SourceSet {
		pack := r.typeCheckSources(packagePath, sourcePaths)
		spec.Packages = append(spec.Packages, pack)
	}
	return spec
}

func (r *Oracle) typeCheckSources(packagePath string, sourcePaths []string) *lib.Package {
	// Build package AST from sources
	fset := token.NewFileSet()
	astFiles := []*ast.File{}
	for _, sourcePath := range sourcePaths {
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
	tPackage, err := typeChecker.Check(packagePath, fset, astFiles, info)
	if err != nil {
		panic(err)
	}

	pack := &lib.Package{}
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

	return pack
}
