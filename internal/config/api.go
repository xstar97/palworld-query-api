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

var HtmlDetailsTemplate =
`
<html>
<head>
    <style>
        body {
            font-family: Arial, sans-serif;
            background-color: #222;
            color: #fff;
            display: flex;
            justify-content: center;
            align-items: center;
            height: 100vh;
        }
        .server-details {
            padding: 20px;
            border: 1px solid #ccc;
            border-radius: 5px;
            background-color: #333;
            position: relative; /* Add relative positioning */
        }
        .server-details h2 {
            color: #ddd;
        }
        .server-details p {
            margin: 5px 0;
            color: #ddd;
        }
        .players {
            font-weight: bold;
        }
        .version-info-container {
            position: absolute;
            top: 5px;
            right: 5px;
        }
        .version-info {
            padding: 5px; /* Add padding */
            margin: 0; /* Add margin */
        }
        .other-details {
            margin-top: 30px; /* Add margin to separate from version-info */
        }
        .lock-icon {
            width: 20px;
            height: 20px;
            fill: #ddd;
            margin-right: 5px;
        }
        table {
            width: 100%;
            border-collapse: collapse;
            margin-top: 10px; /* Add margin-top for table */
        }
        th, td {
            padding: 8px;
            text-align: left;
            border-bottom: 1px solid #ddd;
        }
    </style>
</head>
<body>
<div class="server-details" style="border-color: {{if .IsPassword}}red{{else}}green{{end}}">
    <div class="version-info-container">
        <p class="version-info">{{.Version}} {{if .IsPassword}}<img class="lock-icon" src="https://emojicdn.elk.sh/🔒?style=google" alt="Has Password"/>{{else}}<img class="lock-icon" src="https://emojicdn.elk.sh/🔓?style=google" alt="No Password"/>{{end}}</p>
    </div>
    <div class="other-details">
	<table>
    {{if ne .Name ""}}
    <tr>
        <th>Name</th>
        <td>{{.Name}}</td>
    </tr>
    {{end}}
    {{if ne .Description ""}}
    <tr>
        <th>Description</th>
        <td>{{.Description}}</td>
    </tr>
    {{end}}
    <tr>
        <th>Address</th>
        <td>{{.Address}}:{{.Port}}</td>
    </tr>
    <tr>
        <th>Players</th>
        <td>{{.CurrentPlayers}}/{{.MaxPlayers}}</td>
    </tr>
    <tr>
        <th>Days</th>
        <td>{{.Days}}</td>
    </tr>
    <tr>
        <th>Map Name</th>
        <td>{{.MapName}}</td>
    </tr>
    <tr>
        <th>Type</th>
        <td>{{.Type}}</td>
    </tr>
    <tr>
        <th>Region</th>
        <td>{{.Region}}</td>
    </tr>
    <tr>
        <th>Created At</th>
        <td>{{.CreatedAt}}</td>
    </tr>
    <tr>
        <th>Update At</th>
        <td>{{.UpdateAt}}</td>
    </tr>
    <tr>
        <th>Namespace</th>
        <td>{{.Namespace}}</td>
    </tr>
    <tr>
        <th>Server Time</th>
        <td>{{.ServerTime}}</td>
    </tr>
    <tr>
        <th>World GUID</th>
        <td>{{.WorldGUID}}</td>
    </tr>
    <tr>
        <th>Server ID</th>
        <td>{{.ServerID}}</td>
    </tr>
</table>
    </div>
</div>
</body>
</html>
`

var HtmlListTemplate = 
`
<html>
<head>
    <style>
        body {
            font-family: Arial, sans-serif;
            background-color: #222;
            color: #fff;
        }
        .server-list {
            display: grid;
            grid-template-columns: repeat(4, minmax(200px, 1fr)); /* Adjusted to have 4 items per row */
            grid-gap: 20px;
            padding: 0;
            margin: 0;
            list-style-type: none;
        }
        .server-item:hover {
            background-color: #555; /* Change background color on hover */
        }
        .server-item {
            padding: 20px;
            border: 1px solid #ccc;
            border-radius: 5px;
            background-color: #333;
            display: flex;
            flex-direction: column;
            align-items: center;
            justify-content: center;
            text-align: center;
            position: relative; /* Add relative positioning */
        }
        .server-item h3 {
            color: #ddd;
            margin-top: 0;
            margin-bottom: 10px; /* Add bottom margin */
            padding: 5px; /* Add padding */
        }
        .server-item p {
            margin: 10px 0; /* Add margin top and bottom */
            padding: 0 10px; /* Add padding left and right */
            color: #ddd;
        }
        .version-info {
            position: absolute;
            top: 5px;
            right: 5px; /* Adjusted position for top right */
            padding: 5px; /* Add padding */
            margin: 0; /* Add margin */
            z-index: 1; /* Ensure it's above other content */
        }
		.lock-icon {
			width: 20px;
			height: 20px;
			fill: #ddd;
			margin-right: 5px;
		}
        .players {
            font-weight: bold;
        }
    </style>
</head>
<body>
    <ul class="server-list">
        {{range .}}
        <div class="clickable-item" onclick="location.href='?name={{.Name}}&server_id={{.ServerID}}';">
            <li class="server-item" style="border-color: {{if .IsPassword}}red{{else}}green{{end}}">
                <p class="version-info">{{.Version}} {{if .IsPassword}}<img class="lock-icon" src="https://emojicdn.elk.sh/🔒?style=google" alt="Has Password"/>{{else}}<img class="lock-icon" src="https://emojicdn.elk.sh/🔓?style=google" alt="No Password"/>{{end}}</p>
                {{if ne .Name ""}}
                <h4>{{.Name}}</h4>
                {{end}}
                {{if ne .Description ""}}
                <p>{{.Description}}</p>
                {{end}}
                <p>{{.Address}}:{{.Port}}</p>
                <p class="players">Players: {{.CurrentPlayers}}/{{.MaxPlayers}}</p>
            </li>
        </div>
        {{end}}
    </ul>
</body>
</html>
`