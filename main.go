package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type WalkerServer struct {
	baseDir string
}

func NewWalkerServer() *WalkerServer {
	// Get the current working directory
	wd, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("Failed to get working directory: %v", err))
	}

	return &WalkerServer{
		baseDir: wd,
	}
}

func (s *WalkerServer) walkDirectoryHandler(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[map[string]interface{}]) (*mcp.CallToolResultFor[any], error) {
	// Get the path argument, defaulting to "/"
	pathArg := "/"
	if params.Arguments != nil {
		if path, ok := params.Arguments["path"].(string); ok {
			pathArg = path
		}
	}

	// Remap "/" to the current directory
	if pathArg == "/" {
		pathArg = s.baseDir
	} else {
		// For other paths, join with base directory
		pathArg = filepath.Join(s.baseDir, pathArg)
	}

	// Walk the directory
	var files []string
	err := filepath.Walk(pathArg, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the base directory itself
		if path == pathArg {
			return nil
		}

		// Convert to relative path from base directory
		relPath, err := filepath.Rel(s.baseDir, path)
		if err != nil {
			return err
		}

		// Convert Windows path separators to "/" as specified
		relPath = strings.ReplaceAll(relPath, "\\", "/")

		files = append(files, relPath)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk directory: %v", err)
	}

	// Create the response message
	content := fmt.Sprintf("Found %d files:\n%s", len(files), strings.Join(files, "\n"))

	return &mcp.CallToolResultFor[any]{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: content,
			},
		},
	}, nil
}

func main() {
	server := NewWalkerServer()

	// Create MCP implementation
	impl := &mcp.Implementation{
		Name:    "go-walker-mcp",
		Title:   "Go Walker MCP Server",
		Version: "1.0.0",
	}

	// Create MCP server
	mcpServer := mcp.NewServer(impl, &mcp.ServerOptions{})

	// Create input schema for the tool
	inputSchema := &jsonschema.Schema{
		Type: "object",
		Properties: map[string]*jsonschema.Schema{
			"path": {
				Type:        "string",
				Description: "Path to walk (use '/' for current directory)",
			},
		},
	}

	// Add the walk_directory tool with proper handler
	mcpServer.AddTool(&mcp.Tool{
		Name:        "walk_directory",
		Description: "Recursively walks the directory and returns a list of all files",
		InputSchema: inputSchema,
	}, server.walkDirectoryHandler)

	// Create SSE handler
	sseHandler := mcp.NewSSEHandler(func(request *http.Request) *mcp.Server {
		return mcpServer
	})

	// Set up HTTP server
	http.HandleFunc("/mcp", sseHandler.ServeHTTP)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8124"
	}

	fmt.Printf("Starting MCP server on port %s\n", port)
	fmt.Printf("SSE endpoint available at: http://localhost:%s/mcp\n", port)

	// Start HTTP server
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting HTTP server: %v\n", err)
		os.Exit(1)
	}
}
