package config

import (
)

// Constants for routes
var Routes = struct {
	Index string
	Rcon string
	Api string
	Health  string
}{
	Index: "/",
	Rcon: "/rcon/",
	Api: "/api",
	Health:  "/healthz",
}
var RoutesList = []string{Routes.Health,Routes.Rcon, Routes.Api}
