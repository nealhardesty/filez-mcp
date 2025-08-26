package main

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
)

// Test walkDirectory function with various scenarios
func TestWalkDirectory(t *testing.T) {
	// Create a temporary directory structure for testing
	tempDir, err := os.MkdirTemp("", "mcp-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test directory structure
	testDirs := []string{
		"subdir1",
		"subdir1/subsubdir",
		"subdir2",
	}
	testFiles := []string{
		"file1.txt",
		"subdir1/file2.go",
		"subdir1/subsubdir/file3.json",
		"subdir2/file4.py",
	}

	// Create directories
	for _, dir := range testDirs {
		err := os.MkdirAll(filepath.Join(tempDir, dir), 0755)
		if err != nil {
			t.Fatalf("Failed to create test directory %s: %v", dir, err)
		}
	}

	// Create files
	for _, file := range testFiles {
		fullPath := filepath.Join(tempDir, file)
		err := os.WriteFile(fullPath, []byte("test content"), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", file, err)
		}
	}

	tests := []struct {
		name         string
		rootDir      string
		relativePath string
		wantErr      bool
		checkContent func([]string) bool
	}{
		{
			name:         "Walk root directory",
			rootDir:      tempDir,
			relativePath: "/",
			wantErr:      false,
			checkContent: func(results []string) bool {
				// Should contain the temp directory itself and all subdirectories and files
				expectedMinCount := 1 + len(testDirs) + len(testFiles) // root + dirs + files
				return len(results) >= expectedMinCount
			},
		},
		{
			name:         "Walk subdirectory",
			rootDir:      tempDir,
			relativePath: "/subdir1",
			wantErr:      false,
			checkContent: func(results []string) bool {
				// Should contain subdir1, its subdirectory, and its files
				subdir1Path := filepath.ToSlash(filepath.Join(tempDir, "subdir1"))
				for _, result := range results {
					if strings.Contains(result, subdir1Path) {
						return true
					}
				}
				return false
			},
		},
		{
			name:         "Walk empty relative path defaults to root",
			rootDir:      tempDir,
			relativePath: "",
			wantErr:      false,
			checkContent: func(results []string) bool {
				return len(results) > 0
			},
		},
		{
			name:         "Path outside root directory should fail",
			rootDir:      tempDir,
			relativePath: "/../",
			wantErr:      true,
			checkContent: nil,
		},
		{
			name:         "Non-existent subdirectory should fail",
			rootDir:      tempDir,
			relativePath: "/nonexistent",
			wantErr:      true,
			checkContent: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := walkDirectory(tt.rootDir, tt.relativePath)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("walkDirectory() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("walkDirectory() unexpected error: %v", err)
				return
			}

			// Check that all paths use forward slashes
			for _, result := range results {
				if strings.Contains(result, "\\") {
					t.Errorf("Result contains backslash (should use forward slashes): %s", result)
				}
			}

			// Check that all paths are absolute
			for _, result := range results {
				if !filepath.IsAbs(result) {
					t.Errorf("Result is not absolute path: %s", result)
				}
			}

			// Run custom content check if provided
			if tt.checkContent != nil && !tt.checkContent(results) {
				t.Errorf("walkDirectory() content check failed for results: %v", results)
			}
		})
	}
}

// Test that walkDirectory returns sorted results for deterministic output
func TestWalkDirectorySorted(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "mcp-sort-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create some files
	files := []string{"zebra.txt", "alpha.txt", "beta.txt"}
	for _, file := range files {
		err := os.WriteFile(filepath.Join(tempDir, file), []byte("test"), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	}

	results, err := walkDirectory(tempDir, "/")
	if err != nil {
		t.Fatalf("walkDirectory() failed: %v", err)
	}

	// Check if results are sorted (they should be due to filepath.WalkDir behavior)
	sortedResults := make([]string, len(results))
	copy(sortedResults, results)
	sort.Strings(sortedResults)

	if !reflect.DeepEqual(results, sortedResults) {
		t.Errorf("Results are not sorted. Got: %v, Want: %v", results, sortedResults)
	}
}

// Test path normalization
func TestPathNormalization(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "mcp-norm-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a subdirectory
	subDir := filepath.Join(tempDir, "testdir")
	err = os.Mkdir(subDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}

	tests := []struct {
		relativePath string
		description  string
	}{
		{"/testdir", "with leading slash"},
		{"testdir", "without leading slash"},
		{"/testdir/", "with trailing slash"},
		{"testdir/", "with trailing slash, no leading slash"},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			results, err := walkDirectory(tempDir, tt.relativePath)
			if err != nil {
				t.Errorf("walkDirectory() failed for %s: %v", tt.description, err)
				return
			}

			// Should find the subdirectory
			found := false
			expectedPath := filepath.ToSlash(subDir)
			for _, result := range results {
				if result == expectedPath {
					found = true
					break
				}
			}

			if !found {
				t.Errorf("Expected to find %s in results for %s, got: %v", expectedPath, tt.description, results)
			}
		})
	}
}

// Test error handling for invalid root directory
func TestWalkDirectoryInvalidRoot(t *testing.T) {
	_, err := walkDirectory("/nonexistent/directory", "/")
	if err == nil {
		t.Error("Expected error for non-existent root directory")
	}
}

// Benchmark walkDirectory performance
func BenchmarkWalkDirectory(b *testing.B) {
	tempDir, err := os.MkdirTemp("", "mcp-bench-*")
	if err != nil {
		b.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create some test files
	for i := 0; i < 100; i++ {
		filename := filepath.Join(tempDir, fmt.Sprintf("file_%d.txt", i))
		err := os.WriteFile(filename, []byte("test content"), 0644)
		if err != nil {
			b.Fatalf("Failed to create test file: %v", err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := walkDirectory(tempDir, "/")
		if err != nil {
			b.Fatalf("walkDirectory failed: %v", err)
		}
	}
}
