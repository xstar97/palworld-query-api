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
    routRcon := config.Routes.Rcon
    routeApi := config.Routes.Api

	// Register healthz route
	http.HandleFunc(routeHealth, routes.HealthHandler)

	// Register rcon route
	http.HandleFunc(routRcon, routes.RconHandler)

	// Register api route
	http.HandleFunc(routeApi, routes.ApiHandler)

	// Register root route to list available routes
	http.HandleFunc(routeRoot, routes.IndexHandler)

	log.Printf("server listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
