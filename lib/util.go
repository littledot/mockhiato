package lib

import (
	"bytes"
	"go/types"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

func GetTypeString(obj types.Object) string {
	spew.Dump(obj.Type())
	buf := &bytes.Buffer{}
	getTypeString(buf, obj.Type())
	return buf.String()
}

func getTypeString(buf *bytes.Buffer, obj types.Type) {
	switch objType := obj.(type) {
	case *types.Pointer:
		splits := strings.Split(objType.Elem().String(), "/")
		buf.WriteString("*")
		buf.WriteString(splits[len(splits)-1])
	case *types.Slice:
		buf.WriteString("[]")
		getTypeString(buf, objType.Elem())
	default:
		buf.WriteString(objType.String())
	}
}
