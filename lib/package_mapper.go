package lib

import (
	"go/types"
	"strconv"
)

// PackageMapper indexes packages and generates package aliases if required.
type PackageMapper struct {
	PkgPath     string
	PathToAlias map[string]string
	usedAliases map[string]int
}

// NewPackageMapper creates a new PackageMapper.
func NewPackageMapper(PkgPath string) *PackageMapper {
	return &PackageMapper{
		PkgPath,
		map[string]string{},
		map[string]int{},
	}
}

// IndexImports indexes tPkg's dependencies.
func (r *PackageMapper) IndexImports(tPkg *types.Package) {
	for _, tImport := range tPkg.Imports() {
		r.RecordDependency(tImport)
	}
}

// IndexInterface indexes tInterface's dependencies.
func (r *PackageMapper) IndexInterface(tInterface *types.Interface) {
	for i := 0; i < tInterface.NumMethods(); i++ {
		r.IndexFunc(tInterface.Method(i))
	}
}

// IndexFunc indexes tFunc's dependencies.
func (r *PackageMapper) IndexFunc(tFunc *types.Func) {
	r.IndexSignature(tFunc.Type().(*types.Signature))
}

// IndexSignature indexes tSignature's dependencies.
func (r *PackageMapper) IndexSignature(tSignature *types.Signature) {
	r.IndexTuple(tSignature.Params())
	r.IndexTuple(tSignature.Results())
}

// IndexTuple indexes tTuple's dependencies
func (r *PackageMapper) IndexTuple(tTuple *types.Tuple) {
	for i := 0; i < tTuple.Len(); i++ {
		tVar := tTuple.At(i)
		types.ObjectString(tVar, r.RecordDependency)
	}
}

// RecordDependency indexes tPkg as a dependency and returns its name.
func (r *PackageMapper) RecordDependency(tPkg *types.Package) string {
	if tPkg.Path() == r.PkgPath {
		return ""
	}

	if alias, exists := r.PathToAlias[tPkg.Path()]; exists {
		return alias
	}

	name := tPkg.Name()
	alias := name
	for i := 1; ; i++ {
		if _, exists := r.usedAliases[alias]; !exists {
			r.usedAliases[alias] = 1
			r.PathToAlias[tPkg.Path()] = alias
			break
		}
		alias = name + strconv.Itoa(i)
	}

	return alias
}

// ObjectTypeString returns the string representation of tObject's type.
func (r *PackageMapper) ObjectTypeString(tObject types.Object) string {
	return r.TypeString(tObject.Type())
}

// TypeString returns the string representation of tType.
func (r *PackageMapper) TypeString(tType types.Type) string {
	return types.TypeString(tType, r.RecordDependency)
}
