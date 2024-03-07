package config

import (
    "strings"
    "log"
    "unicode"
	"unicode/utf8"
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

    infoOutput, err := sendCommand(serverName, PALWORLD_RCON_COMMANDS.INFO)
    if err != nil {
        log.Printf("Error running INFO: %v", err)
    }
	log.Printf("infoOutput: %v", infoOutput)

    // Parse server version and name
    serverInfo.Version = ParseServerVersion(infoOutput)
    serverInfo.Name = ParseServerName(infoOutput)

    playersOutput, err := sendCommand(serverName, PALWORLD_RCON_COMMANDS.SHOW_PLAYERS)
    if err != nil {
        log.Printf("Error running INFO: %v", err)
    }
    // Parse player list
    count, players := ParsePlayerList(playersOutput)
    serverInfo.Players.Count = count
    serverInfo.Players.List = players

    return serverInfo, nil
}

func ParseServerVersion(input string) string {
    parts := strings.Split(input, "[")
    if len(parts) < 2 {
        log.Println("Invalid input format in ParseServerVersion")
        return "" // Invalid input format
    }
    version := strings.TrimSpace(strings.TrimSuffix(strings.Split(parts[1], "]")[0], " "))
    return version
}
func ParseServerName(input string) string {
    // Find the index of "]" and "\n"
    start := strings.Index(input, "]")
    end := strings.Index(input, "\n")

    // Check if both markers are found
    if start == -1 || end == -1 {
        return "" // Invalid input format
    }

    // Extract the name between "]" and "\n"
    name := strings.TrimSpace(input[start+1:end])

    // Remove any null terminators from the name
    name = strings.TrimRightFunc(name, func(r rune) bool {
        return r == '\u0000'
    })

    // Remove any other non-printable characters from the name
    name = strings.Map(func(r rune) rune {
        if r < utf8.RuneSelf && !unicode.IsPrint(r) {
            return -1
        }
        return r
    }, name)

    return name
}
func ParsePlayerList(input string) (int, []Player) {
    var players []Player

    lines := strings.Split(input, "\n")
    count := 0
    for i := 3; i < len(lines); i++ { // Adjusted the start index to skip the separator line
        line := strings.TrimSpace(lines[i])
        if line == "" {
            continue // Skip empty lines
        }
        playerData := strings.Split(line, ",")
        if len(playerData) != 3 {
            log.Printf("Malformed player data in line %d: %s", i, line)
            continue // Skip malformed player data
        }
        player := Player{
            Name: strings.TrimSpace(playerData[0]),
            PID:  strings.TrimSpace(playerData[1]),
            SID:  strings.TrimSpace(playerData[2]),
        }
        players = append(players, player)
        count++
    }

    return count, players
}

func sendCommand(serverName string, command string) (string, error) {
    output, err := ExecuteShellCommand(CONFIG.CLI_ROOT, COMMANDS.CONFIG, CONFIG.CLI_CONFIG, COMMANDS.ENV, serverName, command)
    if err != nil {
        return "", err // Return the error directly
    }
    return string(output), nil // Convert output to string before returning
}












/*
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
}*/



	/*
    infoOutput, err := sendCommand(serverName)
    if err != nil {
        return nil, fmt.Errorf("error running: %v", err)
    }
    parseServerInfo(infoOutput, serverInfo)

    playersOutput, err := sendCommand(serverName)
    if err != nil {
        return nil, fmt.Errorf("error running: %v", err)
    }
    parsePlayerList(playersOutput, serverInfo)*/