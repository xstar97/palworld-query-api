package config

import (
	"flag"
	"log"
	"os"
)

// Configuration constants
var Config = struct {
	Port string
    ConfigJson string
	CliConfig string
	LogsPath string
	CachePath string
}{
	Port:         "3000",
    ConfigJson:   "",
	CliConfig:    "/config/rcon.yaml",
	LogsPath:     "/logs",
	CachePath:     "/cache",
}

// setIfNotEmpty sets the value of a string variable if the corresponding environment variable is not empty.
func setIfNotEmpty(key string, value *string) {
	if env := os.Getenv(key); env != "" {
		*value = env
	}
}

// setConfigFromEnv sets configuration values from environment variables.
func setConfigFromEnv() {
	setIfNotEmpty("PORT", &Config.Port)
	setIfNotEmpty("CLI_CONFIG", &Config.CliConfig)
	setIfNotEmpty("CONFIG_JSON", &Config.ConfigJson)
	setIfNotEmpty("LOGS_PATH", &Config.LogsPath)
	setIfNotEmpty("CACHE_PATH", &Config.CachePath)
}

// init parses flags and sets configuration.
func init() {
	// Set configuration from environment variables
	setConfigFromEnv()

	flag.StringVar(&Config.Port, "port", Config.Port, "Server port")
	flag.StringVar(&Config.CliConfig, "cli-config", Config.CliConfig, "path to rcon.yaml")
	flag.StringVar(&Config.ConfigJson, "config-json", Config.ConfigJson, "json object")
	flag.StringVar(&Config.LogsPath, "logs-path", Config.LogsPath, "Logs path")
	flag.StringVar(&Config.CachePath, "cache-path", Config.CachePath, "Cache path")
	flag.Parse()
	// Check if CONFIG_JSON is set
	if Config.ConfigJson != "" {
		// Update the existing config file if it exists, otherwise create a new one
		err := GenerateConfigFromJSON(Config.ConfigJson, Config.CliConfig, Config.LogsPath)
		if err != nil {
			log.Fatalf("Error generating config from JSON: %v", err)
		}
	}

	// Log the set flags
	log.Printf("Server port: %s", Config.Port)
	log.Printf("Root path to rcon.yaml: %s", Config.CliConfig)
	log.Printf("Logs path: %s", Config.LogsPath)
	log.Printf("Cache path: %s", Config.CachePath)
}
