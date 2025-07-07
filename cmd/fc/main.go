package main

import (
	"os"
	"strconv"

	"github.com/bububa/mcp-document-converter/internal"
)

func main() {
	port, _ := strconv.Atoi(os.Getenv("MCP_DOCUMENT_CONVERTER_PORT"))
	internal.StartServer(port)
}
