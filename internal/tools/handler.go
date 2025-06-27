package tools

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/bububa/mcp-document-converter/internal/pandoc"
)

// ConvertHandler handles document conversion requests
func ConvertHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Create converter
	converter, err := pandoc.NewConverter()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Pandoc: %v", err)
	}

	// Extract parameters
	args := req.GetArguments()

	var (
		input  = new(bytes.Buffer)
		output = new(bytes.Buffer)
	)
	convertReq := pandoc.ConvertRequest{
		Input:  input,
		Output: output,
	}
	if val, ok := args["input_file"]; ok {
		if b64, ok := val.(string); ok {
			if bs, err := base64.StdEncoding.DecodeString(b64); err != nil {
				return nil, fmt.Errorf("bad input_file encoding: %v", err)
			} else {
				input.Write(bs)
			}
		} else {
		}
	}
	if val, ok := args["input_format"]; ok {
		if format, ok := val.(string); ok {
			convertReq.InputFormat = format
		}
	}
	if val, ok := args["output_format"]; ok {
		if format, ok := val.(string); ok {
			convertReq.OutputFormat = format
		}
	}
	if input.Len() == 0 {
		return nil, errors.New("empty input_file content")
	}
	if convertReq.InputFormat == "" {
		convertReq.InputFormat = "markdown"
	}
	if convertReq.OutputFormat == "" {
		convertReq.OutputFormat = "docx"
	}

	// Convert file
	if err := converter.Convert(&convertReq); err != nil {
		return nil, err
	}
	b64 := base64.StdEncoding.EncodeToString(output.Bytes())
	return mcp.NewToolResultText(b64), nil
}
