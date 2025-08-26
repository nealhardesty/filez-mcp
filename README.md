# MCP Directory Walker Server

A Model Context Protocol (MCP) server in Go that provides directory walking capabilities via HTTP transport. The server exposes a single tool to recursively list all files and directories under a specified path.

## Features

- **Single Tool**: `walk_directory` - recursively walks directory structure
- **HTTP Transport**: MCP server with HTTP/REST endpoints
- **Path Normalization**: Cross-platform path handling with `/` separators
- **Error Handling**: Graceful handling of permission errors and invalid paths
- **CLI Interface**: Simple command-line interface for server startup

## Installation

### Prerequisites
- Go 1.23 or later
- Make (optional, for build automation)

### Build

```bash
# Using Make
make build

# Or directly with Go
go build -o directory-walker main.go
```

## Usage

### Starting the Server

```bash
# Start with current directory as root
./directory-walker .

# Start with specific directory as root
./directory-walker /home/user/projects
./directory-walker ./my-project
```

### MCP Tool Usage

The server exposes one tool: `walk_directory`

**Tool Schema:**
```json
{
  "name": "walk_directory",
  "description": "Recursively lists all files and directories under the specified path",
  "inputSchema": {
    "type": "object",
    "properties": {
      "path": {
        "type": "string",
        "description": "Directory path to walk (use '/' for root directory)",
        "default": "/"
      }
    }
  }
}
```

**Path Mapping:**
- Input `"/"` → maps to command-line specified root directory
- Input `"/subdir"` → maps to `<root_directory>/subdir`
- All output paths use `/` separator (cross-platform)

**Example Response:**
```json
{
  "content": [
    "/full/path/to/file1.txt",
    "/full/path/to/subdir",
    "/full/path/to/subdir/file2.go",
    "/full/path/to/another/deep/file.json"
  ]
}
```

### Server Endpoints

When running, the server provides:
- **MCP Endpoint**: `http://localhost:8080/mcp`
- **Protocol**: JSON-RPC 2.0 over HTTP POST

## Development

### Build Commands

```bash
make build    # Build binary
make run      # Build and run with current directory
make clean    # Remove built binaries
make deps     # Install/update dependencies
make test     # Run tests
make fmt      # Format code
make lint     # Lint code (requires golangci-lint)
```

### Testing

Run the CLI test suite:

```bash
./scripts/test-cli.sh
```

This script tests:
- Server startup with various directory arguments
- Error handling for invalid directories and missing arguments
- Basic functionality verification

### Project Structure

```
project-root/
├── main.go              # Single file implementation
├── Makefile            # Build configuration
├── go.mod              # Go module definition
├── go.sum              # Go module checksums
├── README.md           # This file
└── scripts/
    └── test-cli.sh     # CLI testing script
```

## Dependencies

- **github.com/metoro-io/mcp-golang**: MCP protocol implementation
- **github.com/gin-gonic/gin**: HTTP web framework (used by MCP library)
- Standard library: `os`, `path/filepath`, `strings`, `fmt`, `log`

## Error Handling

The server handles various error conditions:
- **Invalid root directory**: Server exits with error message
- **Permission denied**: Logged but directory walking continues
- **Invalid MCP requests**: Returns appropriate MCP error responses
- **Missing arguments**: Shows usage information and exits

## Cross-Platform Support

- Path separators automatically normalized to `/` in all outputs
- Handles Windows drive letters appropriately
- Consistent behavior across Unix-like and Windows systems

## License

This project follows the specifications defined in `specs/SPEC-1.md`.