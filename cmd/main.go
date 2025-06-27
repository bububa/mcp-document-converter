package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/bububa/mcp-document-converter/internal/tools"
)

func main() {
	// Create MCP server
	s := server.NewMCPServer(
		"mcp-document-converter",
		"1.0.0",
		server.WithLogging(),
		server.WithResourceCapabilities(true, true),
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

	var (
		host string
		port int
	)
	flag.StringVar(&host, "host", "", "sse server host")
	flag.IntVar(&port, "port", 0, "sse server port")
	flag.Parse()
	if host == "" || port == 0 {
		// Start server via stdio
		if err := server.ServeStdio(s); err != nil {
			os.Exit(1)
		}
		return
	}
	startCtx := context.Background()
	srv := server.NewSSEServer(s)
	stopCtx, stop := signal.NotifyContext(startCtx, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		if err := srv.Start(fmt.Sprintf("%s:%d", host, port)); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalln(err)
		}
	}()
	<-stopCtx.Done()
	stop()
	ctx, cancel := context.WithTimeout(startCtx, 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalln(err)
	}
}
