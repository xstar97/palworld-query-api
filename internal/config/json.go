package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

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