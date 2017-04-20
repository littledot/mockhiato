package lib

import (
	"go/types"
	"os"
)

type Config struct {
	Verbose     bool
	ProjectPath string
	IgnorePaths []string
	OutputPath  string
}

type Formatter interface {
	SetProjectPackage(projectPackage string)
	GenerateMocks(spec *Spec)
	IsMockFile(file *os.File) bool
}

type Project struct {
	SourceSet map[string][]string
}

func NewProject() *Project {
	return &Project{
		SourceSet: map[string][]string{},
	}
}

type Spec struct {
	Packages []*Package
}

type Package struct {
	TPackage   *types.Package
	TImports   []*types.Package
	Interfaces []*Interface
}

type Interface struct {
	TObject    types.Object
	TInterface *types.Interface
}
