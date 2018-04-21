package example

import (
	"database/sql/driver"
	"flag"
	"go/constant"
	"os"
)

// Target is an interface that should be mocked.
type Target interface {
	// Should identify dependent interface and generate mocks in /mocks
	GoInterface(fi os.FileInfo) (err error)
	// Should not identify B as dependent interface and generate mocks in /mocks
	DefineAndUse(B)

	// Should recreate their original package tree in /mocks so that interfaces with the same name are not duplicated under the same package "mocks"
	InterfaceWithSameName(flag.Value)
	InterfaceWithSameName2(driver.Value)
	InterfaceWithSameName3(constant.Value)
}

// B B
type B interface{}
