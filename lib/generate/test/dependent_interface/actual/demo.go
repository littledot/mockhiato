package example

import (
	"os"

	"golang.org/x/text"
	"golang.org/x/tools"
)

// Target is an interface that should be mocked.
type Target interface {
	// Should identify dependent interface and generate mocks in /mocks
	GoInterface(fi os.FileInfo) (err error)
	// Should not identify B as dependent interface and generate mocks in /mocks
	DefineAndUse(B)

	// Should recreate their original package tree in /mocks so that interfaces with the same name are not duplicated in the same package "mocks"
	VendorInterface(text.VendorDep)
	VendorInterfaceWithSameName(tools.VendorDep)
}

// B B
type B interface{}
