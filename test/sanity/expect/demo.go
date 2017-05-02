package example

import (
	"bytes"
	"encoding/json"
	"image/png"
	"os"
)

// A this package contains 2 interfaces that needs to be mocked.
type A interface {
	Hello(a string) (b int, err error)
	World(c int, a ...string) (b map[int]interface{})
	Yes(fi os.FileInfo) (pnge *png.Encoder, err error)
}

// B B
type B interface {
	No(*json.Decoder) *bytes.Buffer
}
