package generate

import (
	log "github.com/sirupsen/logrus"
	"gitlab.com/littledot/mockhiato/lib"
)

// Run executes the command.
func Run(config lib.Config) {
	log.Info("Running generate")

	oracle := NewOracle(config)
	project := oracle.ScanProject()

	log.Info("Scan project complete")

	oracle.TypeCheckProject(project)

	log.Info("Type check complete")

	oracle.GenerateMocks(project)

	log.Info("Generate mocks complete")
}
