# ContainerShip

A simple CLI tool for managing Docker containers with ease. ContainerShip provides a straightforward interface to ship, stop, monitor, and interact with your containerized applications.

## Features

- **Ship Containers**: Deploy your applications in containers quickly
- **Container Management**: Start, stop, restart, and monitor container status
- **Logs & Debugging**: View container logs and execute commands inside containers
- **Volume Support**: Bind mounts and named volumes
- **Environment Management**: Environment variables and env files
- **Health Checks**: Define health checks for services
- **Restart Policies**: Configure container restart behavior
- **Dependency Management**: Handle service dependencies
- **Simple CLI**: Easy-to-use command-line interface
- **Docker Integration**: Built on top of Docker API

## Installation

### Prerequisites

- Go 1.24 or later
- Docker installed and running

### Build from Source

```bash
git clone https://github.com/Femn0X/ContainerShip.git
cd ContainerShip
make build
```

The binary will be available at `bin/cs`.

### Install Globally

```bash
make install
```

This will install the binary to your GOPATH/bin.

## Usage

### Basic Commands

```bash
# Ship a container (deploy application)
cs ship

# Stop running containers
cs stop

# Restart all services
cs restart

# Check status of containers
cs status

# List defined containers
cs list

# View container logs
cs logs

# Execute command in container
cs exec <container_name> <command>
```

### Help

```bash
cs help
```

## Configuration

ContainerShip uses `containership.yaml` for configuration. Here's an example:

```yaml
version: "1.0"
services:
  web:
    image: nginx:latest
    depends_on: [db]
    ports:
      - "8080:80"
    volumes:
      - ./html:/usr/share/nginx/html
    restart: always
    command: ["nginx", "-g", "daemon off;"]
    env_file:
      - .env
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s
  db:
    image: postgres:15
    environment:
      POSTGRES_PASSWORD: defNotApi
    volumes:
      - db_data:/var/lib/postgresql/data
    restart: unless-stopped
```

### Service Configuration Options

- `image`: Docker image to use
- `depends_on`: List of services this service depends on
- `ports`: Port mappings (host:container)
- `volumes`: Volume mounts (host:container or volume:container)
- `restart`: Restart policy (no, always, unless-stopped, on-failure)
- `command`: Override default command
- `env_file`: List of environment files to load
- `healthcheck`: Health check configuration

## Development

### Build

```bash
make build
```

### Test

```bash
make test
```

### Clean

```bash
make clean
```

### Run

```bash
make run
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for version history.