package example

import (
	"github.com/littledot/mockhiato/test/same_import_name/expect/same"
	other_same "github.com/littledot/mockhiato/test/same_import_name/expect/same1"
)

// Target interface references 2 packages from different paths but with the same name. Mockhiato must be able to handle imports correctly for these kind of packages.
type Target interface {
	A(same.ThisInt, other_same.ThatInt)
}
