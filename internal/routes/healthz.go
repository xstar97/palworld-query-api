package routes

import (
    "encoding/json"
    "net/http"
)

func HealthzHandler(w http.ResponseWriter, r *http.Request) {
    // Set response Content-Type header to indicate JSON
    w.Header().Set("Content-Type", "application/json")
    
    // Set status code
    w.WriteHeader(http.StatusOK)
    
    // Define a response JSON object
    jsonResponse := map[string]string{"status": "ok"}

    // Encode the JSON response object into the response body
    if err := json.NewEncoder(w).Encode(jsonResponse); err != nil {
        // If encoding fails, respond with an internal server error
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}
