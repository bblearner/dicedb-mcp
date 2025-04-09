# DiceDB MCP

A [Model Context Protocol (MCP)](https://modelcontextprotocol.io) server implementation for DiceDB to enable interactions between AI applications (hosts/clients) and DiceDB database servers.

This implementation uses the [DiceDB Go SDK](https://github.com/DiceDB/dicedb-go) to communicate with DiceDB.

## Features

- PING a DiceDB server to check connectivity.
- ECHO a message through a DiceDB server.

## Installation

Prerequisites:

- Go 1.24 or higher

```bash
go install github.com/pottekkat/dicedb-mcp@latest
```

## Usage

### With MCP Hosts/Clients

Add this to your `claude_desktop_config.json` for Claude Desktop or `mcp.json` for Cursor:

```json
{
    "mcpServers": {
        "dicedb-mcp": {
            "command": "/Users/pottekkat/Git/dicedb-mcp/dist/dicedb-mcp"
        }
    }
}
```

## Available Tools

### ping

Pings a DiceDB server to check connectivity.

### echo

Echoes a message through the DiceDB server.

## Development

Fork and clone the repository:

```bash
git clone https://github.com/username/dicedb-mcp.git
```

Change into the directory:

```bash
cd dicedb-mcp
```

Install dependencies:

```bash
make deps
```

Build the project:

```bash
make build
```

Update your MCP servers configuration to point to the local build:

```json
{
    "mcpServers": {
        "dicedb-mcp": {
            "command": "/path/to/dicedb-mcp/dist/dicedb-mcp"
        }
    }
}
```

## License

[MIT License](LICENSE)
