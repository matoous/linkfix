package severity

import (
	"github.com/fatih/color"
)

//go:generate stringer -type=Severity

// Severity represents severity of the link issue found
type Severity uint8

const (
	// Error is used when the link doesn't work and there's high chance it won't work ever again
	Error Severity = iota
	// Warn is used for less several link issues such as using http instead of https or if the error is most likely
	// only temporary
	Warn
	Info
)

func (i Severity) Color() *color.Color {
	switch i {
	case Error:
		return color.New(color.FgRed)
	case Warn:
		return color.New(color.FgYellow)
	default:
		return color.New(color.FgBlue)
	}
}
