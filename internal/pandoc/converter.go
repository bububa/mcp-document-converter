package pandoc

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"slices"
)

// PandocConverter represents a document converter based on Pandoc
type PandocConverter struct {
	pandocPath string
}

// NewConverter creates a new document converter
func NewConverter() (*PandocConverter, error) {
	// Try to get path from environment variable
	pandocPath := os.Getenv("PANDOC_PATH")
	if pandocPath == "" {
		// If variable is not set, look for pandoc in system PATH
		path, err := exec.LookPath("pandoc")
		if err != nil {
			return nil, errors.New("pandoc not found in PATH, please set PANDOC_PATH environment variable")
		}
		pandocPath = path
	}

	// Check that file exists and is executable
	info, err := os.Stat(pandocPath)
	if err != nil {
		return nil, fmt.Errorf("error accessing Pandoc at %s: %v", pandocPath, err)
	}

	if info.IsDir() {
		return nil, fmt.Errorf("%s is a directory, not a Pandoc executable", pandocPath)
	}

	// On Windows don't check execution permissions
	if os.Getenv("OS") != "Windows_NT" && info.Mode()&0111 == 0 {
		return nil, fmt.Errorf("%s is not executable", pandocPath)
	}

	return &PandocConverter{
		pandocPath: pandocPath,
	}, nil
}

// ValidateFormat checks if the format is supported
func (p *PandocConverter) ValidateFormat(format string) bool {
	supportedFormats := []string{"markdown", "html", "pdf", "docx", "rst", "latex", "epub", "txt"}
	return slices.Contains(supportedFormats, format)
}

// Convert converts a file from one format to another
func (p *PandocConverter) Convert(req *ConvertRequest) error {
	// Format validation
	if !p.ValidateFormat(req.InputFormat) || !p.ValidateFormat(req.OutputFormat) {
		return fmt.Errorf("unsupported format: input=%s, output=%s", req.InputFormat, req.OutputFormat)
	}

	tmpInput, err := os.CreateTemp("", "pandoc-input-*."+req.InputFormat)
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpInput.Name())
	if _, err := io.Copy(tmpInput, req.Input); err != nil {
		return fmt.Errorf("failed to write input temp file: %v", err)
	}
	inputFile := tmpInput.Name()
	tmpInput.Close()

	tmpOutput, err := os.CreateTemp("", "pandoc-output-*."+req.OutputFormat)
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpOutput.Name())
	outputFile := tmpOutput.Name()
	tmpOutput.Close()

	// Run pandoc
	args := []string{
		"--standalone",
		"--highlight-style", "tango",
		"--wrap", "preserve",
		"--toc",
		"-f", req.InputFormat,
		"-t", req.OutputFormat,
		"-o", outputFile,
	}

	// Add input file at the end
	args = append(args, inputFile)

	cmd := exec.Command(p.pandocPath, args...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("pandoc conversion failed: %v\noutput: %s", err, string(output))
	}
	if fn, err := os.Open(outputFile); err != nil {
		return fmt.Errorf("failed to read output file: %v", err)
	} else if _, err := io.Copy(req.Output, fn); err != nil {
		return fmt.Errorf("failed to copy output file to output reader: %v", err)
	}

	return nil
}
