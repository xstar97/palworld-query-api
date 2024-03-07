package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"log"
	"os/exec"
	"gopkg.in/yaml.v2"
	"github.com/fsnotify/fsnotify"
)

func GenerateConfigFromJSON(configJSON string, outputPath string, logs string) error {
	var JsonConfigData JsonConfigData
	err := json.Unmarshal([]byte(configJSON), &JsonConfigData)
	if err != nil {
		return fmt.Errorf("error parsing JSON: %v", err)
	}

	// Prepare YAML content
	yamlContent := ""
	for _, server := range JsonConfigData.Servers {
		yamlContent += fmt.Sprintf("%s:\n", server.Name)
		yamlContent += fmt.Sprintf("  address: \"%s\"\n", server.Address)
		yamlContent += fmt.Sprintf("  password: \"%s\"\n", server.Password)
		yamlContent += fmt.Sprintf("  log: \"%s/%s.log\"\n", logs, server.Name)
		yamlContent += fmt.Sprintf("  type: \"%s\"\n", server.Type)
		yamlContent += fmt.Sprintf("  timeout: \"%s\"\n", server.Timeout)
	}

	// Write YAML content to file
	err = ioutil.WriteFile(outputPath, []byte(yamlContent), 0644)
	if err != nil {
		return fmt.Errorf("error writing YAML to file: %v", err)
	}

	fmt.Println("Config file generated successfully")
	return nil
}

func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

// Function to read the YAML config file and return the content
func ReadConfig() (map[string]ConfigServer, error) {
	filePath := CONFIG.CLI_CONFIG

	// Log the file path
	log.Printf("Reading config from file: %s", filePath)

	// Check if the file exists
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return nil, err
	}

	// Read YAML file
	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Unmarshal YAML into map
	var data map[string]ConfigServer
	err = yaml.Unmarshal(yamlFile, &data)
	if err != nil {
		return nil, err
	}

	// Create a new watcher to monitor changes to the file
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	defer watcher.Close()

	// Watch for changes to the file
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("Config file modified. Reloading config...")
					// Reload config from file
					reloadData, err := ReadConfig()
					if err != nil {
						log.Println("Error reloading config:", err)
						continue
					}
					log.Println("Config reloaded successfully")
					// Update the existing data with the reloaded data
					data = reloadData
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("Error watching file:", err)
			}
		}
	}()

	// Add the file to the watcher
	err = watcher.Add(filePath)
	if err != nil {
		return nil, err
	}

	log.Println("Config read successfully from file.")
	return data, nil
}

// reads the YAML config file and returns the configuration for a specific server
func GetServer(serverName string) (ConfigServer, error) {
	data, err := ReadConfig()
	if err != nil {
		return ConfigServer{}, err
	}

	// Check if the server name exists
	config, ok := data[serverName]
	if !ok {
		return ConfigServer{}, fmt.Errorf("server '%s' not found", serverName)
	}

	return config, nil
}

// ExecuteShellCommand executes a shell command with provided arguments and returns its output
func ExecuteShellCommand(command string, args ...string) ([]byte, error) {
    // Set the command to execute
    cmd := exec.Command(command, args...)
    
    // Capture the output of the command
    output, err := cmd.CombinedOutput()
    if err != nil {
        log.Println("Error executing command:", err)
    }
    
	log.Println("output: ", string(output))
    return output, nil // Always return output and nil error
}