// internal/server/server.go
package server

import (
	"fmt"
	"os/exec"
	"strings"
	"palworld-query-api/internal/config"
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

func GetServerData(config *config.Config) (*ServerInfo, error) {
	serverInfo := &ServerInfo{}

	infoOutput, err := runRCONCommand(config, "info")
	if err != nil {
		return nil, fmt.Errorf("error running 'rcon-cli info': %v", err)
	}
	parseServerInfo(infoOutput, serverInfo)

	playersOutput, err := runRCONCommand(config, "showplayers")
	if err != nil {
		return nil, fmt.Errorf("error running 'rcon-cli showplayers': %v", err)
	}
	parsePlayerList(playersOutput, serverInfo)

	return serverInfo, nil
}

func runRCONCommand(config *config.Config, command string) (string, error) {
	cmd := exec.Command(config.RconCLIPath, "-config", config.RconCLIConfig, command)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
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
	players := Players{}
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
