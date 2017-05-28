package example

import (
	"os"

	"gitgud.com/bigdot"
)

// Target is an interface that should be mocked.
type Target interface {
	// Should generate mocks for dependent interfaces in /mocks
	GoInterface(fi os.FileInfo) (err error)
	VendorInterface(bigdot.VendorDep)

	// Should not identify B as dependent interface and generate mocks in /mocks
	DefineAndUse(B)
}

// B B
type B interface{}
