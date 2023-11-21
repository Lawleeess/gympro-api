package config

import (
	"log"
	"sync"

	"github.com/spf13/viper"
)

var (
	onceConfigin sync.Once
	CfgIn        Config
)

// Config stores all configuration of the application
// The values are ready by viper from a config file or env variables
type Config struct {
	// HTTP server configuration
	ServerPort string `mapstructure:"SERVER_PORT"`

	GoogleProjectID string `mapstructure:"GOOGLE_PROJECT_ID"`

	// Firebase REST API
	FirebaseHost string `mapstructure:"FIREBASE_HOST"`
	FirebaseKey  string `mapstructure:"FIREBASE_KEY"`

	// Google Cloud credential
	ServiceCredentialJSON string `mapstructure:"SERVICE_CREDENTIAL_JSON"`
}

// LoadConfig reads configuration from file or environment variables
func LoadConfig() {
	onceConfigin.Do(func() {
		// load config path when is called from handlers test
		viper.AddConfigPath("../../../../")
		viper.AddConfigPath("../../../")
		// load cinfig path when is called from main
		viper.AddConfigPath(".")
		viper.SetConfigName("app")
		viper.SetConfigType("env")

		// viper.AutomaticEnv allows overrides from environment variables
		// If a value from configuration is already defined by env var then
		// It will ignore the value in env file
		viper.AutomaticEnv()

		if err := viper.ReadInConfig(); err != nil {
			log.Fatal("failed to load config: " + err.Error())
		}

		if err := viper.Unmarshal(&CfgIn); err != nil {
			log.Fatal("failed to unmarshal configuration: " + err.Error())
		}
	})
}

func init() {
	LoadConfig()
}
