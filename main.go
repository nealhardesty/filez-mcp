package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	mcp "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/http"
	"github.com/gin-gonic/gin"
)

type WalkDirectoryArgs struct {
	Path string `json:"path" jsonschema:"description=Directory path to walk (use '/' for root directory),default=/"`
}

type WalkDirectoryResult struct {
	Content []string `json:"content"`
}

type DirectoryWalkerServer struct {
	rootDirectory string
}

func NewDirectoryWalkerServer(rootDir string) *DirectoryWalkerServer {
	return &DirectoryWalkerServer{
		rootDirectory: rootDir,
	}
}

func (s *DirectoryWalkerServer) WalkDirectory(args WalkDirectoryArgs) (*mcp.ToolResponse, error) {
	// Handle path mapping: "/" maps to root directory, "/subdir" maps to root/subdir
	targetPath := s.rootDirectory
	if args.Path != "/" {
		// Remove leading "/" if present and join with root
		cleanPath := strings.TrimPrefix(args.Path, "/")
		targetPath = filepath.Join(s.rootDirectory, cleanPath)
	}

	// Convert to absolute path
	absPath, err := filepath.Abs(targetPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Verify path exists and is accessible
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("path does not exist: %s", absPath)
	} else if err != nil {
		return nil, fmt.Errorf("failed to access path: %w", err)
	}

	var paths []string

	// Walk the directory tree
	err = filepath.Walk(absPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// Log error but continue walking (graceful handling of permission errors)
			log.Printf("Error accessing %s: %v", path, err)
			return nil
		}

		// Convert to absolute path and normalize separators to '/'
		absWalkPath, err := filepath.Abs(path)
		if err != nil {
			log.Printf("Error getting absolute path for %s: %v", path, err)
			return nil
		}

		// Normalize path separators to '/' for cross-platform consistency
		normalizedPath := filepath.ToSlash(absWalkPath)
		paths = append(paths, normalizedPath)

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}

	// Create response text with file paths
	pathList := strings.Join(paths, "\\n")
	responseText := fmt.Sprintf("Found %d items:\\n%s", len(paths), pathList)
	return mcp.NewToolResponse(mcp.NewTextContent(responseText)), nil
}

func main() {
	// Validate command line arguments
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <root_directory>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Examples:\n")
		fmt.Fprintf(os.Stderr, "  %s /home/user/projects\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s ./my-project\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s .\n", os.Args[0])
		os.Exit(1)
	}

	rootDir := os.Args[1]

	// Validate root directory exists
	if _, err := os.Stat(rootDir); os.IsNotExist(err) {
		log.Fatalf("Root directory does not exist: %s", rootDir)
	} else if err != nil {
		log.Fatalf("Failed to access root directory: %v", err)
	}

	// Convert to absolute path for consistency
	absRootDir, err := filepath.Abs(rootDir)
	if err != nil {
		log.Fatalf("Failed to get absolute path for root directory: %v", err)
	}

	log.Printf("Starting MCP Directory Walker Server with root: %s", absRootDir)

	// Create server instance
	walkerServer := NewDirectoryWalkerServer(absRootDir)

	// Create Gin router
	router := gin.Default()

	// Create HTTP transport with Gin
	transport := http.NewGinTransport()
	router.POST("/mcp", transport.Handler())

	// Create MCP server
	mcpServer := mcp.NewServer(transport)

	// Register the walk_directory tool
	err = mcpServer.RegisterTool("walk_directory", "Recursively lists all files and directories under the specified path", walkerServer.WalkDirectory)
	if err != nil {
		log.Fatalf("Failed to register tool: %v", err)
	}

	// Start HTTP server
	log.Println("Starting MCP server on HTTP transport at localhost:8080...")
	log.Printf("MCP endpoint available at: http://localhost:8080/mcp")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}