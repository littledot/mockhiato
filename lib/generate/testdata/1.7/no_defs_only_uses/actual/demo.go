package example

import "io"

// Target should generate mocks for io.Reader even if it does define any mocks iteself
func Target(io.Reader) {}
