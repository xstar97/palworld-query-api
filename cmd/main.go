package main

import (
    "fmt"
    "log"
    "net/http"
    "palworld-query-api/internal/config"
    "palworld-query-api/internal/routes"
)

func main() {
    port := fmt.Sprintf(":%s", config.CONFIG.PORT)

    // Register healthz route
    http.HandleFunc(config.ROUTES.HEALTH, routes.HealthzHandler)

    // Register servers route
    http.HandleFunc(config.ROUTES.SERVERS, routes.IndexHandler) 
    
    log.Printf("server listening on port %s\n", port)
    log.Fatal(http.ListenAndServe(port, nil))
}
