package printers

import (
	"fmt"
	"strings"

	"github.com/fatih/color"

	"github.com/matoous/linkfix/internal/log"
	"github.com/matoous/linkfix/models"
)

type texter struct {
	useColors bool
}

func NewText(useColors bool) *texter {
	return &texter{
		useColors: useColors,
	}
}

func (p *texter) PrintAll(issues []models.Fix) error {
	for i := range issues {
		err := p.Print(issues[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *texter) Print(issue models.Fix) error {
	var s strings.Builder
	s.WriteString(color.New(color.Bold).Sprintf("%s:%d:%d", issue.Path, issue.Line, issue.Index))
	s.WriteRune('\t')
	s.WriteString(issue.URL.String())
	s.WriteRune('\t')
	s.WriteString(issue.Severity.Color().Sprint(issue.Reason))
	s.WriteRune('\n')
	_, err := fmt.Fprint(log.GetStdOut(), s.String())
	return err
}
