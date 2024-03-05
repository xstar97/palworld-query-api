// internal/config/config.go
package config

import (
	"flag"
	"log"
	"os"
	"palworld-query-api/internal/utils"
	"strconv"
)

type Config struct {
	RconCLIPath   string
	RconCLIConfig string
	Port          int
	LogsPath      string // New field for logs directory
}

func ParseFlags() *Config {
	var config Config

	// Check environment variables
	config.RconCLIPath = os.Getenv("RCON_CLI_PATH")
	config.RconCLIConfig = os.Getenv("RCON_CLI_CONFIG")
	config.LogsPath = os.Getenv("LOGS_PATH") // Read the LOGS_PATH environment variable
	portEnv := os.Getenv("PORT")
	if portEnv != "" {
		p, err := strconv.Atoi(portEnv)
		if err != nil {
			log.Fatalf("Invalid value for PORT environment variable: %v", err)
		}
		config.Port = p
	}

	// Parse flags if environment variables not set
	flag.StringVar(&config.RconCLIPath, "rcon-cli-path", "/app/rcon/rcon", "Path to the rcon-cli executable")
	flag.StringVar(&config.RconCLIConfig, "rcon-cli-config", "/config/rcon.yaml", "Path to the rcon-cli config file")
	flag.IntVar(&config.Port, "port", 3000, "server port")
	flag.StringVar(&config.LogsPath, "logs-path", "/logs", "Path to the directory for log files") // Add logs-path flag
	flag.Parse()

	// Check if RconCLIPath exists
	if _, err := os.Stat(config.RconCLIPath); os.IsNotExist(err) {
		log.Fatalf("RconCLIPath '%s' does not exist", config.RconCLIPath)
	}

	// Check if CONFIG_JSON is set
	configJSON := os.Getenv("CONFIG_JSON")
	if configJSON != "" {
		// Update the existing config file if it exists, otherwise create a new one
		err := utils.GenerateConfigFromJSON(configJSON, config.RconCLIConfig, config.LogsPath)
		if err != nil {
			log.Fatalf("Error generating config from JSON: %v", err)
		}
	}

	return &config
}
