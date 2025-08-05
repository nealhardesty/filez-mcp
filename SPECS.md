This is a golang MCP server that can walk all of the files underneath it's directory.

it exposes a single MCP server tool named `walk_directory` that returns a list of files recursively.

Use '/' as file separators, no matter what the environment.

It should remap '/' to the current directory it is launched from.

Requirements:
* implemented in modern golang
* as simple as possible (single golang file is fine)
* uses the official golang MCP library