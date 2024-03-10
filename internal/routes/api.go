package routes

import (
    "encoding/json"
    "fmt"
    "html/template"
    "log"
    "net/http"
    "net/url"
    "palworld-query-api/internal/config"
    "reflect"
    "strings"
)

type Server struct {
    ServerID            string `json:"server_id"`
    Namespace           string `json:"namespace"`
    Type                string `json:"type"`
    Region              string `json:"region"`
    Name                string `json:"name"`
    MapName             string `json:"map_name"`
    Description         string `json:"description"`
    Address             string `json:"address"`
    Port                int    `json:"port"`
    IsPassword          bool   `json:"is_password"`
    Version             string `json:"version"`
    CreatedAt           int64  `json:"created_at"`
    UpdateAt            int64  `json:"update_at"`
    WorldGUID           string `json:"world_guid"`
    CurrentPlayers      int    `json:"current_players"`
    MaxPlayers          int    `json:"max_players"`
    Days                int    `json:"days"`
    ServerTime          int    `json:"server_time"`
}

type ServerListResponse struct {
    CurrentPage int      `json:"current_page"`
    PageSize    int      `json:"page_size"`
    SortType    string   `json:"sort_type"`
    ServerType  string   `json:"server_type"`
    Region      string   `json:"region"`
    IsNextPage  bool     `json:"is_next_page"`
    ServerList  []Server `json:"server_list"`
    NextPageURL string   `json:"next_page_url"`
}

func ApiHandler(w http.ResponseWriter, r *http.Request) {
    // Parse query parameters
    queryParams := r.URL.Query()

    // Check if the "name" query parameter is provided
    nameQuery := queryParams.Get("name")
    if nameQuery == "" {
        http.Error(w, "Name query parameter is required", http.StatusBadRequest)
        return
    }

    // Set the default search field to "q" if not provided
    if queryParams.Get("q") == "" {
        queryParams.Set("q", nameQuery)
    }

    // Construct the initial search URL
    searchURL := config.ApiConfig.Base + config.ApiConfig.Search + "?q=" + url.QueryEscape(queryParams.Get("q"))

    var allServers []Server

    // Pagination loop
    for {
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

        // Append servers to the result
        allServers = append(allServers, serverListResponse.ServerList...)

        // Check if there are more pages
        if !serverListResponse.IsNextPage {
            break
        }

        // Update the search URL for the next page
        searchURL = config.ApiConfig.Base + serverListResponse.NextPageURL
    }

    // Filter servers based on query parameters other than "q"
    filteredServers := allServers
    for key, values := range queryParams {
        // Skip filtering if the query parameter is "q" or "name"
        if key == "q" || key == "name" {
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

    // Return the result
    if len(filteredServers) == 1 {
        // If only one server found after filtering, return it as a single object
        serverJSON, err := json.Marshal(filteredServers[0])
        if err != nil {
            log.Printf("Error marshalling server to JSON: %s", err)
            http.Error(w, fmt.Sprintf("Error marshalling server to JSON: %s", err), http.StatusInternalServerError)
            return
        }
        contentType := "application/json"
        if acceptsHTML(r) {
            contentType = "text/html"
            // Render HTML instead of JSON
            renderHTML(w, filteredServers[0])
        } else {
            w.Header().Set("Content-Type", contentType)
            w.WriteHeader(http.StatusOK)
            w.Write(serverJSON)
        }
    } else if len(filteredServers) > 1 {
        // If multiple servers found after filtering, return them as an array
        contentType := "application/json"
        if acceptsHTML(r) {
            contentType = "text/html"
            // Render HTML instead of JSON
            renderHTMLList(w, filteredServers)
        } else {
            // Return JSON
            serverListJSON, err := json.Marshal(filteredServers)
            if err != nil {
                log.Printf("Error marshalling server list to JSON: %s", err)
                http.Error(w, fmt.Sprintf("Error marshalling server list to JSON: %s", err), http.StatusInternalServerError)
                return
            }
            w.Header().Set("Content-Type", contentType)
            w.WriteHeader(http.StatusOK)
            w.Write(serverListJSON)
        }
    } else {
        // No servers found after filtering
        log.Println("No servers found after filtering.")
        http.Error(w, "No servers found after filtering.", http.StatusNotFound)
    }
}

// Function to render HTML for a single server
func renderHTML(w http.ResponseWriter, server Server) {
    // Define your HTML template for rendering a single server
    htmlTemplate := config.HtmlDetailsTemplate

    // Parse the HTML template
    tmpl, err := template.New("htmlTemplate").Parse(htmlTemplate)
    if err != nil {
        http.Error(w, fmt.Sprintf("Error parsing HTML template: %s", err), http.StatusInternalServerError)
        return
    }

    // Execute the HTML template with the server data
    err = tmpl.Execute(w, server)
    if err != nil {
        http.Error(w, fmt.Sprintf("Error executing HTML template: %s", err), http.StatusInternalServerError)
        return
    }
}

// Function to render HTML for an array of servers
func renderHTMLList(w http.ResponseWriter, servers []Server) {
    // Define your HTML template for rendering a list
    htmlTemplate := config.HtmlListTemplate

    // Parse the HTML template
    tmpl, err := template.New("htmlTemplate").Parse(htmlTemplate)
    if err != nil {
        http.Error(w, fmt.Sprintf("Error parsing HTML template: %s", err), http.StatusInternalServerError)
        return
    }

    // Execute the HTML template with the array of servers
    err = tmpl.Execute(w, servers)
    if err != nil {
        http.Error(w, fmt.Sprintf("Error executing HTML template: %s", err), http.StatusInternalServerError)
        return
    }
}

// Function to check if the request accepts HTML
func acceptsHTML(r *http.Request) bool {
    accept := r.Header.Get("Accept")
    return strings.Contains(accept, "text/html")
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