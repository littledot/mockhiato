package main

import (
	"bytes"
	"encoding/json"
	"image/png"
	"os"

	"gitlab.com/littledot/mockhiato/demo"
)

type AAA interface {
	Hello(a string) (b int, err error)
	World(c int, a ...string) (b map[int]interface{})
	Yes(fi os.FileInfo) (pnge *png.Encoder, err error)
	No([]*json.Decoder) *bytes.Buffer
	demo.B
	// os.FileInfo
}

type Beta interface {
	No(*json.Decoder) *bytes.Buffer
}
