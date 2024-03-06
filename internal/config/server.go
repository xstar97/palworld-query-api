package config

import (
    "fmt"
    "strings"
    "log"
)

type ServerInfo struct {
    Name    string  `json:"serverName"`
    Version string  `json:"serverVer"`
    Players Players `json:"players"`
}

type Player struct {
    Name string `json:"name"`
    PID  string `json:"pid"`
    SID  string `json:"sid"`
}

type Players struct {
    Count int      `json:"count"`
    List  []Player `json:"list"`
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
            log.Printf("Parsing line: %s\n", line)
            parts := strings.SplitN(line, "]", 2)
            if len(parts) == 2 {
                version := strings.TrimSpace(parts[0][strings.Index(parts[0], "[")+1:])
                serverName := strings.TrimSpace(parts[1])
                log.Printf("Extracted Version: %s, Server Name: %s\n", version, serverName)
                serverInfo.Version = version
                serverInfo.Name = serverName
                log.Printf("ServerInfo after parsing: %+v\n", serverInfo)
                break
            }
        }
    }
}

func parsePlayerList(output string, serverInfo *ServerInfo) {
    lines := strings.Split(output, "\n")
    players := Players{
        List: make([]Player, 0), // Initialize the list with an empty slice
    }
    for _, line := range lines {
        if !strings.HasPrefix(line, "name,playeruid,steamid") && line != "" {
            log.Printf("Parsing player list line: %s\n", line)
            playerData := strings.Split(line, ",")
            if len(playerData) >= 3 {
                playerName := playerData[0]
                playerID := playerData[1]
                steamID := playerData[2]
                log.Printf("Extracted player name: %s, player ID: %s, steam ID: %s\n", playerName, playerID, steamID)
                player := Player{
                    Name: playerName,
                    PID:  playerID,
                    SID:  steamID,
                }
                players.List = append(players.List, player)
            }
        }
    }
    players.Count = len(players.List)
    serverInfo.Players = players
    log.Printf("Player list parsed. Total players: %d\n", players.Count)
}
