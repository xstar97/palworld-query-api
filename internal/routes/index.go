package routes

import (
    "net/http"
    "palworld-query-api/internal/config"
    "fmt"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
    // Create clickable links for each route
    availableRoutes := ""
    routes := config.RoutesList
    for _, route := range routes {
        availableRoutes += fmt.Sprintf("<a href=\"%s\">%s</a><br>", route, route)
    }

    // Create HTML content
    htmlContent := fmt.Sprintf("<html><body>%s</body></html>", availableRoutes)

    // Set response headers and write HTML content
    w.Header().Set("Content-Type", "text/html")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(htmlContent))
}
