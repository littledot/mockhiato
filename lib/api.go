package lib

import (
	"go/types"
	"os"

	"golang.org/x/tools/go/loader"
)

// Config configures Mockhiato behavior.
type Config struct {
	Verbose      bool
	ProjectPath  string
	MockFileName string
}

// Formatter formats mock files.
type Formatter interface {
	GenerateMocks(project *Project)
	IsMockFile(file *os.File) bool
}

// Project represents the project Mockhiato is operating on.
type Project struct {
	// ProjectAbsPath is the project's absolute path.
	ProjectAbsPath string
	// GoAbsPath is $GOPATH.
	GoAbsPath string
	// GoSrcAbsPath is $GOPATH/src.
	GoSrcAbsPath string
	// PackagePath is the project's package path, which should be the relative path to $GOPATH/src.
	PackagePath string
	// VendorPath is the project's vedor path, which should be PackagePath/vendor
	VendorPath string

	// Program is the loaded project
	Program *loader.Program
	// Packages is a list of packages with interfaces that can be mocked.
	Packages []*Package

	//GenPaths contains a list of generated file paths
	GenAbsPaths []string
}

// NewProject creates a new Project.
func NewProject() *Project {
	return &Project{}
}

// Package contains metadata for a package discovered in the project tree. Formatters rely on this to generate mocks.
type Package struct {
	// PackageInfo is the loader's package info
	PackageInfo *loader.PackageInfo
	// Context is the package of the project.
	Context *types.Package
	// Interfaces contains interface definitions found in the package.
	Interfaces []*Interface
}

// Interface contains metadata for an interface definition. Formatters rely on this to generate mocks.
type Interface struct {
	TObject    types.Object
	TInterface *types.Interface
}
