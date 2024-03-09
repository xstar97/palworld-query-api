package config

// API constants
var ApiConfig = struct {
	Base string
	Search string
	List string
}{
	Base: "https://api.palworldgame.com",
	Search: "/server/search", // query by name
	List: "/server/list", // paginated list
}
