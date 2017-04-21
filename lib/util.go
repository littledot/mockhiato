package lib

import (
	"go/types"
	"path/filepath"
	"strconv"
	"strings"
)

// GetPackagePath returns the package path relative to $GOPATH/src
func GetPackagePath(path string) string {
	absPath, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	const src = "src/"
	srcPos := strings.Index(absPath, src)
	return absPath[srcPos+len(src) : len(absPath)]
}

// PackageFormatter indexes packages and generates package aliases if required.
type PackageFormatter struct {
	Context     *types.Package
	PathToAlias map[string]string
	usedAliases map[string]int
}

// NewPackageFormatter creates a new PackageFormatter.
func NewPackageFormatter(context *types.Package) *PackageFormatter {
	return &PackageFormatter{
		Context:     context,
		PathToAlias: map[string]string{},
		usedAliases: map[string]int{},
	}
}

// IndexImports indexes tPackage's dependencies.
func (r *PackageFormatter) IndexImports(tPackage *types.Package) {
	for _, tImport := range tPackage.Imports() {
		r.RecordDependency(tImport)
	}
}

// IndexInterface indexes tInterface's dependencies.
func (r *PackageFormatter) IndexInterface(tInterface *types.Interface) {
	for i := 0; i < tInterface.NumMethods(); i++ {
		r.IndexFunc(tInterface.Method(i))
	}
}

// IndexFunc indexes tFunc's dependencies.
func (r *PackageFormatter) IndexFunc(tFunc *types.Func) {
	r.IndexSignature(tFunc.Type().(*types.Signature))
}

// IndexSignature indexes tSignature's dependencies.
func (r *PackageFormatter) IndexSignature(tSignature *types.Signature) {
	r.IndexTuple(tSignature.Params())
	r.IndexTuple(tSignature.Results())
}

// IndexTuple indexes tTuple's dependencies
func (r *PackageFormatter) IndexTuple(tTuple *types.Tuple) {
	for i := 0; i < tTuple.Len(); i++ {
		tVar := tTuple.At(i)
		types.ObjectString(tVar, r.RecordDependency)
	}
}

// RecordDependency indexes tPackage as a dependency and returns its name.
func (r *PackageFormatter) RecordDependency(tPackage *types.Package) string {
	if tPackage.Path() == r.Context.Path() {
		return ""
	}

	if alias, exists := r.PathToAlias[tPackage.Path()]; exists {
		return alias
	}

	name := (*types.Package).Name(tPackage)
	alias := name
	for i := 1; ; i++ {
		if _, exists := r.usedAliases[alias]; !exists {
			r.usedAliases[alias] = 1
			r.PathToAlias[tPackage.Path()] = alias
			break
		}
		alias = name + strconv.Itoa(i)
	}

	return alias
}

// ObjectTypeString returns the string representation of tObject's type.
func (r *PackageFormatter) ObjectTypeString(tObject types.Object) string {
	return r.TypeString(tObject.Type())
}

// TypeString returns the string representation of tType.
func (r *PackageFormatter) TypeString(tType types.Type) string {
	return types.TypeString(tType, r.RecordDependency)
}
