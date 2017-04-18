package generate

import (
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"gitlab.com/littledot/mockhiato/lib"
	"gitlab.com/littledot/mockhiato/lib/plugin/github.com/stretchr/testify"
	"gopkg.in/yaml.v2"
)

func Run() {
	oracle := NewOracle()
	project := oracle.ScanProject()
	spew.Dump(project)
	spec := oracle.TypeCheckProject(project)
	spew.Dump(spec)
	oracle.GenerateMocks(spec)
}

// newDefaultConfig provides default behavior
func newDefaultConfig() *lib.Config {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return &lib.Config{
		ProjectPath: wd,
		IgnorePaths: []string{"vendor"},
	}
}

func NewConfig() *lib.Config {
	config := newDefaultConfig()
	parseConfigYaml(config)

	spew.Dump(config)
	return config
}

func parseConfigYaml(config *lib.Config) {
	// Unmarshal config yaml
	configFile, err := os.Open("mockhiato.yaml")
	if err != nil {
		return
	}
	defer configFile.Close()

	configBytes, err := ioutil.ReadAll(configFile)
	if err != nil {
		return
	}
	if err := yaml.Unmarshal(configBytes, config); err != nil {
		return
	}
}

// Oracle parses Go projects, looking for interfaces to mock.
type Oracle struct {
	lib.Formatter

	config lib.Config
}

func NewOracle() *Oracle {
	config := *NewConfig()
	oracle := &Oracle{
		Formatter: testify.NewTestifyFormatter(config),
		config:    config,
	}
	return oracle
}

// ScanProject walks the project directory, indexing valid Go source code
func (r *Oracle) ScanProject() *lib.Project {
	project := lib.NewProject()
	err := filepath.Walk(r.config.ProjectPath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil { // Something wrong? Skip
			return nil
		}
		if !strings.HasSuffix(filePath, ".go") { // Not Go source? Skip
			return nil
		}
		if strings.HasSuffix(filePath, "_test.go") { // Go test? Skip
			return nil
		}

		relPath, err := filepath.Rel(r.config.ProjectPath, filePath)
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

		dirPath := filepath.Dir(filePath)
		project.SourceSet[dirPath] = append(project.SourceSet[dirPath], filePath)
		return nil
	})
	if err != nil {
		panic(err)
	}
	return project
}

func (r *Oracle) TypeCheckProject(project *lib.Project) *lib.Spec {
	spec := &lib.Spec{}
	for dir, sources := range project.SourceSet {
		pack := r.typeCheckSources(dir, sources)
		spec.Packages = append(spec.Packages, pack)
	}
	return spec
}

func (r *Oracle) typeCheckSources(dir string, sources []string) *lib.Package {
	// Build package AST from sources
	fset := token.NewFileSet()
	astFiles := []*ast.File{}
	for _, source := range sources {
		astFile, err := parser.ParseFile(fset, source, nil, parser.AllErrors)
		if err != nil {
			panic(err)
		}
		astFiles = append(astFiles, astFile)
	}

	// Type-check package AST
	typeChecker := &types.Config{
		Importer: importer.Default(),
	}
	info := &types.Info{
		Defs: map[*ast.Ident]types.Object{},
		Uses: map[*ast.Ident]types.Object{},
	}
	if _, err := typeChecker.Check(dir, fset, astFiles, info); err != nil {
		panic(err)
	}

	pack := &lib.Package{}

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
			Object:    def,
			Interface: interfaceDef,
		}
		pack.Interfaces = append(pack.Interfaces, iface)
	}

	// Index imports used by the package
	for _, use := range info.Uses {
		pkgNameUse, ok := use.(*types.PkgName)
		if !ok {
			continue
		}
		pack.Imports = append(pack.Imports, pkgNameUse.Imported())
	}

	return pack
}
