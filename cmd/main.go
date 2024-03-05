// main.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"palworld-query-api/internal/config"
	"palworld-query-api/internal/server"
)

func main() {
	cfg := config.ParseFlags()

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		// Respond with 200 OK status
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received API request")
		serverData, err := server.GetServerData(cfg)
		if err != nil {
			log.Printf("Error getting server data: %v\n", err)
			http.Error(w, "Error getting server data", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(serverData)
		log.Println("Sent server data to client")
	})

	log.Printf("server listening on port %d\n", cfg.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), nil))
}
