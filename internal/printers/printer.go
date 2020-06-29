package printers

import (
	"github.com/matoous/linkfix/models"
)

type Printer interface {
	Print(issue models.Fix) error
	PrintAll(issues []models.Fix) error
}

func GetPrinter(typ string) Printer {
	switch typ {
	case "github":
		return NewGithub()
	case "json":
		return NewJSON()
	default:
		return NewText(true)
	}
}
