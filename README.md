# palworld-query-api

Streamlined web API for effortlessly managing and querying Palworld game servers using the RCON protocol.

## Features

- Supports querying multiple Palworld game servers.
- Dynamic routing for retrieving server data by name.

## Installation

To install and run *palworld-query-api*, follow these steps:

1. Clone the repository: `git clone https://github.com/xstar97/palworld-query-api.git`
2. Navigate to the project directory: `cd palworld-query-api`
3. Build the project: `go build cmd/main.go`
4. Run the compiled binary: `./palworld-query-api`

Make sure you have Go installed and properly configured on your system before proceeding.

### Command-Line Installation

To install and run *palworld-query-api* from the command line, follow these steps:

1. Clone the repository: `git clone https://github.com/xstar97/palworld-query-api.git`
2. Navigate to the project directory: `cd palworld-query-api`
3. Build the project: `go build`
4. Run the compiled binary: `./palworld-query-api`

Make sure you have Go installed and properly configured on your system before proceeding.

#### Command-Line Flags

You can customize the behavior of *palworld-query-api* using the following command-line flags:

| Flag               | Description                           | Default Value      |
|--------------------|---------------------------------------|--------------------|
| `-port`            | Web port                              | `3000`             |
| `-cli-root`        | Root path to rcon file                | `/app/rcon/rcon`   |
| `-cli-config`      | Root path to rcon.yaml                | `/config/rcon.yaml`|
| `-logs-path`       | Logs path                             | `/logs`            |

Replace the default values as needed when running the binary.

### Docker Installation

Alternatively, you can use the Docker image hosted on GitHub. Use the following `docker-compose.yml` file:

```yaml
version: '3.8'

services:
  palworld-query-api:
    image: ghcr.io/xstar97/palworld-query-api:latest
    environment:
      - PORT=3000
      # default values; really dont need to be changed!
      # - CLI_ROOT=/app/rcon/rcon
      # - CLI_CONFIG=/config/rcon.yaml
      # - LOGS_PATH=/logs
      # generates the yaml from this json array (optional, but recommended)
      # - CONFIG_JSON='{"servers":[{"name":"default","address":"localhost:25575","password":"1234567890","type":"rcon","timeout":"10s"}]}'
    ports:
      - "3000:3000"
    volumes:
      - ./config:/config
      - ./logs:/logs
```

an env variable `CONFIG_JSON` can be set to automatically create the rcon.yaml file needed for the rcon-cli tool.

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

### License

This project is licensed under the MIT License - see the [LICENSE](./LICENSE) file for details.

## Contributing

Contributions are welcome! Please see the [CONTRIBUTING.md](./CONTRIBUTING.md) file for more details.

## Acknowledgements

- [rcon-cli](https://github.com/gorcon/rcon-cli) - The underlying CLI tool for RCON communication.
- [palworld](https://palworld.gg/) - The game server platform supported by this tool.
