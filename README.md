# palworld-query-api

Streamlined web API for effortlessly managing and querying PalWorld game servers using the RCON protocol.

## Features

- Supports querying multiple PalWorld game servers.
- Dynamic routing for retrieving server data by name.

## Installation

To install and run *palworld-query-api*, follow these steps:

1. Clone the repository: `git clone https://github.com/xstar97/palworld-query-api.git`
2. Navigate to the project directory: `cd palworld-query-api`
3. Build the project: `go build cmd/main.go`
4. Run the compiled binary: `./palworld-query-api`

Make sure you have Go installed and properly configured on your system before proceeding.

### Command-Line Flags

You can customize the behavior of *palworld-query-api* using the following command-line flags:

| Flag               | Description                           | Default Value      |
|--------------------|---------------------------------------|--------------------|
| `-port`            | Web port                              | `3000`             |
| `-cli-config`      | Root path to rcon.yaml                | `/config/rcon.yaml`|
| `-logs-path`       | Logs path                             | `/logs`            |

Replace the default values as needed when running the binary.

### Routes

- `/healthz`: This route is used to check the health status of the server.

- `/rcon/:name`: This route is used to retrieve server information by specifying the server name.

- `/rcon/`: This route lists all available servers and their information.
- /api: The PalWorld server list api
  - accepts query params
  - requires a ?name param to search by server name.
  - additional params can further filter the list.

#### API Route Params

| Query Key       | Example                | Required |
|-----------------|------------------------|----------|
| name            | ?name=MyServer         | &#10003; |
| server_id       | ?server_id=123         | &#9744;  |
| namespace       | ?namespace=example     | &#9744;  |
| type            | ?type=public           | &#9744;  |
| region          | ?region=US             | &#9744;  |
| map_name        | ?map_name=Map1         | &#9744;  |
| description     | ?description=Description | &#9744;  |
| address         | ?address=example.com   | &#9744;  |
| port            | ?port=8080             | &#9744;  |
| is_password     | ?is_password=true      | &#9744;  |
| version         | ?version=1.0           | &#9744;  |
| created_at      | ?created_at=1234567890 | &#9744;  |
| update_at       | ?update_at=1234567890  | &#9744;  |
| world_guid      | ?world_guid=123        | &#9744;  |
| current_players | ?current_players=5     | &#9744;  |
| max_players     | ?max_players=10        | &#9744;  |
| days            | ?days=30               | &#9744;  |
| server_time     | ?server_time=1234567890| &#9744;  |

### Docker Installation

Alternatively, you can use the Docker image hosted on GitHub. Use the following `docker-compose.yml` file:

```yaml
version: '3.8'

services:
  palworld-query-api:
    image: ghcr.io/xstar97/palworld-query-api:latest
    environment:
      - PORT=3000
      # generates the yaml from this json array (optional, but recommended)
      # - CONFIG_JSON='{"servers":[{"name":"default","address":"localhost:25575","password":"1234567890","type":"rcon","timeout":"10s"}]}'
    ports:
      - "3000:3000"
    volumes:
      - ./config:/config
      - ./logs:/logs
```

an env variable `CONFIG_JSON` can be set to automatically create the rcon.yaml file needed for the rcon-cli dependency.

```json
{
  "servers": [
    {
      "name": "default",
      "address": "localhost:25575",
      "password": "1234567890",
      "type": "rcon",
      "timeout": "60s"
    }
  ]
}
```

### HomePage Integration CustomAPI API route

Integrate PalWorld server information seamlessly into your homepage using the CustomAPI widget. By specifying the server environment name, you can display key details such as server name, version, and current player count. Keep your users informed with real-time updates on server status.

the output of this route /api/?name=SERVER_NAME

```json
{
  "server_id": "Server ID",
  "namespace": "production",
  "type": "community",
  "region": "all",
  "name": "Server Name",
  "map_name": "Mainworld5",
  "description": "My Awesome Server",
  "address": "XXX.XXX.XX.XXX",
  "port": 8211,
  "is_password": true,
  "version": "v0.1.5.1",
  "created_at": 050505050505,
  "update_at": 060606060,
  "world_guid": "WORLD_GUID",
  "current_players": 9000, # OVER 9000!
  "max_players": 32,
  "days": 7366,
  "server_time": 37
}
```

The custom widget for homepage:

```yaml
    - PalWorld:
        icon: https://tech.palworldgame.com/img/logo.jpg
        description: A clone PKM game
        widget:
          type: customapi
          url: "http://localhost:3000/api?name=SERVER_NAME!"
          refreshInterval: 10000
          method: GET
          mappings:
            - field: name
              label: Name
              format: text
            - field: version
              label: Version
              format: text
            - field: current_players
              label: Current Players
              format: number
            - field: days
              label: Days Count
              format: text
```

### HomePage Integration CustomAPI Rcon route

Integrate PalWorld server information seamlessly into your homepage using the CustomAPI widget. By specifying the server environment name, you can display key details such as server name, version, and current player count. Keep your users informed with real-time updates on server status.

the output of this route /rcon/:name

```json
{
  "online": false,
  "serverName": "",
  "serverVer": "",
  "players": {
    "count": 0,
    "list": []
  }
}
```

The custom widget for homepage:

```yaml
    - PalWorld:
        icon: https://tech.palworldgame.com/img/logo.jpg
        description: A clone PKM game
        widget:
          type: customapi
          url: "http://localhost:3000/rcon/default" # change the name given to server env (not palworld server name!)
          refreshInterval: 10000
          method: GET
          mappings:
            - field: serverName
              label: Name
              format: text
            - field: serverVer
              label: Version
              format: text
            - field:
                players: count
              label: Current Players
              format: number
            - field: online
              label: Status
              format: text
              remap:
                - value: false
                  to: Not Online
                - value: true
                  to: Online
```

### License

This project is licensed under the MIT License - see the [LICENSE](./LICENSE) file for details.

## Contributing

Contributions are welcome! Please see the [CONTRIBUTING.md](./CONTRIBUTING.md) file for more details.

## Acknowledgements

- [rcon](https://github.com/gorcon/rcon) - The underlying RCON communication.
- [palworld](https://palworld.gg/) - The game server platform supported by this tool.
