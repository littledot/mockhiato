package generate

import "github.com/littledot/mockhiato/lib"

type byInterfaceName []*lib.Interface

func (r byInterfaceName) Len() int           { return len(r) }
func (r byInterfaceName) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r byInterfaceName) Less(i, j int) bool { return r[i].TObject.Name() < r[j].TObject.Name() }
