package generate

import (
	"go/ast"
	"go/types"

	"github.com/littledot/mockhiato/lib"
	log "github.com/sirupsen/logrus"
)

// MockFinder finds objects that can be mocked.
type MockFinder struct {
	packages map[*types.Package]*lib.Package

	// allMockables indexes interfaces that will be mocked
	allMockables map[types.Object]int
}

// NewMockFinder constructor.
func NewMockFinder() *MockFinder {
	return &MockFinder{
		packages:     map[*types.Package]*lib.Package{},
		allMockables: map[types.Object]int{},
	}
}

// findMockable indexes objects that can be mocked.
func (r *MockFinder) findMockable(objs map[*ast.Ident]types.Object) {
	for _, obj := range objs {
		if obj == nil {
			continue
		}
		if obj.Pkg() == nil {
			continue
		}
		if _, ok := obj.(*types.TypeName); !ok {
			continue
		}
		if _, exists := r.allMockables[obj]; exists { // Object already indexed? Skip
			log.Debugf("Ignore object %s because it is already indexed", obj.Name())
			continue
		}

		switch v := obj.Type().Underlying().(type) {
		case *types.Interface:
			iface := &lib.Interface{
				TObject:    obj,
				TInterface: v.Complete(),
			}
			r.recordInterface(iface)

		case *types.Signature:
			signature := &lib.Signature{
				TObject:    obj,
				TSignature: v,
			}
			r.recordSignature(signature)
		}
	}
}

func (r *MockFinder) recordInterface(iface *lib.Interface) {
	pkg := r.recordPackage(iface.TObject)
	pkg.Interfaces = append(pkg.Interfaces, iface)
	log.Debugf("Record interface: %s", iface.TObject.Name())
}

func (r *MockFinder) recordSignature(signature *lib.Signature) {
	pkg := r.recordPackage(signature.TObject)
	pkg.Signatures = append(pkg.Signatures, signature)
	log.Debugf("Record signature: %s", signature.TObject.Name())
}

func (r *MockFinder) recordPackage(obj types.Object) *lib.Package {
	r.allMockables[obj]++

	tPkg := obj.Pkg()
	pkg := r.packages[tPkg]
	if pkg == nil {
		pkg = lib.NewPackage(tPkg)
		r.packages[tPkg] = pkg
		log.Debugf("Record package: %s (%s)", tPkg.Name(), tPkg.Path())
	}
	return pkg
}
