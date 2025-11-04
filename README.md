# Parallel SSH (PSSH)

PSSH is a command-line tool for executing SSH commands and scripts in parallel
across multiple hosts. It's designed for my own use, but feel free to use it as
well.

## Features

- **Parallel Execution**: Run commands concurrently on multiple hosts
- **Script Support**: Execute both inline commands and script files
- **Sudo Support**: Run commands with sudo privileges using stored passwords
- **Host Filtering**: Filter target hosts by tags for selective execution
- **Concurrency Control**: Limit the number of concurrent connections
- **Dry Run Mode**: Preview actions without executing them
- **SSH Agent Integration**: Uses SSH agent for authentication
- **Known Hosts Management**: Automatically manages SSH known hosts

## Installation

### Prerequisites

- Go 1.25.1 or later
- SSH agent running with your keys loaded
- Access to target hosts via SSH key authentication

### Build from Source

```bash
# Clone the repository
git clone https://github.com/txgruppi/pssh.git
cd pssh

# Build the binary
go build -o pssh .

# Or use the justfile
just build
```

## Usage

### Basic Syntax

```bash
pssh --hosts <hosts-file> [options] <command|script-file>
```

### Examples

#### Run a simple command on all hosts

```bash
pssh --hosts hosts.json "uptime"
```

#### Run a command with sudo

```bash
pssh --hosts hosts.json --sudo "systemctl restart nginx"
```

#### Filter hosts by tag

```bash
pssh --hosts hosts.json --tag production "df -h"
```

#### Execute a script file

```bash
pssh --hosts hosts.json ./deploy.sh
```

#### Limit concurrent connections

```bash
pssh --hosts hosts.json --limit 5 "ps aux | grep nginx"
```

#### Dry run to preview actions

```bash
pssh --hosts hosts.json --dry-run --sudo "systemctl restart apache2"
```

#### Filter by multiple tags

```bash
pssh --hosts hosts.json --tag web --tag production "nginx -t"
```

### Command Line Options

- `--hosts` (required): Path to JSON file containing host definitions
- `--tag`: Filter hosts by tag (can be specified multiple times)
- `--sudo`: Execute commands with sudo privileges
- `--dry-run`: Preview actions without executing them
- `--limit`: Limit number of concurrent connections (default: unlimited)

## Hosts File Format

The hosts file is a JSON array containing host definitions:

```json
[
    {
        "name": "web-server-01",
        "addr": "192.168.1.10",
        "user": "deploy",
        "sudo_password": "your_sudo_password",
        "tags": ["web", "production", "frontend"]
    },
    {
        "name": "db-server-01",
        "addr": "192.168.1.20",
        "user": "admin",
        "sudo_password": "admin_sudo_password",
        "tags": ["database", "production", "backend"]
    },
    {
        "name": "dev-server-01",
        "addr": "192.168.1.100",
        "user": "developer",
        "sudo_password": "dev_password",
        "tags": ["development", "testing"]
    }
]
```

### Host Configuration Fields

- `name`: Human-readable identifier for the host
- `addr`: IP address or hostname
- `user`: SSH username for connection
- `sudo_password`: Password for sudo operations (store securely!)
- `tags`: Array of tags for filtering hosts

## Security Considerations

⚠️ **Important**: The hosts file contains sudo passwords in plain text. Ensure
proper file permissions and secure storage:

```bash
# Set restrictive permissions on hosts file
chmod 600 hosts.json
```

## How It Works

1. **Authentication**: Uses SSH agent for key-based authentication
2. **File Transfer**: Uploads scripts to target hosts via SFTP
3. **Execution**: Runs commands/scripts on remote hosts
4. **Output Collection**: Aggregates and prefixes output by hostname
5. **Cleanup**: Removes temporary files from target hosts

---

_This README content is AI-generated and human-reviewed._
