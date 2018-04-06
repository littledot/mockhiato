package lib

// Config configures Mockhiato behavior. See mockhiato/cmd for documentation.
type Config struct {
	Verbose                    bool
	ProjectPath                string
	MockFileName               string
	DependentMocksPath         string
	StructNameFormat           string
	DependentPackageNameFormat string
}
