package lib

import (
	"go/types"
	"os"
)

// Config configures Mockhiato behavior.
type Config struct {
	Verbose     bool
	ProjectPath string
	IgnorePaths []string
	OutputPath  string
}

// Formatter formats mock files.
type Formatter interface {
	SetProjectPackage(projectPackage string)
	GenerateMocks(project *Project)
	IsMockFile(file *os.File) bool
}

// Project represents the project Mockhiato is operating on.
type Project struct {
	// PathToPackage contains a mapping of package's absolute path to the package.
	PathToPackage map[string]*Package
}

// NewProject creates a new Project.
func NewProject() *Project {
	return &Project{
		PathToPackage: map[string]*Package{},
	}
}

// Package contains metadata for a package discovered in the project tree. Formatters rely on this to generate mocks.
type Package struct {
	// AbsPath is the package's absolute path.
	AbsPath string
	// PackagePath is the package's path, which should be the relative path to $GOPATH/src/.
	PackagePath string
	// SourcePath contains paths for Go source code.
	SourcePaths []string

	// Context is the package of the project.
	Context *types.Package
	// Imports contains dependencies used by the package.
	Imports []*types.Package
	// Interfaces contains interface definitions found in the package.
	Interfaces []*Interface
}

// Interface contains metadata for an interface definition. Formatters rely on this to generate mocks.
type Interface struct {
	TObject    types.Object
	TInterface *types.Interface
}
