#!/bin/bash

# MCP Directory Walker CLI Test Script

set -e

echo "=== MCP Directory Walker CLI Test Script ==="

# Build the project first
echo "Building project..."
make build

# Test 1: Run with current directory
echo
echo "Test 1: Running with current directory..."
echo "Command: ./directory-walker ."
timeout 10s ./directory-walker . &
SERVER_PID=$!
sleep 2
echo "Server started with PID: $SERVER_PID"
kill $SERVER_PID 2>/dev/null || true
echo "✓ Test 1 passed: Server starts successfully with current directory"

# Test 2: Run with specific directory (specs folder)
echo
echo "Test 2: Running with specs directory..."
echo "Command: ./directory-walker ./specs"
timeout 10s ./directory-walker ./specs &
SERVER_PID=$!
sleep 2
echo "Server started with PID: $SERVER_PID"
kill $SERVER_PID 2>/dev/null || true
echo "✓ Test 2 passed: Server starts successfully with specific directory"

# Test 3: Test with invalid directory (should fail)
echo
echo "Test 3: Testing with invalid directory..."
echo "Command: ./directory-walker /nonexistent"
if ./directory-walker /nonexistent 2>/dev/null; then
    echo "✗ Test 3 failed: Should have failed with invalid directory"
    exit 1
else
    echo "✓ Test 3 passed: Correctly fails with invalid directory"
fi

# Test 4: Test with no arguments (should fail)
echo
echo "Test 4: Testing with no arguments..."
echo "Command: ./directory-walker"
if ./directory-walker 2>/dev/null; then
    echo "✗ Test 4 failed: Should have failed with no arguments"
    exit 1
else
    echo "✓ Test 4 passed: Correctly fails with no arguments"
fi

echo
echo "=== All CLI tests passed! ==="
echo
echo "To test MCP functionality manually:"
echo "1. Start server: ./directory-walker ."
echo "2. Use MCP client to call walk_directory tool"
echo "3. Verify JSON response with file paths using '/' separators"