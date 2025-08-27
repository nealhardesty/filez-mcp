package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// walkDirectoryTool implements the walk_directory tool handler
func walkDirectoryTool(rootDir string) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Log the tool call with arguments
		if request.Params.Arguments != nil {
			if args, ok := request.Params.Arguments.(json.RawMessage); ok {
				log.Printf("[TOOL CALL] %s with arguments: %s", request.Params.Name, string(args))
			} else {
				log.Printf("[TOOL CALL] %s with arguments: %v", request.Params.Name, request.Params.Arguments)
			}
		} else {
			log.Printf("[TOOL CALL] %s with no arguments", request.Params.Name)
		}
		// Parse the arguments
		var args struct {
			Path string `json:"path"`
		}
		
		// Set default path to "/"
		args.Path = "/"
		
		// Parse arguments if provided
		if err := request.BindArguments(&args); err != nil {
			return nil, fmt.Errorf("failed to parse arguments: %w", err)
		}
		
		// Map the input path to actual filesystem path
		var targetPath string
		if args.Path == "/" {
			targetPath = rootDir
		} else {
			// Remove leading slash and join with root
			cleanPath := strings.TrimPrefix(args.Path, "/")
			targetPath = filepath.Join(rootDir, cleanPath)
		}
		
		// Ensure the target path is within the root directory (security check)
		absRoot, err := filepath.Abs(rootDir)
		if err != nil {
			return nil, fmt.Errorf("failed to get absolute root path: %w", err)
		}
		
		absTarget, err := filepath.Abs(targetPath)
		if err != nil {
			return nil, fmt.Errorf("failed to get absolute target path: %w", err)
		}
		
		if !strings.HasPrefix(absTarget, absRoot) {
			return nil, fmt.Errorf("path is outside root directory")
		}
		
		// Check if target path exists
		if _, err := os.Stat(absTarget); os.IsNotExist(err) {
			return nil, fmt.Errorf("path does not exist: %s", args.Path)
		}
		
		// Walk the directory tree
		var files []string
		err = filepath.WalkDir(absTarget, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				// Log permission errors but continue walking
				if os.IsPermission(err) {
					log.Printf("Permission denied: %s", path)
					return nil
				}
				return err
			}
			
			// Convert to absolute path and normalize separators
			absPath, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			
			// Convert to forward slashes for cross-platform consistency
			normalizedPath := filepath.ToSlash(absPath)
			files = append(files, normalizedPath)
			
			return nil
		})
		
		if err != nil {
			return nil, fmt.Errorf("failed to walk directory: %w", err)
		}
		
		// Create result with structured content
		result := mcp.NewToolResultStructured(files, fmt.Sprintf("Found %d files and directories", len(files)))
		
		// Log the completion
		log.Printf("[TOOL COMPLETED] %s - found %d files/directories", request.Params.Name, len(files))
		
		return result, nil
	}
}

// loggingMiddleware logs all incoming HTTP requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Log request details
		log.Printf("[HTTP REQUEST] %s %s from %s", r.Method, r.RequestURI, r.RemoteAddr)
		
		// Call the next handler
		next.ServeHTTP(w, r)
		
		// Log completion
		log.Printf("[HTTP COMPLETED] %s %s - took %s", r.Method, r.RequestURI, time.Since(start))
	})
}

func main() {
	// Parse command line arguments
	var useStdio bool
	flag.BoolVar(&useStdio, "s", false, "Use stdio transport instead of HTTP")
	flag.Parse()
	
	// Get root directory argument
	args := flag.Args()
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s [-s] <root_directory>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  -s: Use stdio transport instead of HTTP\n")
		os.Exit(1)
	}
	
	rootDir := args[0]
	
	// Validate root directory exists
	if _, err := os.Stat(rootDir); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: Root directory does not exist: %s\n", rootDir)
		os.Exit(1)
	}
	
	// Convert to absolute path for consistency
	absRootDir, err := filepath.Abs(rootDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to get absolute path: %v\n", err)
		os.Exit(1)
	}
	
	// Create MCP server with logging
	mcpServer := server.NewMCPServer("directory-walker", "1.0.0")
	log.Printf("MCP Server created: directory-walker v1.0.0")
	
	// Define the walk_directory tool
	walkTool := server.ServerTool{
		Tool: mcp.Tool{
			Name:        "walk_directory",
			Description: "Recursively lists all files and directories under the specified path",
			InputSchema: mcp.ToolInputSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"path": map[string]interface{}{
						"type":        "string",
						"description": "Directory path to walk (use '/' for root directory)",
						"default":     "/",
					},
				},
			},
		},
		Handler: walkDirectoryTool(absRootDir),
	}
	
	// Register the tool
	mcpServer.AddTools(walkTool)
	log.Printf("Registered tool: %s", walkTool.Tool.Name)
	
	// Start server based on transport mode
	if useStdio {
		fmt.Fprintf(os.Stderr, "Starting MCP Directory Walker Server (stdio) for root: %s\n", absRootDir)
		err = server.ServeStdio(mcpServer)
	} else {
		// HTTP server
		port := os.Getenv("PORT")
		if port == "" {
			port = "5001"
		}
		
		portNum, err := strconv.Atoi(port)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Invalid PORT value: %s\n", port)
			os.Exit(1)
		}
		
		fmt.Fprintf(os.Stderr, "Starting MCP Directory Walker Server (HTTP) on port %d for root: %s\n", portNum, absRootDir)
		
		// Create HTTP server - using StreamableHTTPServer for the /mcp path
		httpServer := server.NewStreamableHTTPServer(mcpServer, 
			server.WithEndpointPath("/mcp"),
			server.WithStateLess(true),
		)
		
		// Add basic request logging information
		log.Printf("HTTP MCP Server ready to accept requests on http://localhost:%d/mcp", portNum)
		
		err = httpServer.Start(fmt.Sprintf(":%d", portNum))
	}
	
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error starting server: %v\n", err)
		os.Exit(1)
	}
}
