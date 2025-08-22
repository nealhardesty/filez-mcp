# MCP Directory Walker Server - Technical Specification

## Overview
Build a Model Context Protocol (MCP) server in Go that provides directory walking capabilities via HTTP transport. The server exposes a single tool to recursively list all files and directories under a specified path.

## Core Requirements

### Language & Architecture
- **Language**: Go (latest stable version)
- **Architecture**: Single file implementation (`main.go`)
- **Transport**: HTTP MCP server
- **Dependencies**: Official Go MCP library

### Functionality
- **Tool Name**: `walk_directory`
- **Operation**: Recursively walk directory structure starting from configured root
- **Output**: List of full paths for all files and directories found
- **Path Format**: Always use `/` as path separator (cross-platform normalization)
- **Root Mapping**: Map `/` requests to the server's launch directory

### Command Line Interface
```bash
./directory-walker <root_directory>
```

**Parameters:**
- `<root_directory>` (required): Absolute or relative path to the root directory for walking operations

**Example Usage:**
```bash
./directory-walker /home/user/projects
./directory-walker ./my-project
./directory-walker .
```

## Technical Implementation Details

### MCP Tool Specification
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

### Path Handling Rules
1. Input `"/"` maps to the command-line specified root directory
2. Input `"/subdir"` maps to `<root_directory>/subdir`
3. All output paths use `/` separator regardless of OS
4. Return full absolute paths in results
5. Handle symbolic links appropriately (follow or skip based on Go defaults)

### Response Format
Return JSON array of path strings:
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

### Error Handling
- Validate command line arguments (exactly one required)
- Handle permission denied errors gracefully
- Return appropriate MCP error responses for invalid paths
- Log errors to stderr without breaking MCP protocol

## Build System

### Makefile Requirements
Create `Makefile` with the following targets:

```makefile
# Build target - compile the binary
build:
    # Compile to ./directory-walker executable
    # Enable appropriate Go build flags

# Run target - build and run with example directory
run: build
    # Run server with current directory as root
    # Should demonstrate the server starting up

# Clean target (optional but recommended)
clean:
    # Remove built binaries

# Test target (if tests are added)
test:
    # Run go test
```

**Usage Examples:**
```bash
make build
make run
make clean
```

## Project Structure
```
project-root/
├── main.go           # Single file implementation
├── Makefile         # Build configuration
├── go.mod           # Go module definition
├── go.sum           # Go module checksums
└── README.md        # Basic usage instructions
```

## Dependencies
- Official Go MCP library (import path to be determined based on available library)
- Standard library packages:
  - `os`
  - `path/filepath`
  - `strings`
  - `encoding/json`
  - `net/http`
  - `log`

## Implementation Notes

### File Walking Logic
- Use `filepath.Walk()` or `filepath.WalkDir()` for directory traversal
- Convert all paths to use `/` separator using `filepath.ToSlash()`
- Handle both files and directories in the output
- Preserve directory structure in the returned paths

### MCP Server Setup
- Initialize HTTP MCP server on available port
- Register the `walk_directory` tool
- Handle MCP protocol messages correctly
- Provide proper tool metadata and capabilities

### Cross-Platform Considerations
- Normalize path separators to `/` in all outputs
- Handle Windows drive letters appropriately
- Ensure consistent behavior across Unix-like and Windows systems

## Validation Criteria
1. **Functional**: Tool successfully lists all files/directories recursively
2. **Protocol Compliance**: Proper MCP HTTP server implementation
3. **Path Handling**: Correct path normalization and root mapping
4. **Build System**: Makefile targets work as specified
5. **Single File**: Implementation contained in one Go file
6. **Error Resilience**: Graceful handling of filesystem errors

## Success Metrics
- Server starts without errors when given valid root directory
- `walk_directory` tool returns comprehensive file listing
- All paths use `/` separator format
- MCP clients can successfully invoke the tool
- Build and run targets execute successfully