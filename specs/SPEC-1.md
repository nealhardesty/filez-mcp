# MCP Directory Walker Server

## Overview

This document specifies the requirements for a Model Context Protocol (MCP) server written in Go. The server will provide directory walking functionality, exposing a single tool that recursively lists all files and directories from a specified root path. It will support both HTTP and standard I/O (stdio) transport methods.

-----

## Core Requirements

### Language and Architecture

  * **Language**: Go (latest stable version)
  * **MCP Library**: Use https://github.com/modelcontextprotocol/go-sdk (the official MCP SDK).  Make sure you read the current implementation details in full (https://github.com/modelcontextprotocol/go-sdk).
  * **Architecture**: A single-file implementation (`main.go`).
  * **Transport**: The server must support both **HTTP** (on the `/mcp` path) and **stdio** transport. A command-line flag (`-s`) will switch between the two, with HTTP being the default.
  * **Dependencies**: The official Go MCP library is required.

-----

## Functional Specification

### `walk_directory` Tool

  * **Name**: `walk_directory`
  * **Operation**: Recursively traverses the directory tree starting from a configured root.
  * **Input**: The tool accepts an optional `path` string.
      * An input of `"/"` maps to the server's launch directory, which is the command-line specified root.
      * An input of `"/subdir"` maps to `<root_directory>/subdir`.
  * **Output**: A JSON array of full, absolute paths for all discovered files and directories. All output paths must use `/` as the path separator, regardless of the operating system.

\<br\>

**Tool Definition:**

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

-----

**Example Tool Response:**

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

-----

## Command Line Interface (CLI)

### Main Application

The application will be run with the following command-line signature:

```bash
./directory-walker [-s] <root_directory>
```

  * `<root_directory>` (**required**): An absolute or relative path to the root directory for walking operations.
  * `-s` (**optional**): If specified, the application will use the stdio MCP implementation instead of the default HTTP one.

\<br\>

**Example Usage:**

  * `./directory-walker /home/user/projects`
  * `./directory-walker ./my-project`
  * `./directory-walker -s .`

### Testing Tools

The build process must also create two command-line testing tools in the `/test` directory:

1.  `test-stdio`: A Go tool to test the MCP server via stdio.
2.  `test-http`: A Go tool to test the MCP server via HTTP, assuming it's running locally on the default port.

-----

## Build System

A `Makefile` will manage the build process with the following targets:

  * **`build`**: Compiles the `main.go` file into a `./directory-walker` executable, ensuring appropriate Go build flags are enabled.
  * **`run`**: Builds the executable and runs the server using the current directory as the root. This target should demonstrate the server starting successfully.
  * **`test`**: Runs all unit tests.
  * **`clean`**: Removes the built binaries.

-----

## Technical Details

### Path Handling

  * Input paths are mapped relative to the command-line specified root directory.
  * All output paths must use the forward slash (`/`) as a separator, irrespective of the operating system.
  * The application should return **full absolute paths** in its results.
  * Symbolic links should be handled appropriately, following Go's default behavior.

### Error Handling

  * Validate that exactly one `root_directory` argument is provided.
  * Handle `permission denied` errors gracefully within the `walk_directory` logic.
  * Return appropriate MCP error responses for invalid paths.
  * Log errors to `stderr` without disrupting the MCP protocol stream.

-----

## Project Structure

```
project-root/
├── main.go               # Single file implementation
├── main_test.go          # Unit tests for the main application
├── Makefile              # Build configuration
├── go.mod                # Go module definition
├── go.sum                # Go module checksums
├── README.md             # Setup and usage instructions (always keep this up to date with every action)
└── test/
    ├── test-http.go      # HTTP test tool
    └── test-stdio.go     # Stdio test tool
```

-----

## Dependencies

  * Official Go MCP library
  * Standard library packages: `os`, `path/filepath`, `strings`, `encoding/json`, `net/http`, `log`

-----

## Implementation Notes

  * **File Walking**: Use `filepath.Walk()` or `filepath.WalkDir()` for efficient directory traversal. Convert paths to use `/` separators with `filepath.ToSlash()`.
  * **MCP Server Setup**: Initialize the HTTP MCP server to listen on an available port (default 5001, but overrideable via the `PORT` environment variable). The stdio implementation should be enabled when the `-s` flag is present.
  * **Cross-Platform Considerations**: Pay close attention to path normalization, especially for Windows drive letters, to ensure consistent behavior across all operating systems.

-----

## Validation Criteria

1.  **Functionality**: The `walk_directory` tool correctly lists all files and directories.
2.  **Protocol Compliance**: The server adheres to the MCP HTTP and stdio protocols.
3.  **Path Handling**: Paths are normalized and mapped correctly.
4.  **Build System**: All `Makefile` targets function as specified.
5.  **Code Structure**: The core logic is contained within a single Go file (`main.go`).
6.  **Error Resilience**: The application handles filesystem and argument errors gracefully.