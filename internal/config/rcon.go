package config

type Command struct {
	Info        string
	ShowPlayers string
}

var Rcon = struct {
	Command Command
}{
	Command: Command{
		Info:        "info",
		ShowPlayers: "showplayers",
	},
}