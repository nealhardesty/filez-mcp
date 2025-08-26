# filez-mcp - Directory Walker MCP Server

A Model Context Protocol (MCP) server written in Go that provides directory walking functionality. The server exposes a single tool that recursively lists all files and directories from a specified root path, supporting both HTTP and stdio transport methods.

## Features

- **MCP Compliant**: Implements the Model Context Protocol specification
- **Dual Transport Support**: Both HTTP (streamable) and stdio transport methods
- **Security**: Path traversal protection ensuring requests stay within the configured root directory
- **Cross-Platform**: Works on Windows, macOS, and Linux with consistent forward-slash path output
- **Comprehensive Testing**: Unit tests and integration tests for both transport methods

## Installation

### Prerequisites

- Go 1.23.0 or later
- Make (optional, for using the Makefile)

### Building from Source

```bash
# Clone the repository
git clone <repository-url>
cd filez-mcp

# Install dependencies
go mod tidy

# Build the binary
make build
# Or without make:
go build -ldflags="-s -w" -o directory-walker main.go
```

## Usage

### Command Line Interface

```bash
./directory-walker [-s] <root_directory>
```

#### Arguments

- `<root_directory>` (required): An absolute or relative path to the root directory for walking operations
- `-s` (optional): Use stdio transport instead of HTTP transport (default)

#### Examples

```bash
# Run with HTTP transport (default)
./directory-walker /home/user/projects
./directory-walker ./my-project

# Run with stdio transport  
./directory-walker -s .
./directory-walker -s /tmp
```

### Transport Methods

#### HTTP Transport (Default)

The server listens on port 5001 by default (configurable via `PORT` environment variable) and serves the MCP protocol at the `/mcp` endpoint.

```bash
# Start HTTP server
./directory-walker /path/to/directory

# Custom port
PORT=8080 ./directory-walker /path/to/directory
```

Server will be available at: `http://localhost:5001/mcp`

#### Stdio Transport

Use the `-s` flag to enable stdio transport, suitable for direct integration with MCP clients.

```bash
./directory-walker -s /path/to/directory
```

## MCP Tool: walk_directory

The server provides a single MCP tool called `walk_directory` that recursively traverses directory trees.

### Tool Definition

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

### Path Mapping

- `"/"` or empty string: Maps to the server's configured root directory
- `"/subdir"`: Maps to `<root_directory>/subdir`
- Paths are automatically normalized and security-checked

### Response Format

Returns a JSON array of absolute paths using forward slashes as separators:

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

## Development

### Project Structure

```
filez-mcp/
├── main.go              # Main application
├── main_test.go          # Unit tests
├── Makefile              # Build configuration
├── go.mod                # Go module definition
├── go.sum                # Go module checksums
├── README.md             # This file
└── test/
    ├── test-http.go      # HTTP transport test tool
    └── test-stdio.go     # Stdio transport test tool
```

### Available Make Targets

```bash
make build           # Build the directory-walker binary
make run             # Run server in HTTP mode with current directory
make run-stdio       # Run server in stdio mode with current directory
make test            # Run unit tests
make test-coverage   # Run tests with coverage report
make test-http       # Test HTTP transport
make test-stdio      # Test stdio transport
make test-integration# Run all integration tests
make clean           # Remove build artifacts
make fmt             # Format code
make lint            # Run linter
make check           # Run fmt, lint, and test
make deps            # Install/update dependencies
make help            # Show available targets
```

### Running Tests

```bash
# Unit tests
make test

# Integration tests
make test-integration

# Individual transport tests
make test-http
make test-stdio

# Coverage report
make test-coverage
```

### Code Quality

The project follows Go best practices and includes:

- Comprehensive unit tests with race detection
- Integration tests for both transport methods
- Code formatting with `gofmt`
- Linting with `go vet`
- Test coverage reporting

## Dependencies

- [github.com/modelcontextprotocol/go-sdk](https://github.com/modelcontextprotocol/go-sdk) - Official MCP Go SDK
- [github.com/google/jsonschema-go](https://github.com/google/jsonschema-go) - JSON Schema support
- Go standard library packages

## Error Handling

The server handles various error conditions gracefully:

- **Invalid Arguments**: Validates exactly one root directory argument is provided
- **Permission Denied**: Logs errors and continues walking accessible paths
- **Invalid Paths**: Returns appropriate MCP error responses for paths outside the root
- **Non-existent Paths**: Returns errors for paths that don't exist

## Security

- **Path Traversal Protection**: All paths are validated to ensure they remain within the configured root directory
- **Input Validation**: Tool parameters are validated against the JSON schema
- **Error Logging**: Security-relevant errors are logged to stderr without disrupting the MCP protocol

## License

[Add your license information here]

## Contributing

[Add contribution guidelines here]

## Support

For issues and questions:
- Check the [Model Context Protocol documentation](https://modelcontextprotocol.io/)
- Review the [Go SDK documentation](https://github.com/modelcontextprotocol/go-sdk)
- File issues in the project repository
