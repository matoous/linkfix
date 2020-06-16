package printers

import (
	"fmt"
	"strings"

	"github.com/fatih/color"

	"github.com/matoous/linkfix/internal/log"
    "github.com/matoous/linkfix/models"
)

type Text struct {
	useColors bool
}

func NewText(useColors bool) *Text {
	return &Text{
		useColors: useColors,
	}
}

func (p *Text) PrintAll(issues []models.Fix) error {
	for i := range issues {
		err := p.Print(issues[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Text) Print(issue models.Fix) error {
	var s strings.Builder
	s.WriteString(color.New(color.Bold).Sprintf("%s:%d:%d", issue.Path, issue.Line, issue.Index))
	s.WriteRune('\t')
	s.WriteString(issue.URL.String())
	s.WriteRune('\t')
	if issue.Severity == "error" {
		s.WriteString(color.New(color.FgRed).Sprintf("%s", issue.Reason))
	} else {
		s.WriteString(color.New(color.FgYellow).Sprintf("%s", issue.Reason))
	}
	s.WriteRune('\n')
	_, err := fmt.Fprint(log.GetStdOut(), s.String())
	return err
}
