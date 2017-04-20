package generate

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"gitlab.com/littledot/mockhiato/lib"
)

func Run(config lib.Config) {
	oracle := NewOracle(config)
	project := oracle.ScanProject()
	fmt.Println("dumping project...")
	spew.Dump(project)
	spec := oracle.TypeCheckProject(project)
	fmt.Println("dumping spec...")
	spew.Dump(spec)
	oracle.GenerateMocks(spec)
}