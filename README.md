# Go Walker MCP Server

A simple Go MCP (Model Context Protocol) server that recursively walks directories and returns a list of all files.

## Features

- Exposes a single MCP tool named `walk_directory`
- Recursively lists all files in the specified directory
- Uses "/" as file separators regardless of the operating system
- Remaps "/" to the current working directory when launched
- Built with the official Go MCP library

## Requirements

- Go 1.21 or later
- The official Go MCP library

## Installation

1. Clone or download this repository
2. Navigate to the project directory
3. Install dependencies:

```bash
go mod tidy
```

## Usage

### Building

```bash
go build -o go-walker-mcp main.go
```

### Running

```bash
./go-walker-mcp
```

The server will start and listen for MCP connections. It will walk the directory from which it was launched.

## MCP Tool

### walk_directory

**Description**: Recursively walks the directory and returns a list of all files

**Parameters**:
- `path` (string, optional): Path to walk. Use "/" for the current directory (default: "/")

**Returns**: A list of all files found in the directory tree, with paths using "/" separators.

## Example

When called with `path: "/"`, the tool will:
1. Use the current working directory as the base
2. Recursively walk all subdirectories
3. Return a list of all files found
4. Convert all path separators to "/" regardless of the operating system

## Project Structure

```
go-walker-mcp/
├── main.go          # Main server implementation
├── go.mod           # Go module definition
├── README.md        # This file
└── SPECS.md         # Original specifications
```

## Implementation Details

- The server uses `filepath.Walk` to recursively traverse directories
- All paths are converted to use "/" separators as specified
- The server captures the current working directory at startup
- Paths are returned relative to the base directory
- Error handling is included for invalid paths and I/O errors 