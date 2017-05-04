package generate

import (
	"github.com/littledot/mockhiato/lib"
	log "github.com/sirupsen/logrus"
)

// Run executes the command.
func Run(config lib.Config) {
	log.Infof("Running generate")

	project := &lib.Project{}
	oracle := NewOracle(config)
	oracle.ScanProject(project)
	oracle.TypeCheckProject(project)
	oracle.GenerateMocks(project)

	log.Infof("Enjoy!")
}
