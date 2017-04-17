package lib

import (
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

type Gen struct {
}

func Run() {
	gen := &Gen{}
	sources := gen.scanFiles()
	spew.Dump(sources)
	// gen.parseFile(sources)
	loaded := gen.loader(sources)
	gen.gen(loaded)
}

type scanResult struct {
	sources []string
}

func (r *Gen) scanFiles() *scanResult {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	result := &scanResult{}
	err = filepath.Walk(wd, func(path string, info os.FileInfo, err error) error {
		if err != nil { // Something wrong? Skip
			return nil
		}
		if !strings.HasSuffix(path, ".go") { // Not Go source? Skip
			return nil
		}
		if strings.HasSuffix(path, "_test.go") { // Go test? Skip
			return nil
		}

		rel, err := filepath.Rel(wd, path)
		if err != nil {
			return nil
		}
		if strings.HasPrefix(rel, "vendor") { // Vendor directories? Skip TODO: make this configurable
			return filepath.SkipDir
		}

		if !strings.Contains(path, "demo.go") {
			return nil
		}

		result.sources = append(result.sources, path)
		return nil
	})
	if err != nil {
		panic(err)
	}

	return result
}

type AstInterface struct {
	Name    string
	Methods []*AstFunction
}

type AstFunction struct {
	Name string
}

func (r *Gen) parseFile(scan *scanResult) {
	fset := token.NewFileSet()
	for _, source := range scan.sources {
		f, err := parser.ParseFile(fset, source, nil, parser.AllErrors)
		if err != nil {
			panic(err)
		}
		ast.Print(fset, f)

		ast.Inspect(f, func(node ast.Node) bool {
			if genDecl, ok := node.(*ast.GenDecl); ok {
				if genDecl.Tok == token.TYPE {
					if typeSpec, ok := genDecl.Specs[0].(*ast.TypeSpec); ok {
						if interfaceType, ok := typeSpec.Type.(*ast.InterfaceType); ok {
							astInterface := &AstInterface{}
							astInterface.Name = typeSpec.Name.Name
							spew.Dump("iface:", astInterface.Name)

							for _, method := range interfaceType.Methods.List {
								astFunction := &AstFunction{}
								astFunction.Name = method.Names[0].Name
								astInterface.Methods = append(astInterface.Methods, astFunction)
								spew.Dump("meth:", astFunction.Name)
							}
							spew.Dump(astInterface)
							return false
						}
					}
				}
			}
			return true
		})
	}
}

type loaderResult struct {
	types []*typeResult
}

type typeResult struct {
	object types.Object
	iface  *types.Interface
}

func (r *Gen) loader(scan *scanResult) *loaderResult {
	cfg := &types.Config{
		Importer: importer.Default(),
	}
	fset := token.NewFileSet()
	files := []*ast.File{}
	for _, source := range scan.sources {
		f, err := parser.ParseFile(fset, source, nil, parser.AllErrors)
		if err != nil {
			panic(err)
		}
		files = append(files, f)
	}

	info := &types.Info{
		Defs: map[*ast.Ident]types.Object{},
	}
	_, err := cfg.Check("", fset, files, info)
	if err != nil {
		panic(err)
	}
	// spew.Dump(pkg)
	spew.Dump(info)

	typeResults := []*typeResult{}
	for _, def := range info.Defs {
		if def == nil {
			continue
		}
		if def.Pkg() == nil {
			continue
		}
		if !types.IsInterface(def.Type()) {
			continue
		}

		interfaceDef := def.Type().Underlying().(*types.Interface).Complete()
		result := &typeResult{
			object: def,
			iface:  interfaceDef,
		}
		typeResults = append(typeResults, result)
		spew.Printf("Found iface %s with %d func\n", def.Name(), interfaceDef.NumMethods())
		for c := 0; c < interfaceDef.NumMethods(); c++ {
			method := interfaceDef.Method(c)
			spew.Dump(method.Pkg())
			signature := method.Type().(*types.Signature)
			spew.Dump(signature)
		}
	}
	return &loaderResult{
		types: typeResults,
	}
}

func (r *Gen) gen(loaded *loaderResult) {
	for _, load := range loaded.types {
		spew.Dump(types.ObjectString(load.object, nil))
		spew.Dump(types.TypeString(load.iface, nil))
	}
}
