# ContainerShip

A simple CLI tool for managing Docker containers with ease. ContainerShip provides a straightforward interface to ship, stop, monitor, and interact with your containerized applications.

## Features

- **Ship Containers**: Deploy your applications in containers quickly
- **Container Management**: Start, stop, and monitor container status
- **Logs & Debugging**: View container logs and execute commands inside containers
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

ContainerShip uses `containership.yaml` for configuration. See the example file for details.

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