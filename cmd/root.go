package cmd

import (
	"strings"

	"github.com/littledot/mockhiato/lib"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "mockhiato",
	Short: "A mock generation CLI tool for the Go",
	Long:  `Mockhiato is a mock generation CLI tool for the Go programming language. It is designed to be fast, configurable and correct.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) {
	// 	spew.Dump(viper.AllSettings())
	// },
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is mockhiato.yaml)")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.PersistentFlags().BoolP("Verbose", "v", false, "Make some noise!")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName("mockhiato") // name of config file (without extension)
	viper.AddConfigPath(".")

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Infof("Configuring mockhiato with %s", viper.ConfigFileUsed())
	}
}

func getConfig(cmd *cobra.Command) lib.Config {
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		panic(err)
	}
	config := lib.Config{}
	if err := viper.Unmarshal(&config); err != nil {
		panic(err)
	}
	if config.Verbose {
		log.SetLevel(log.DebugLevel)
	}
	if !strings.HasSuffix(config.MockFileName, ".go") { // Ensure mock files end with ".go"
		config.MockFileName += ".go"
	}
	log.Debugf("Configs: %#v", config)
	return config
}
