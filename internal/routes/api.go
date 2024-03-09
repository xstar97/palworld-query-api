package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"palworld-query-api/internal/config"
	"reflect"
)

type Server struct {
	ServerID       string `json:"server_id"`
	Namespace      string `json:"namespace"`
	Type           string `json:"type"`
	Region         string `json:"region"`
	Name           string `json:"name"`
	MapName        string `json:"map_name"`
	Description    string `json:"description"`
	Address        string `json:"address"`
	Port           int    `json:"port"`
	IsPassword     bool   `json:"is_password"`
	Version        string `json:"version"`
	CreatedAt      int    `json:"created_at"`
	UpdateAt       int    `json:"update_at"`
	WorldGUID      string `json:"world_guid"`
	CurrentPlayers int    `json:"current_players"`
	MaxPlayers     int    `json:"max_players"`
	Days           int    `json:"days"`
	ServerTime     int    `json:"server_time"`
}

type ServerListResponse struct {
	CurrentPage int      `json:"current_page"`
	PageSize    int      `json:"page_size"`
	SortType    string   `json:"sort_type"`
	ServerType  string   `json:"server_type"`
	Region      string   `json:"region"`
	IsNextPage  bool     `json:"is_next_page"`
	ServerList  []Server `json:"server_list"`
}

func ApiHandler(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	queryParams := r.URL.Query()

	// Construct the search URL
	searchURL := config.ApiConfig.Base + config.ApiConfig.Search + "?q=" + url.QueryEscape(queryParams.Get("q"))

	// Send a request to the search endpoint
	response, err := http.Get(searchURL)
	if err != nil {
		log.Printf("Error searching for server: %s", err)
		http.Error(w, fmt.Sprintf("Error searching for server: %s", err), http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	// Decode the response JSON
	var serverListResponse ServerListResponse
	err = json.NewDecoder(response.Body).Decode(&serverListResponse)
	if err != nil {
		log.Printf("Error decoding search response: %s", err)
		http.Error(w, fmt.Sprintf("Error decoding search response: %s", err), http.StatusInternalServerError)
		return
	}

	// If no servers found, return empty response
	if len(serverListResponse.ServerList) == 0 {
		log.Println("No servers found.")
		http.Error(w, "No servers found.", http.StatusNotFound)
		return
	}

	// If there is only one result, return the object directly
	if len(serverListResponse.ServerList) == 1 {
		serverJSON, err := json.Marshal(serverListResponse.ServerList[0])
		if err != nil {
			log.Printf("Error marshalling server to JSON: %s", err)
			http.Error(w, fmt.Sprintf("Error marshalling server to JSON: %s", err), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(serverJSON)
		return
	}

	// Filter servers based on query parameters other than "q"
	filteredServers := serverListResponse.ServerList
	for key, values := range queryParams {
		// Skip filtering if the query parameter is "q"
		if key == "q" {
			continue
		}
		if len(values) > 0 {
			// Special handling for "address" query parameter
			if key == "address" {
				// Check if the value is a domain and try to resolve it to a public IP address
				if config.IsValidDomain(values[0]) {
					publicIP, err := config.GetPublicIP(values[0])
					if err != nil {
						log.Printf("Error resolving domain %s: %s", values[0], err)
						continue
					}
					filteredServers = filterServersByParamByKey(filteredServers, key, publicIP)
				} else {
					// If not a domain, filter servers by address directly
					filteredServers = filterServersByParamByKey(filteredServers, key, values[0])
				}
			} else {
				// For other query parameters, filter servers by matching the key with struct field tags
				filteredServers = filterServersByParamByKey(filteredServers, key, values[0])
			}
		}
	}

	// If there is only one result after filtering, return the object directly
	if len(filteredServers) == 1 {
		serverJSON, err := json.Marshal(filteredServers[0])
		if err != nil {
			log.Printf("Error marshalling server to JSON: %s", err)
			http.Error(w, fmt.Sprintf("Error marshalling server to JSON: %s", err), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(serverJSON)
		return
	}

	// More than one server found after filtering, return an error
	log.Println("More than one server found after filtering. Please provide a more specific query.")
	http.Error(w, "More than one server found after filtering. Please provide a more specific query.", http.StatusBadRequest)
}

// Define a new function to filter servers by matching the query parameter key with struct field tags
func filterServersByParamByKey(servers []Server, key string, value string) []Server {
	filteredServers := make([]Server, 0)
	log.Printf("Received query parameter:\nquery: %s\nValue: %s\n", key, value)

	// Iterate over each server
	for _, server := range servers {
		// Get the server struct type
		serverType := reflect.TypeOf(server)
		// Iterate over each field of the server struct
		for i := 0; i < serverType.NumField(); i++ {
			field := serverType.Field(i)
			// Check if the field tag matches the query parameter key
			if field.Tag.Get("json") == key {
				// Get the field value
				fieldValue := reflect.ValueOf(server).FieldByName(field.Name).Interface()
				// Convert the field value to a string for comparison
				fieldValueStr := fmt.Sprintf("%v", fieldValue)
				// Log the server name and the value being compared
				log.Printf("\nServer Name: %s\nField: %s,\nField Value: %s\nQuery Value: %s\n", server.Name, key, fieldValueStr, value)
				// Compare the field value with the query parameter value
				if fieldValueStr == value {
					filteredServers = append(filteredServers, server)
				}
				break // Move to the next server once a match is found
			}
		}
	}

	log.Printf("\nNumber of servers before filtering: %d\n", len(servers))
	log.Printf("Number of servers after filtering: %d\n", len(filteredServers))
	return filteredServers
}