package config

import (
    "strings"
    "log"
    "unicode"
	"unicode/utf8"
)

type ServerInfo struct {
    Online    bool    `json:"online"`
    Name      string  `json:"serverName"`
    Version   string  `json:"serverVer"`
    Players   Players `json:"players"`
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
    serverInfo.Name = ParseServerName(serverInfo.Version, infoOutput)

    // Set Online field based on server name and version
    if serverInfo.Name != "" && serverInfo.Version != "" {
        serverInfo.Online = true
    } else {
        serverInfo.Online = false
    }

    playersOutput, err := sendCommand(serverName, PALWORLD_RCON_COMMANDS.SHOW_PLAYERS)
    if err != nil {
        log.Printf("Error running INFO: %v", err)
    }
    log.Printf("playersOutput: %v", playersOutput)
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

func ParseServerName(version, input string) string {
    // Define the prefix pattern
    prefix := "Welcome to Pal Server[" + version + "] "

    // Check if the last character of the input string is a newline
    if !strings.HasSuffix(input, "\n") {
        // Append a newline character to the input string
        input += "\n"
    }

    // Find the index of the prefix
    prefixIndex := strings.Index(input, prefix)
    if prefixIndex == -1 {
        return "" // Prefix not found
    }

    // Find the index of the newline character after the prefix
    newlineIndex := strings.Index(input[prefixIndex:], "\n")
    if newlineIndex == -1 {
        return "" // Newline not found
    }

    // Extract the name between the prefix and the newline character
    name := input[prefixIndex+len(prefix) : prefixIndex+newlineIndex]

    // Trim any leading or trailing whitespace
    name = strings.TrimSpace(name)

    // Remove any null characters (\u0000) from the name
    name = strings.Map(func(r rune) rune {
        if r == '\u0000' {
            return -1
        }
        return r
    }, name)

    // Remove any non-printable characters from the name
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
	for i := 1; i < len(lines); i++ { // Adjusted the start index to skip the header line
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

	if count == 0 {
		return count, []Player{} // Return an empty array if count is 0
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
