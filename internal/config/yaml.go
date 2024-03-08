package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"log"
	"gopkg.in/yaml.v2"
	"github.com/fsnotify/fsnotify"
)
type ConfigServer struct {
	Address  string `json:"address"`
	Password string `json:"password"`
	Type     string `json:"type"`
	Timeout  string `json:"timeout"`
}

// Function to read the YAML config file and return the content
func GetConfig() (map[string]ConfigServer, error) {
	filePath := Config.CliConfig

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
					reloadData, err := GetConfig()
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
func GetServerConfig(serverName string) (ConfigServer, error) {
	data, err := GetConfig()
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
