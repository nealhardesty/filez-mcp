package main

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

// createTempTestDir creates a temporary directory structure for testing
func createTempTestDir(t *testing.T) string {
	tempDir, err := os.MkdirTemp("", "mcp-walker-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Create test structure:
	// tempDir/
	//   ├── file1.txt
	//   ├── subdir/
	//   │   ├── file2.go
	//   │   └── deep/
	//   │       └── file3.json
	//   └── emptydir/

	if err := os.WriteFile(filepath.Join(tempDir, "file1.txt"), []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create file1.txt: %v", err)
	}

	subdirPath := filepath.Join(tempDir, "subdir")
	if err := os.MkdirAll(subdirPath, 0755); err != nil {
		t.Fatalf("Failed to create subdir: %v", err)
	}

	if err := os.WriteFile(filepath.Join(subdirPath, "file2.go"), []byte("package main"), 0644); err != nil {
		t.Fatalf("Failed to create file2.go: %v", err)
	}

	deepPath := filepath.Join(subdirPath, "deep")
	if err := os.MkdirAll(deepPath, 0755); err != nil {
		t.Fatalf("Failed to create deep dir: %v", err)
	}

	if err := os.WriteFile(filepath.Join(deepPath, "file3.json"), []byte(`{"test": true}`), 0644); err != nil {
		t.Fatalf("Failed to create file3.json: %v", err)
	}

	emptyDirPath := filepath.Join(tempDir, "emptydir")
	if err := os.MkdirAll(emptyDirPath, 0755); err != nil {
		t.Fatalf("Failed to create emptydir: %v", err)
	}

	return tempDir
}

func TestWalkDirectoryTool_RootPath(t *testing.T) {
	tempDir := createTempTestDir(t)
	defer os.RemoveAll(tempDir)

	handler := walkDirectoryTool(tempDir)

	// Create request for root path
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      "walk_directory",
			Arguments: json.RawMessage(`{"path": "/"}`),
		},
	}

	result, err := handler(context.Background(), request)
	if err != nil {
		t.Fatalf("Handler returned error: %v", err)
	}

	if result == nil {
		t.Fatal("Result is nil")
	}

	// Check that we have structured content
	if result.StructuredContent == nil {
		t.Fatal("StructuredContent is nil")
	}

	// Extract the file list from structured content
	files, ok := result.StructuredContent.([]string)
	if !ok {
		t.Fatalf("StructuredContent is not []string, got %T", result.StructuredContent)
	}

	// We should have at least 6 items: tempDir, file1.txt, subdir, file2.go, deep, file3.json, emptydir
	if len(files) < 6 {
		t.Fatalf("Expected at least 6 files/dirs, got %d: %v", len(files), files)
	}

	// Convert to normalized paths for comparison
	expectedFiles := []string{
		filepath.ToSlash(tempDir),
		filepath.ToSlash(filepath.Join(tempDir, "file1.txt")),
		filepath.ToSlash(filepath.Join(tempDir, "subdir")),
		filepath.ToSlash(filepath.Join(tempDir, "subdir", "file2.go")),
		filepath.ToSlash(filepath.Join(tempDir, "subdir", "deep")),
		filepath.ToSlash(filepath.Join(tempDir, "subdir", "deep", "file3.json")),
		filepath.ToSlash(filepath.Join(tempDir, "emptydir")),
	}

	// Check that all expected files are present
	fileMap := make(map[string]bool)
	for _, file := range files {
		fileMap[file] = true
	}

	for _, expected := range expectedFiles {
		if !fileMap[expected] {
			t.Errorf("Expected file %s not found in results", expected)
		}
	}
}

func TestWalkDirectoryTool_SubdirPath(t *testing.T) {
	tempDir := createTempTestDir(t)
	defer os.RemoveAll(tempDir)

	handler := walkDirectoryTool(tempDir)

	// Create request for subdir path
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      "walk_directory",
			Arguments: json.RawMessage(`{"path": "/subdir"}`),
		},
	}

	result, err := handler(context.Background(), request)
	if err != nil {
		t.Fatalf("Handler returned error: %v", err)
	}

	files, ok := result.StructuredContent.([]string)
	if !ok {
		t.Fatalf("StructuredContent is not []string, got %T", result.StructuredContent)
	}

	// Should contain subdir, file2.go, deep dir, and file3.json
	if len(files) < 4 {
		t.Fatalf("Expected at least 4 files/dirs in subdir, got %d: %v", len(files), files)
	}

	// Check for specific files
	fileMap := make(map[string]bool)
	for _, file := range files {
		fileMap[file] = true
	}

	expectedInSubdir := []string{
		filepath.ToSlash(filepath.Join(tempDir, "subdir", "file2.go")),
		filepath.ToSlash(filepath.Join(tempDir, "subdir", "deep")),
		filepath.ToSlash(filepath.Join(tempDir, "subdir", "deep", "file3.json")),
	}

	for _, expected := range expectedInSubdir {
		if !fileMap[expected] {
			t.Errorf("Expected file %s not found in subdir results", expected)
		}
	}
}

func TestWalkDirectoryTool_DefaultPath(t *testing.T) {
	tempDir := createTempTestDir(t)
	defer os.RemoveAll(tempDir)

	handler := walkDirectoryTool(tempDir)

	// Create request with no arguments (should default to "/")
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "walk_directory",
		},
	}

	result, err := handler(context.Background(), request)
	if err != nil {
		t.Fatalf("Handler returned error: %v", err)
	}

	files, ok := result.StructuredContent.([]string)
	if !ok {
		t.Fatalf("StructuredContent is not []string, got %T", result.StructuredContent)
	}

	// Should include all files since it defaults to root
	if len(files) < 6 {
		t.Fatalf("Expected at least 6 files/dirs with default path, got %d: %v", len(files), files)
	}
}

func TestWalkDirectoryTool_NonexistentPath(t *testing.T) {
	tempDir := createTempTestDir(t)
	defer os.RemoveAll(tempDir)

	handler := walkDirectoryTool(tempDir)

	// Create request for nonexistent path
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      "walk_directory",
			Arguments: json.RawMessage(`{"path": "/nonexistent"}`),
		},
	}

	_, err := handler(context.Background(), request)
	if err == nil {
		t.Fatal("Expected error for nonexistent path, but got none")
	}

	if !contains(err.Error(), "does not exist") {
		t.Errorf("Expected 'does not exist' error, got: %v", err)
	}
}

func TestWalkDirectoryTool_SecurityCheck(t *testing.T) {
	tempDir := createTempTestDir(t)
	defer os.RemoveAll(tempDir)

	handler := walkDirectoryTool(tempDir)

	// Try to access parent directory (security check)
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      "walk_directory",
			Arguments: json.RawMessage(`{"path": "/../.."}`),
		},
	}

	_, err := handler(context.Background(), request)
	if err == nil {
		t.Fatal("Expected error for path outside root, but got none")
	}

	if !contains(err.Error(), "outside root directory") {
		t.Errorf("Expected 'outside root directory' error, got: %v", err)
	}
}

func TestWalkDirectoryTool_InvalidJSON(t *testing.T) {
	tempDir := createTempTestDir(t)
	defer os.RemoveAll(tempDir)

	handler := walkDirectoryTool(tempDir)

	// Create request with invalid JSON
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      "walk_directory",
			Arguments: json.RawMessage(`{"path": }`), // Invalid JSON
		},
	}

	_, err := handler(context.Background(), request)
	if err == nil {
		t.Fatal("Expected error for invalid JSON, but got none")
	}

	if !contains(err.Error(), "failed to parse arguments") {
		t.Errorf("Expected 'failed to parse arguments' error, got: %v", err)
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
		(len(s) > len(substr) && containsHelper(s, substr)))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
