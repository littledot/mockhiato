package lib

import "go/types"

func TObjectTypeToString(obj types.Object) string {
	return TTypeToString(obj.Type())
}

func TTypeToString(obj types.Type) string {
	return types.TypeString(obj, (*types.Package).Name)
}
