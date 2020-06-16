package log

import (
	"io"

	"github.com/fatih/color"
)

func GetStdOut() io.Writer {
	return color.Output
}

func GetStdErr() io.Writer {
	return color.Error
}
