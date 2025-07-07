# mcp-document-converter

A simple and powerful document conversion MCP server using Pandoc.

## When to use this

✅ **When you need to convert documents between different formats**  
Converting markdown to Word, HTML to PDF, or any combination of supported formats.

✅ **When you want to generate professional documents in Cursor**  
Creating Word documents, PDFs, or other formats directly from your markdown content.

✅ **When you need consistent document branding**  
Automatically adds copyright and branding to all generated documents.

✅ **When you need a reliable, high-quality document converter**  
Built on the industry-standard Pandoc conversion engine.

## Features

- Fast document conversion through the Cursor MCP API
- Supports markdown, HTML, PDF, DOCX, RST, LaTeX, EPUB, TXT
- Automatic path normalization for Windows compatibility
- Multiple conversion modes: string-to-string, string-to-file, file-to-file
- Automatic copyright addition to all generated documents

## Manual Installation

### Prerequisites

- [Pandoc](https://pandoc.org/installing.html) must be installed and available in your PATH
  - Download and install manually from the [official website](https://pandoc.org/installing.html)
  - **Important**: Restart your computer after installing Pandoc
  - You can verify installation by running `pandoc --version` in terminal
- For PDF generation, a LaTeX distribution is required (MiKTeX recommended for Windows)
- Git and Go programming language must be installed

### Manual Build

1. Clone this repository
2. Build the server:
   ```
   go build -o mcp-document-converter ./cmd/main.go
   ```
3. Run the server:
   ```
   ./mcp-document-converter
   ```

## Integration with Cursor IDE

To integrate with Cursor IDE, add the following to your MCP configuration file (`mcp.json`):

```json
// In the "servers" section
"pandoc_mcp_go": {
  "type": "stdio",
  "command": "mcp-ducument-converter",
}

// In the "roots" section
"pandoc": {
  "type": "document-converter",
  "server": "pandoc_mcp_go"
}
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
