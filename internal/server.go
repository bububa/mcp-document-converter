package internal

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/bububa/mcp-document-converter/internal/tools"
)

const (
	AppName    = "mcp-document-converter"
	AppVersion = "1.0.0"
)

func StartServer(port int) {
	// Create MCP server
	s := server.NewMCPServer(
		AppName,
		AppVersion,
		server.WithLogging(),
		server.WithToolCapabilities(true),
	)

	// Register convert_contents tool
	convertTool := mcp.NewTool("convert",
		mcp.WithDescription("Convert document between different formats using Pandoc"),
		mcp.WithString("input_file",
			mcp.Description("Base64 encoded data of input file"),
			mcp.Required(),
		),
		mcp.WithString("input_format",
			mcp.Description("Source format of the content"),
			mcp.DefaultString("markdown"),
			mcp.Enum("markdown", "html", "pdf", "docx", "rst", "latex", "epub", "txt"),
		),
		mcp.WithString("output_format",
			mcp.Description("Target format"),
			mcp.DefaultString("docx"),
			mcp.Enum("markdown", "html", "pdf", "docx", "rst", "latex", "epub", "txt"),
		),
	)

	// Add tool handler
	s.AddTool(convertTool, tools.ConvertHandler)
	if port == 0 {
		// Start server via stdio
		if err := server.ServeStdio(s); err != nil {
			os.Exit(1)
		}
		return
	}
	mux := http.NewServeMux()
	registerHealthAndVersion(mux)
	srv := server.NewStreamableHTTPServer(s)
	mux.Handle("/", srv)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), mux); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalln(err)
	}
}

// registerHealthAndVersion adds the /health and /version endpoints.
func registerHealthAndVersion(mux *http.ServeMux) {
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		// Matching the documented output from the project's README.
		fmt.Fprint(w, `{"status":"ok"}`)
	})

	mux.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"name":"%s","version":"%s"}`, AppName, AppVersion)
	})
}
