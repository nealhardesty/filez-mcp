package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <directory-walker-binary> [root-directory]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Example: %s ../directory-walker .\n", os.Args[0])
		os.Exit(1)
	}

	binary := os.Args[1]
	rootDir := "."
	if len(os.Args) > 2 {
		rootDir = os.Args[2]
	}

	fmt.Printf("Testing MCP Directory Walker via stdio\n")
	fmt.Printf("Binary: %s\n", binary)
	fmt.Printf("Root Directory: %s\n", rootDir)
	fmt.Println("=" + strings.Repeat("=", 50))

	// Start the directory-walker process with stdio transport
	cmd := exec.Command(binary, "-s", rootDir)
	
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatalf("Failed to get stdin pipe: %v", err)
	}
	
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("Failed to get stdout pipe: %v", err)
	}
	
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatalf("Failed to get stderr pipe: %v", err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatalf("Failed to start command: %v", err)
	}

	// Monitor stderr in a goroutine
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			fmt.Printf("[SERVER STDERR] %s\n", scanner.Text())
		}
	}()

	// Create a scanner for reading responses
	scanner := bufio.NewScanner(stdout)

	// Send initialize request
	fmt.Println("\n1. Sending initialize request...")
	initRequest := mcp.InitializeRequest{
		Request: mcp.Request{
			Method: "initialize",
		},
		Params: mcp.InitializeParams{
			ProtocolVersion: "2024-11-05",
			Capabilities: mcp.ClientCapabilities{
				Experimental: map[string]interface{}{},
				Sampling:     &struct{}{},
			},
			ClientInfo: mcp.Implementation{
				Name:    "test-stdio",
				Version: "1.0.0",
			},
		},
	}

	if err := sendRequest(stdin, initRequest); err != nil {
		log.Fatalf("Failed to send initialize request: %v", err)
	}

	// Read initialize response
	response, err := readResponse(scanner)
	if err != nil {
		log.Fatalf("Failed to read initialize response: %v", err)
	}
	fmt.Printf("Initialize response: %s\n", response)

	// Send initialized notification
	fmt.Println("\n2. Sending initialized notification...")
	initializedNotification := mcp.InitializedNotification{
		Notification: mcp.Notification{
			Method: "notifications/initialized",
		},
	}

	if err := sendRequest(stdin, initializedNotification); err != nil {
		log.Fatalf("Failed to send initialized notification: %v", err)
	}

	// Send list_tools request
	fmt.Println("\n3. Sending list_tools request...")
	listToolsRequest := mcp.ListToolsRequest{
		PaginatedRequest: mcp.PaginatedRequest{
			Request: mcp.Request{
				Method: "tools/list",
			},
		},
	}

	if err := sendRequest(stdin, listToolsRequest); err != nil {
		log.Fatalf("Failed to send list_tools request: %v", err)
	}

	// Read list_tools response
	response, err = readResponse(scanner)
	if err != nil {
		log.Fatalf("Failed to read list_tools response: %v", err)
	}
	fmt.Printf("List tools response: %s\n", response)

	// Send walk_directory tool call with default path
	fmt.Println("\n4. Sending walk_directory tool call (default path)...")
	callToolRequest := mcp.CallToolRequest{
		Request: mcp.Request{
			Method: "tools/call",
		},
		Params: mcp.CallToolParams{
			Name: "walk_directory",
		},
	}

	if err := sendRequest(stdin, callToolRequest); err != nil {
		log.Fatalf("Failed to send call_tool request: %v", err)
	}

	// Read call_tool response
	response, err = readResponse(scanner)
	if err != nil {
		log.Fatalf("Failed to read call_tool response: %v", err)
	}
	fmt.Printf("Walk directory response (default): %s\n", response)

	// Send walk_directory tool call with specific path
	fmt.Println("\n5. Sending walk_directory tool call (specific path)...")
	callToolRequestWithPath := mcp.CallToolRequest{
		Request: mcp.Request{
			Method: "tools/call",
		},
		Params: mcp.CallToolParams{
			Name:      "walk_directory",
			Arguments: json.RawMessage(`{"path": "/"}`),
		},
	}

	if err := sendRequest(stdin, callToolRequestWithPath); err != nil {
		log.Fatalf("Failed to send call_tool request with path: %v", err)
	}

	// Read call_tool response
	response, err = readResponse(scanner)
	if err != nil {
		log.Fatalf("Failed to read call_tool response with path: %v", err)
	}
	fmt.Printf("Walk directory response (with path): %s\n", response)

	// Close stdin to signal end of communication
	stdin.Close()

	// Wait for process to finish
	if err := cmd.Wait(); err != nil {
		log.Printf("Process finished with error: %v", err)
	}

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("Test completed successfully!")
}

func sendRequest(stdin io.WriteCloser, request interface{}) error {
	data, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	_, err = fmt.Fprintf(stdin, "%s\n", data)
	return err
}

func readResponse(scanner *bufio.Scanner) (string, error) {
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return "", fmt.Errorf("scanner error: %w", err)
		}
		return "", fmt.Errorf("no response received")
	}
	return scanner.Text(), nil
}
