package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Application configuration
type Config struct {
	RootDirectory string
	UseStdio      bool
	Port          string
}

// walkDirectory implements the walk_directory tool
func walkDirectory(rootDir, relativePath string) ([]string, error) {
	// Map relative path to absolute path
	var targetPath string
	if relativePath == "/" || relativePath == "" {
		targetPath = rootDir
	} else {
		// Remove leading slash if present
		cleanPath := strings.TrimPrefix(relativePath, "/")
		targetPath = filepath.Join(rootDir, cleanPath)
	}

	// Ensure the target path is within the root directory
	absTarget, err := filepath.Abs(targetPath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve absolute path: %w", err)
	}

	absRoot, err := filepath.Abs(rootDir)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve root directory: %w", err)
	}

	// Security check: ensure target is within root
	if !strings.HasPrefix(absTarget, absRoot) {
		return nil, fmt.Errorf("path is outside of allowed root directory")
	}

	// Check if target path exists before walking
	if _, err := os.Stat(absTarget); os.IsNotExist(err) {
		return nil, fmt.Errorf("path does not exist: %s", targetPath)
	}

	var results []string

	err = filepath.WalkDir(absTarget, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			// Log error but continue walking for permission issues
			log.Printf("Error accessing %s: %v", path, err)
			return nil
		}

		// Convert to absolute path and normalize separators
		absPath, err := filepath.Abs(path)
		if err != nil {
			log.Printf("Error getting absolute path for %s: %v", path, err)
			return nil
		}

		// Convert to forward slashes for consistent output
		normalizedPath := filepath.ToSlash(absPath)
		results = append(results, normalizedPath)

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}

	return results, nil
}

func main() {
	// Parse command line arguments
	var config Config
	flag.BoolVar(&config.UseStdio, "s", false, "Use stdio transport instead of HTTP")
	flag.Parse()

	// Validate required root directory argument
	args := flag.Args()
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s [-s] <root_directory>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  -s: Use stdio transport (default: HTTP)\n")
		os.Exit(1)
	}

	config.RootDirectory = args[0]

	// Validate root directory exists
	if _, err := os.Stat(config.RootDirectory); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: Root directory '%s' does not exist\n", config.RootDirectory)
		os.Exit(1)
	}

	// Get port from environment or use default
	config.Port = os.Getenv("PORT")
	if config.Port == "" {
		config.Port = "5001"
	}

	// Create MCP server
	impl := &mcp.Implementation{
		Name:    "filez-mcp",
		Title:   "Directory Walker MCP Server",
		Version: "1.0.0",
	}

	options := &mcp.ServerOptions{
		HasTools: true,
	}

	server := mcp.NewServer(impl, options)

	// Create input schema for walk_directory tool
	inputSchema := &jsonschema.Schema{
		Type: "object",
		Properties: map[string]*jsonschema.Schema{
			"path": {
				Type:        "string",
				Description: "Directory path to walk (use '/' for root directory)",
				Default:     []byte(`"/"`),
			},
		},
	}

	// Register the walk_directory tool
	tool := &mcp.Tool{
		Name:        "walk_directory",
		Description: "Recursively lists all files and directories under the specified path",
		InputSchema: inputSchema,
	}

	handler := func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		path := "/"
		if args, ok := req.Params.Arguments.(map[string]interface{}); ok {
			if p, ok := args["path"].(string); ok {
				path = p
			}
		}

		results, err := walkDirectory(config.RootDirectory, path)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					&mcp.TextContent{
						Text: fmt.Sprintf("Error walking directory: %v", err),
					},
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: strings.Join(results, "\n"),
				},
			},
		}, nil
	}

	server.AddTool(tool, handler)

	// Start server based on transport type
	if config.UseStdio {
		log.Println("Starting MCP server with stdio transport")
		transport := mcp.NewStdioTransport()
		if err := server.Run(context.Background(), transport); err != nil {
			log.Fatalf("Failed to start stdio server: %v", err)
		}
	} else {
		log.Printf("Starting MCP server with HTTP transport on port %s", config.Port)

		// Create streamable HTTP handler
		getServer := func(req *http.Request) *mcp.Server {
			return server
		}

		handler := mcp.NewStreamableHTTPHandler(getServer, &mcp.StreamableHTTPOptions{})

		// Set up HTTP server with MCP endpoint
		mux := http.NewServeMux()
		mux.Handle("/mcp", handler)

		addr := ":" + config.Port
		log.Printf("MCP server listening on %s/mcp", addr)

		if err := http.ListenAndServe(addr, mux); err != nil {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}
}
