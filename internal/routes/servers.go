package routes

import (
    "encoding/json"
    "log"
    "net/http"
    "palworld-query-api/internal/config"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
    servers, err := config.ReadConfig()
    if err != nil {
        http.Error(w, "Failed to read server configurations", http.StatusInternalServerError)
        log.Println("Failed to read server configurations:", err)
        return
    }

    log.Printf("Received API request: %s\n", r.URL.Path)

    path := r.URL.Path
    
    // Check if the path is the servers route
    if path == config.ROUTES.SERVERS {
        // Index route for all servers - List all server data keyed by server names
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
    serverName := path[len(config.ROUTES.SERVERS):]
    if serverName == "" {
        // If no server name provided, return bad request
        http.Error(w, "Invalid server name", http.StatusBadRequest)
        log.Println("Invalid server name")
        return
    }

    // Get server data by name
    serverDataInfo, err := config.GetServerData(serverName)
    if err != nil {
        log.Printf("Error getting server data for %s: %v\n", serverName, err)
        http.Error(w, "Error getting server data", http.StatusInternalServerError)
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
        serverDataInfo, err := config.GetServerData(name)
        if err != nil {
            return nil, err
        }
        serverDataMap[name] = serverDataInfo
    }
    return serverDataMap, nil
}
