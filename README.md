# MCP Directory Walker Server

A Model Context Protocol (MCP) server implementation in Go that provides directory walking functionality. The server exposes a single tool that recursively lists all files and directories from a specified root path, supporting both HTTP and stdio transport methods.

## Features

- **Single Tool**: `walk_directory` - Recursively lists all files and directories
- **Dual Transport**: Supports both HTTP and stdio transport protocols
- **Security**: Path validation to prevent directory traversal attacks
- **Cross-Platform**: Consistent forward-slash path separators across all operating systems
- **Error Handling**: Graceful handling of permission errors and invalid paths

## Installation

### Prerequisites

- Go 1.23 or higher
- Git

### Build from Source

1. Clone the repository:
```bash
git clone <repository-url>
cd filez-mcp
```

2. Install dependencies:
```bash
make deps
```

3. Build the binary:
```bash
make build
```

This creates the `directory-walker` executable in the current directory.

## Usage

### Command Line Interface

```bash
./directory-walker [-s] <root_directory>
```

**Arguments:**
- `<root_directory>` (required): Absolute or relative path to the root directory for walking operations
- `-s` (optional): Use stdio transport instead of HTTP (default: HTTP)

**Examples:**
```bash
# Start HTTP server for /home/user/projects
./directory-walker /home/user/projects

# Start HTTP server for current directory
./directory-walker .

# Start stdio server for current directory
./directory-walker -s .
```

### HTTP Transport (Default)

The HTTP server listens on port 5001 by default (configurable via `PORT` environment variable) and serves the MCP protocol on the `/mcp` endpoint.

**Starting the server:**
```bash
./directory-walker /path/to/directory
```

The server will be available at: `http://localhost:5001/mcp`

**Custom port:**
```bash
PORT=8080 ./directory-walker /path/to/directory
```

### Stdio Transport

For stdio transport, use the `-s` flag:

```bash
./directory-walker -s /path/to/directory
```

## Tool Reference

### `walk_directory`

Recursively lists all files and directories under the specified path.

**Input Schema:**
```json
{
  "type": "object",
  "properties": {
    "path": {
      "type": "string",
      "description": "Directory path to walk (use '/' for root directory)",
      "default": "/"
    }
  }
}
```

**Path Mapping:**
- Input `"/"` → Maps to the server's configured root directory
- Input `"/subdir"` → Maps to `<root_directory>/subdir`
- All paths are validated to ensure they stay within the root directory

**Example Responses:**

Success response:
```json
{
  "content": [],
  "structuredContent": [
    "/full/path/to/file1.txt",
    "/full/path/to/subdir",
    "/full/path/to/subdir/file2.go",
    "/full/path/to/another/deep/file.json"
  ]
}
```

Error response:
```json
{
  "content": [
    {
      "type": "text",
      "text": "path does not exist: /nonexistent"
    }
  ],
  "isError": true
}
```

## Development

### Build System

The project uses a Makefile for common development tasks:

```bash
# Build the main binary
make build

# Build and run HTTP server for current directory
make run

# Build and run stdio server for current directory  
make run-stdio

# Run unit tests
make test

# Run integration tests
make test-integration

# Build test tools
make build-test-tools

# Clean build artifacts
make clean

# Show all available targets
make help
```

### Testing

#### Unit Tests

Run unit tests for the core functionality:

```bash
make test
```

#### Integration Tests

##### Stdio Testing

Test the stdio transport:

```bash
make test-stdio
```

This builds the server and test tool, then runs a comprehensive test of the stdio interface.

##### HTTP Testing

1. Start the server in one terminal:
```bash
make run
```

2. Run HTTP tests in another terminal:
```bash
make test-http
```

#### Manual Testing

You can also test manually using the built test tools:

```bash
# Build test tools
make build-test-tools

# Test stdio (replace with your binary path and root directory)
./test/test-stdio ./directory-walker .

# Test HTTP (assumes server is running on localhost:5001)
./test/test-http
```

### Project Structure

```
filez-mcp/
├── main.go               # Single file implementation
├── main_test.go          # Unit tests for the main application
├── Makefile              # Build configuration
├── go.mod                # Go module definition
├── go.sum                # Go module checksums
├── README.md             # This file
├── specs/
│   ├── SPEC-1.md         # Original specification
│   └── TAO-1.md          # Implementation documentation
└── test/
    ├── test-http.go      # HTTP test tool
    └── test-stdio.go     # Stdio test tool
```

## Technical Details

### Path Handling

- Input paths are mapped relative to the command-line specified root directory
- All output paths use forward slash (`/`) separators for cross-platform consistency
- The application returns full absolute paths in results
- Symbolic links are handled using Go's default behavior
- Security validation prevents access outside the root directory

### Error Handling

- Validates exactly one `root_directory` argument is provided
- Handles permission denied errors gracefully within the `walk_directory` logic
- Returns appropriate MCP error responses for invalid paths
- Logs errors to stderr without disrupting the MCP protocol stream

### MCP Compliance

- Implements MCP protocol version 2024-11-05
- Supports both stdio and HTTP transports as specified
- Follows JSON-RPC 2.0 protocol for all communications
- Proper error handling and response formatting

## Dependencies

- [github.com/mark3labs/mcp-go](https://github.com/mark3labs/mcp-go) v0.38.0 - Official Go MCP library
- Go standard library packages: `os`, `path/filepath`, `strings`, `net/http`, `log`, `context`, `flag`, `fmt`, `strconv`

## License

[Add your license information here]

## Contributing

[Add contribution guidelines here]
