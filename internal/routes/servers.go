package routes

import (
    "encoding/json"
    "log"
    "net/http"
    "palworld-query-api/internal/config"
)

func ServersHandler(w http.ResponseWriter, r *http.Request) {
    routeServers := config.Routes.Servers
    path := r.URL.Path
    servers, err := config.GetConfig()

    if path == routeServers {
        if err != nil {
            http.Error(w, "Failed to read server configurations", http.StatusInternalServerError)
            log.Println("Failed to read server configurations:", err)
            return
        }

        log.Printf("Received API request: %s\n", path)

        serverDataMap, err := getAllServerData(servers)
        if err != nil {
            http.Error(w, "Error getting all server data", http.StatusInternalServerError)
            log.Println("Error getting all server data:", err)
            return
        }

        // Encode and send the response
        w.Header().Set("Content-Type", "application/json")
        if err := json.NewEncoder(w).Encode(serverDataMap); err != nil {
            http.Error(w, "Error encoding server data", http.StatusInternalServerError)
            log.Println("Error encoding server data:", err)
            return
        }
        log.Println("Sent all server data to client")
        return
    }

    // Extract server name from path
    serverName := path[len(routeServers):]

    // Validate if the serverName exists in the servers map
    serverData, ok := servers[serverName]
    if !ok {
        // Server does not exist
        emptyResponse := map[string]interface{}{"message": "Server does not exist"}
        w.Header().Set("Content-Type", "application/json")
        if err := json.NewEncoder(w).Encode(emptyResponse); err != nil {
            http.Error(w, "Error encoding server data", http.StatusInternalServerError)
            log.Println("Error encoding server data:", err)
            return
        }
        log.Printf("Server %s does not exist", serverName)
        return
    }

    // Get server data by name
    serverDataInfo, err := config.GetServerData(serverData)
    if err != nil {
        log.Printf("Error getting server data for %s: %v\n", serverName, err)
        // Return empty JSON object indicating that the server does not exist
        emptyResponse := map[string]interface{}{"message": "Server data retrieval error"}
        w.Header().Set("Content-Type", "application/json")
        if err := json.NewEncoder(w).Encode(emptyResponse); err != nil {
            http.Error(w, "Error encoding server data", http.StatusInternalServerError)
            log.Println("Error encoding server data:", err)
        }
        return
    }

    // Encode and send the response
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(serverDataInfo); err != nil {
        http.Error(w, "Error encoding server data", http.StatusInternalServerError)
        log.Println("Error encoding server data:", err)
        return
    }
    log.Printf("Sent server data for %s to client", serverName)
}

func getAllServerData(servers map[string]config.ConfigServer) (map[string]interface{}, error) {
    serverDataMap := make(map[string]interface{})
    for name := range servers {
        serverDataInfo, err := config.GetServerData(servers[name])
        if err != nil {
            return nil, err
        }
        serverDataMap[name] = serverDataInfo
    }
    return serverDataMap, nil
}
