package config

import (
	"flag"
	"log"
	"os"
)

type ConfigServer struct {
	Address  string `json:"address"`
	Password string `json:"password"`
	Type     string `json:"type"`
	Timeout  string `json:"timeout"`
}

type JsonServerConfig struct {
	Name     string `json:"name"`
	Address  string `json:"address"`
	Password string `json:"password"`
	Type     string `json:"type"`
	Timeout  string `json:"timeout"`
}

type JsonConfigData struct {
	Servers []JsonServerConfig `json:"servers"`
}

// Constants for routes
var ROUTES = struct {
    SERVERS       string
    HEALTH        string
}{
    SERVERS:         "/servers/",
    HEALTH:         "/healthz",
}

// Configuration constants
var CONFIG = struct {
    // Web port
    PORT string
    // Root path to rcon file
    CLI_ROOT string
    // Root path to rcon.yaml
    CLI_CONFIG string
    // Default rcon env
    CLI_DEFAULT_SERVER string
    // Default rcon env
    LOGS_PATH string
}{}

// Constants for commands
var COMMANDS = struct {
    ENV     string
    CONFIG  string
}{
    ENV:     "--env",
    CONFIG:  "--config",
}

// Constants for commands
var PALWORLD_RCON_COMMANDS = struct {
    INFO     string
    SHOW_PLAYERS  string
}{
    INFO:          "info",
    SHOW_PLAYERS:  "showplayers",
}

// Function to set configuration from environment variables
func setConfigFromEnv() {
    setIfNotEmpty := func(key string, value *string) {
        if env := os.Getenv(key); env != "" {
            *value = env
        }
    }

    setIfNotEmpty("PORT", &CONFIG.PORT)
    setIfNotEmpty("CLI_ROOT", &CONFIG.CLI_ROOT)
    setIfNotEmpty("CLI_CONFIG", &CONFIG.CLI_CONFIG)
    setIfNotEmpty("CLI_DEFAULT_SERVER", &CONFIG.CLI_DEFAULT_SERVER)
    setIfNotEmpty("LOGS_PATH", &CONFIG.LOGS_PATH)
}

// Parse flags
func init() {
    // Set configuration from environment variables
    setConfigFromEnv()

    flag.StringVar(&CONFIG.PORT, "port", "3000", "Server port")
    flag.StringVar(&CONFIG.CLI_ROOT, "cli-root", "/app/rcon/rcon", "Root path to rcon file")
    flag.StringVar(&CONFIG.CLI_CONFIG, "cli-config", "/config/rcon.yaml", "Root path to rcon.yaml")
    flag.StringVar(&CONFIG.CLI_DEFAULT_SERVER, "cli-def-server", "default", "Default rcon env")
    flag.StringVar(&CONFIG.LOGS_PATH, "logs-path", "/logs", "Logs path")
    flag.Parse()

	// Check if CONFIG_JSON is set
	configJSON := os.Getenv("CONFIG_JSON")
	if configJSON != "" {
		// Update the existing config file if it exists, otherwise create a new one
		err := GenerateConfigFromJSON(configJSON, CONFIG.CLI_CONFIG, CONFIG.LOGS_PATH)
		if err != nil {
			log.Fatalf("Error generating config from JSON: %v", err)
		}
	}
    // Log the set flags
    log.Printf("Server port: %s", CONFIG.PORT)
    log.Printf("Root path to rcon file: %s", CONFIG.CLI_ROOT)
    log.Printf("Root path to rcon.yaml: %s", CONFIG.CLI_CONFIG)
    log.Printf("Default rcon env: %s", CONFIG.CLI_DEFAULT_SERVER)
    log.Printf("Logs path: %s", CONFIG.LOGS_PATH)
}
