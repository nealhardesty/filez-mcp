package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

func main() {
	baseURL := "http://localhost:5001"
	if len(os.Args) > 1 {
		baseURL = os.Args[1]
	}

	// Add /mcp to the base URL if not present
	if !strings.HasSuffix(baseURL, "/mcp") {
		baseURL = baseURL + "/mcp"
	}

	fmt.Printf("Testing MCP Directory Walker via HTTP\n")
	fmt.Printf("Base URL: %s\n", baseURL)
	fmt.Println("=" + strings.Repeat("=", 50))

	// Wait a moment for server to be ready
	fmt.Println("Waiting for server to be ready...")
	time.Sleep(2 * time.Second)

	// Test 1: Initialize
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
				Name:    "test-http",
				Version: "1.0.0",
			},
		},
	}

	response, err := sendHTTPRequest(baseURL, initRequest)
	if err != nil {
		fmt.Printf("Failed to send initialize request: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Initialize response: %s\n", response)

	// Test 2: Send initialized notification
	fmt.Println("\n2. Sending initialized notification...")
	initializedNotification := mcp.InitializedNotification{
		Notification: mcp.Notification{
			Method: "notifications/initialized",
		},
	}

	_, err = sendHTTPRequest(baseURL, initializedNotification)
	if err != nil {
		fmt.Printf("Failed to send initialized notification: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Initialized notification sent successfully")

	// Test 3: List tools
	fmt.Println("\n3. Sending list_tools request...")
	listToolsRequest := mcp.ListToolsRequest{
		PaginatedRequest: mcp.PaginatedRequest{
			Request: mcp.Request{
				Method: "tools/list",
			},
		},
	}

	response, err = sendHTTPRequest(baseURL, listToolsRequest)
	if err != nil {
		fmt.Printf("Failed to send list_tools request: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("List tools response: %s\n", response)

	// Test 4: Call walk_directory tool with default path
	fmt.Println("\n4. Sending walk_directory tool call (default path)...")
	callToolRequest := mcp.CallToolRequest{
		Request: mcp.Request{
			Method: "tools/call",
		},
		Params: mcp.CallToolParams{
			Name: "walk_directory",
		},
	}

	response, err = sendHTTPRequest(baseURL, callToolRequest)
	if err != nil {
		fmt.Printf("Failed to send call_tool request: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Walk directory response (default): %s\n", response)

	// Test 5: Call walk_directory tool with specific path
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

	response, err = sendHTTPRequest(baseURL, callToolRequestWithPath)
	if err != nil {
		fmt.Printf("Failed to send call_tool request with path: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Walk directory response (with path): %s\n", response)

	// Test 6: Test with invalid path
	fmt.Println("\n6. Sending walk_directory tool call (invalid path)...")
	callToolRequestInvalid := mcp.CallToolRequest{
		Request: mcp.Request{
			Method: "tools/call",
		},
		Params: mcp.CallToolParams{
			Name:      "walk_directory",
			Arguments: json.RawMessage(`{"path": "/nonexistent"}`),
		},
	}

	response, err = sendHTTPRequest(baseURL, callToolRequestInvalid)
	if err != nil {
		fmt.Printf("Failed to send call_tool request with invalid path: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Walk directory response (invalid path): %s\n", response)

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("Test completed successfully!")
}

func sendHTTPRequest(baseURL string, request interface{}) (string, error) {
	// Create a proper JSON-RPC wrapper
	var jsonrpcRequest interface{}
	
	switch req := request.(type) {
	case mcp.InitializeRequest:
		jsonrpcRequest = map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      1,
			"method":  "initialize",
			"params":  req.Params,
		}
	case mcp.InitializedNotification:
		jsonrpcRequest = map[string]interface{}{
			"jsonrpc": "2.0",
			"method":  "notifications/initialized",
		}
	case mcp.ListToolsRequest:
		jsonrpcRequest = map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      2,
			"method":  "tools/list",
		}
	case mcp.CallToolRequest:
		// Generate unique ID for each call
		id := 3
		if req.Params.Name == "walk_directory" {
			if req.Params.Arguments != nil {
				// Check if it's the invalid path request
				if args, ok := req.Params.Arguments.(json.RawMessage); ok {
					if string(args) == `{"path": "/nonexistent"}` {
						id = 5
					} else {
						id = 4
					}
				}
			}
		}
		jsonrpcRequest = map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      id,
			"method":  "tools/call",
			"params":  req.Params,
		}
	default:
		// Fallback - just add JSON-RPC wrapper
		jsonrpcRequest = map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      1,
			"method":  "unknown",
			"params":  request,
		}
	}

	data, err := json.Marshal(jsonrpcRequest)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := http.Post(baseURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return "", fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Accept both 200 (OK) and 202 (Accepted) as successful responses
	// 202 is commonly returned for notifications
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("HTTP request failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	return string(body), nil
}
