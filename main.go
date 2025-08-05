package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/modelcontextprotocol/go-mcp"
	"github.com/modelcontextprotocol/go-mcp/protocol"
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

func (s *WalkerServer) ListTools(ctx context.Context) ([]protocol.Tool, error) {
	return []protocol.Tool{
		{
			Name:        "walk_directory",
			Description: "Recursively walks the directory and returns a list of all files",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"path": map[string]interface{}{
						"type":        "string",
						"description": "Path to walk (use '/' for current directory)",
						"default":     "/",
					},
				},
			},
		},
	}, nil
}

func (s *WalkerServer) CallTool(ctx context.Context, name string, arguments map[string]interface{}) ([]protocol.Message, error) {
	if name != "walk_directory" {
		return nil, fmt.Errorf("unknown tool: %s", name)
	}

	// Get the path argument, defaulting to "/"
	pathArg, ok := arguments["path"].(string)
	if !ok {
		pathArg = "/"
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

	return []protocol.Message{
		{
			Role:    "tool",
			Content: []protocol.TextContent{{Type: "text", Text: content}},
		},
	}, nil
}

func main() {
	server := NewWalkerServer()

	// Create MCP server
	mcpServer := mcp.NewServer(server)

	// Run the server
	if err := mcpServer.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running MCP server: %v\n", err)
		os.Exit(1)
	}
}
