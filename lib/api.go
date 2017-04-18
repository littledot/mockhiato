package lib

import "go/types"

type Formatter interface {
	GenerateMocks(spec *Spec)
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
	Interfaces []*Interface
	Imports    []*types.Package
}

type Interface struct {
	Object    types.Object
	Interface *types.Interface
}
