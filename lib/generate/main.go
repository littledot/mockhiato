package generate

import (
	log "github.com/sirupsen/logrus"
	"gitlab.com/littledot/mockhiato/lib"
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
