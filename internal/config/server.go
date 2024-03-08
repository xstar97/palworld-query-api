package config

import (
    "strings"
    "log"
    "unicode"
	"unicode/utf8"
	"github.com/gorcon/rcon"
    "fmt"
    "errors"
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

func GetServerData(configServer ConfigServer) (*ServerInfo, error) {
    serverInfo := &ServerInfo{}
    cmdInfo := Rcon.Command.Info
    cmdShowPlayers := Rcon.Command.ShowPlayers

    infoCommandOutput, err := sendCommand(configServer, cmdInfo)
    if err != nil {
        log.Printf("Error running INFO: %v", err)
    }
    log.Printf("infoCommandOutput: %v", infoCommandOutput)

    // Parse server version and name
    serverInfo.Version = ParseServerVersion(infoCommandOutput)
    serverInfo.Name = ParseServerName(serverInfo.Version, infoCommandOutput)

    // Set Online field based on server name and version
    if serverInfo.Name != "" && serverInfo.Version != "" {
        serverInfo.Online = true
    } else {
        serverInfo.Online = false
    }

    playersCommandOutput, err := sendCommand(configServer, cmdShowPlayers)
    if err != nil {
        log.Printf("Error running INFO: %v", err)
    }
    log.Printf("playersCommandOutput: %v", playersCommandOutput)

    // Parse player list
    count, players := ParsePlayerList(playersCommandOutput)
    serverInfo.Players.Count = count
    serverInfo.Players.List = players
    return serverInfo, nil
}

func ParseServerVersion(input string) string {
    if input == "" {
        log.Println("Input is null or empty in ParseServerVersion")
        return ""
    }
    
    parts := strings.Split(input, "[")
    if len(parts) < 2 {
        log.Println("Invalid input format in ParseServerVersion")
        return "" // Invalid input format
    }
    version := strings.TrimSpace(strings.TrimSuffix(strings.Split(parts[1], "]")[0], " "))
    return version
}

func ParseServerName(version, input string) string {
    // Check if the input is null or empty
    if input == "" {
        log.Println("Input is null or empty in ParseServerName")
        return ""
    }

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

func sendCommand(configServer ConfigServer, command string) (string, error) {
	if configServer.Address == "" {
		return "", errors.New("RCON server address is empty")
	}
	if configServer.Password == "" {
		return "", errors.New("RCON server password is empty")
	}

	conn, err := rcon.Dial(configServer.Address, configServer.Password, rcon.SetDialTimeout(configServer.Timeout), rcon.SetDeadline(configServer.Timeout))
	if err != nil {
		log.Println("Error connecting to RCON server:", err)
		return "", err
	}
	defer conn.Close()

	response, err := conn.Execute(command)
	if err != nil {
		log.Println("Error executing command:", err)
		return "", err
	}

	fmt.Println(response)
	return string(response), nil // Convert output to string before returning
}
