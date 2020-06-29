package printers

import (
	"fmt"

	"github.com/matoous/linkfix/internal/log"
	"github.com/matoous/linkfix/models"
)

type githuber struct {
}

// Github output format outputs issues according to Github actions format:
// https://help.github.com/en/actions/reference/workflow-commands-for-github-actions#setting-an-error-message
func NewGithub() Printer {
	return &githuber{}
}

// PrintAll prints all issues using the Print function. See `Print` for detailed explanation of the format.
func (g *githuber) PrintAll(issues []models.Fix) error {
	for i := range issues {
		err := g.Print(issues[i])
		if err != nil {
			return err
		}
	}
	return nil
}

// Print prints single link fix in the github format.
// Format is: `::warn file=file.md,line=10,col=15::Something went wrong (old URL => new URL)`.
func (g *githuber) Print(issue models.Fix) error {
	_, err := fmt.Fprintf(
		log.GetStdOut(),
		"::warn file=%s,line=%d,col=%d::%s (%s => %s)\n",
		issue.Path, issue.Line, issue.Index, issue.Reason, issue.URL, issue.Suggestion,
	)
	return err
}
