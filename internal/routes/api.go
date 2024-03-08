package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"palworld-query-api/internal/config"
	"reflect"
	"sync"
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

var (
	cacheReady  = make(chan struct{})
	updateMutex sync.Mutex
)

func ApiHandler(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	queryParams := r.URL.Query()

	// Wait until the cache is ready
	<-cacheReady

	// Get cached data
	data := config.GetCachedData()

	// Decode the cached JSON data directly into a slice of Server structs
	var servers []Server
	err := json.Unmarshal(data, &servers)
	if err != nil {
		log.Printf("Error decoding cached JSON: %s", err)
		http.Error(w, fmt.Sprintf("Error decoding cached JSON: %s", err), http.StatusInternalServerError)
		return
	}

	// Filter servers based on query parameters
	filteredServers := servers
	for key, values := range queryParams {
		if len(values) > 0 {
			if key == "address" {
				// Check if the value is a domain and try to resolve it to a public IP address
				if config.IsValidDomain(values[0]) {
					publicIP, err := config.GetPublicIP(values[0])
					if err != nil {
						log.Printf("Error resolving domain %s: %s", values[0], err)
						continue
					}
					filteredServers = filterServersByParam(filteredServers, key, publicIP)
				} else {
					filteredServers = filterServersByParam(filteredServers, key, values[0])
				}
			} else {
				filteredServers = filterServersByParam(filteredServers, key, values[0])
			}
		}
	}

	// Convert filtered server list to JSON
	filteredServersJSON, err := json.Marshal(filteredServers)
	if err != nil {
		log.Printf("Error marshalling filtered servers to JSON: %s", err)
		http.Error(w, fmt.Sprintf("Error marshalling filtered servers to JSON: %s", err), http.StatusInternalServerError)
		return
	}

	// If there is only one result, return the object directly
	if len(filteredServers) == 1 {
		filteredServersJSON, err = json.Marshal(filteredServers[0])
		if err != nil {
			log.Printf("Error marshalling filtered server to JSON: %s", err)
			http.Error(w, fmt.Sprintf("Error marshalling filtered server to JSON: %s", err), http.StatusInternalServerError)
			return
		}
	}

	// Set response headers and write JSON content
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(filteredServersJSON)
}

func filterServersByParam(servers []Server, key string, value string) []Server {
	filteredServers := make([]Server, 0)
	log.Printf("Received query parameter: %s, Value: %s", key, value)

	// Iterate over each server
	for _, server := range servers {
		// Get the server struct type
		serverType := reflect.TypeOf(server)
		// Iterate over each field of the server struct
		for i := 0; i < serverType.NumField(); i++ {
			field := serverType.Field(i)
			// Check if the field name matches the query parameter key
			if field.Tag.Get("json") == key {
				// Get the field value
				fieldValue := reflect.ValueOf(server).FieldByName(field.Name).Interface()
				// Convert the field value to a string for comparison
				fieldValueStr := fmt.Sprintf("%v", fieldValue)
				// Compare the field value with the query parameter value
				if fieldValueStr == value {
					filteredServers = append(filteredServers, server)
				}
				break // Move to the next server once a match is found
			}
		}
	}

	log.Printf("Number of servers before filtering: %d", len(servers))
	log.Printf("Number of servers after filtering: %d", len(filteredServers))
	return filteredServers
}

func fetchAndUpdateData() {
	// Lock to prevent concurrent updates
	updateMutex.Lock()
	defer updateMutex.Unlock()

	// Construct the API request URL
	apiUrl := "https://api.palworldgame.com/server/list"

	// Initialize slice to store all servers
	var allServers []Server

	// Start with page 1
	page := 1

	// Loop until there are no more pages
	for {
		// Send a request to the API endpoint for the current page
		log.Printf("Fetching data for page %d...", page)
		response, err := http.Get(fmt.Sprintf("%s?page=%d", apiUrl, page))
		if err != nil {
			log.Printf("Error fetching data for page %d: %s", page, err)
			fmt.Println("Error fetching data:", err)
			return
		}
		defer response.Body.Close()

		// Decode the response JSON
		var serverListResponse ServerListResponse
		err = json.NewDecoder(response.Body).Decode(&serverListResponse)
		if err != nil {
			log.Printf("Error decoding JSON for page %d: %s", page, err)
			fmt.Println("Error decoding JSON:", err)
			return
		}

		// Append servers from the current page
		allServers = append(allServers, serverListResponse.ServerList...)

		// Check if there are more pages
		if !serverListResponse.IsNextPage {
			log.Println("All pages fetched.")
			break
		}

		// Move to the next page
		page++
	}

	// Convert the accumulated server list to JSON
	data, err := json.Marshal(allServers)
	if err != nil {
		log.Printf("Error marshalling servers to JSON: %s", err)
		fmt.Println("Error marshalling servers to JSON:", err)
		return
	}

	// Update the cached data
	config.SetCachedData(data)

	// Signal that the cache is ready
	log.Println("Cache is ready.")
	close(cacheReady)
}

func init() {
	// Start fetching data in the background
	go fetchAndUpdateData()
}
