package config

import (
)

// Constants for routes
var Routes = struct {
	Index string
	Servers string
	Health  string
}{
	Index: "/",
	Servers: "/servers/",
	Health:  "/healthz",
}
var RoutesList = []string{Routes.Servers, Routes.Health}
