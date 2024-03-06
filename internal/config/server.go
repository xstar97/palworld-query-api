package config

import (
	"fmt"
	"strings"
)

type ServerInfo struct {
	Name    string  `json:"serverName"`
	Version string  `json:"serverVer"`
	Players Players `json:"players"`
}

type Players struct {
	Count int      `json:"count"`
	List  []string `json:"list"`
}

func GetServerData(serverName string) (*ServerInfo, error) {
	serverInfo := &ServerInfo{}

	infoOutput, err := runRCONCommand(serverName, "info")
	if err != nil {
		return nil, fmt.Errorf("error running 'rcon-cli info': %v", err)
	}
	parseServerInfo(infoOutput, serverInfo)

	playersOutput, err := runRCONCommand(serverName, "showplayers")
	if err != nil {
		return nil, fmt.Errorf("error running 'rcon-cli showplayers': %v", err)
	}
	parsePlayerList(playersOutput, serverInfo)

	return serverInfo, nil
}

func runRCONCommand(serverName string, command string) (string, error) {
    output, err := ExecuteShellCommand(CONFIG.CLI_ROOT, COMMANDS.CONFIG, CONFIG.CLI_CONFIG, COMMANDS.ENV, serverName, command)
    if err != nil {
        return "", fmt.Errorf("failed to run rcon-cli: %v", err)
    }
    return string(output), nil // Convert output to string before returning
}

func parseServerInfo(output string, serverInfo *ServerInfo) {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Welcome to Pal Server") {
			parts := strings.Split(line, "]")
			if len(parts) == 2 {
				serverInfo.Version = strings.TrimSpace(parts[0][strings.Index(parts[0], "[")+1:])
				serverInfo.Name = strings.TrimSpace(parts[1])
				break
			}
		}
	}
}

func parsePlayerList(output string, serverInfo *ServerInfo) {
	lines := strings.Split(output, "\n")
	players := Players{
		List: make([]string, 0), // Initialize the list with an empty slice
	}
	for _, line := range lines {
		if !strings.HasPrefix(line, "name,playeruid,steamid") && line != "" {
			playerData := strings.Split(line, ",")
			if len(playerData) > 0 {
				players.List = append(players.List, playerData[0])
			}
		}
	}
	players.Count = len(players.List)
	serverInfo.Players = players
}
