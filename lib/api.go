package lib

import (
	"go/types"
	"os"
)

// Config configures Mockhiato behavior
type Config struct {
	Verbose     bool
	ProjectPath string
	IgnorePaths []string
	OutputPath  string
}

// Formatter formats mock files.
type Formatter interface {
	SetProjectPackage(projectPackage string)
	GenerateMocks(spec *Spec)
	IsMockFile(file *os.File) bool
}

type Project struct {
	SourceSet map[string][]string
}

// NewProject creates a new Project.
func NewProject() *Project {
	return &Project{
		SourceSet: map[string][]string{},
	}
}

// Spec contains a list of packages discovered in the project tree. Formatters rely on this to generate mocks.
type Spec struct {
	Packages []*Package
}

// Package contains metadata for a package discovered in the project tree. Formatters rely on this to generate mocks.
type Package struct {
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
