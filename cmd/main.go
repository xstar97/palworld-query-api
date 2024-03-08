package main

import (
	"fmt"
	"log"
	"net/http"
	"palworld-query-api/internal/config"
	"palworld-query-api/internal/routes"
)

func main() {
	port := fmt.Sprintf(":%s", config.Config.Port)
	routeRoot := config.Routes.Index
    routeHealth := config.Routes.Health
    routeServers := config.Routes.Servers

	// Register healthz route
	http.HandleFunc(routeHealth, routes.HealthHandler)

	// Register servers route
	http.HandleFunc(routeServers, routes.ServersHandler)

	// Register root route to list available routes
	http.HandleFunc(routeRoot, routes.IndexHandler)

	log.Printf("server listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
