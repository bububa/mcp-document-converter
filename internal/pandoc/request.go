package pandoc

import "io"

type ConvertRequest struct {
	Input        io.Reader
	Output       io.ReadWriter
	InputFormat  string
	OutputFormat string
}
