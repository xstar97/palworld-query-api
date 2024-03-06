package routes

import (
    "net/http"
)

func HealthzHandler(w http.ResponseWriter, r *http.Request) {
    // Respond with 200 OK status
    w.WriteHeader(http.StatusOK)
}
