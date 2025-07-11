package main

import (
	"flag"

	"github.com/bububa/mcp-document-converter/internal"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 0, "sse server port")
	flag.Parse()
	internal.StartServer(port)
}
