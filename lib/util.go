package lib

import (
	"go/types"
	"strconv"
)

func TObjectTypeToString(obj types.Object) string {
	return TTypeToString(obj.Type())
}

func TTypeToString(obj types.Type) string {
	return types.TypeString(obj, (*types.Package).Name)
}

type PackageFormatter struct {
	PathToAlias map[string]string

	dependencies map[string]*types.Package
	usedAliases  map[string]int
}

func NewPackageFormatter() *PackageFormatter {
	return &PackageFormatter{
		dependencies: map[string]*types.Package{},
		PathToAlias:  map[string]string{},
		usedAliases:  map[string]int{},
	}
}

func (r *PackageFormatter) AnalyzePackage(tPackage *types.Package) {
	r.RecordDependency(tPackage)
	for _, tImport := range tPackage.Imports() {
		r.RecordDependency(tImport)
	}
}

func (r *PackageFormatter) AnalyzeInterface(tInterface *types.Interface) {
	for i := 0; i < tInterface.NumMethods(); i++ {
		r.AnalyzeFunc(tInterface.Method(i))
	}
}

func (r *PackageFormatter) AnalyzeFunc(tFunc *types.Func) {
	r.AnalyzeSignature(tFunc.Type().(*types.Signature))
}

func (r *PackageFormatter) AnalyzeSignature(tSignature *types.Signature) {
	r.AnalyzeTuple(tSignature.Params())
	r.AnalyzeTuple(tSignature.Results())
}

func (r *PackageFormatter) AnalyzeTuple(tTuple *types.Tuple) {
	for i := 0; i < tTuple.Len(); i++ {
		tVar := tTuple.At(i)
		types.ObjectString(tVar, r.RecordDependency)
	}
}

func (r *PackageFormatter) RecordDependency(tPackage *types.Package) string {
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

func (r *PackageFormatter) ObjectTypeString(t types.Object) string {
	return r.TypeString(t.Type())
}

func (r *PackageFormatter) TypeString(tType types.Type) string {
	return types.TypeString(tType, r.RecordDependency)
}
