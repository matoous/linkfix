package printers

import (
	"encoding/json"
	"fmt"

	"github.com/matoous/linkfix/internal/log"
	"github.com/matoous/linkfix/models"
)

type jsoner struct {
}

func NewJSON() *jsoner {
	return &jsoner{}
}

type JSONResult struct {
	Issues []models.Fix `json:"issues"`
}

func (p jsoner) PrintAll(issues []models.Fix) error {
	res := JSONResult{
		Issues: issues,
	}

	outputJSON, err := json.Marshal(res)
	if err != nil {
		return err
	}

	_, err = fmt.Fprint(log.GetStdOut(), string(outputJSON))
	return err
}

func (p jsoner) Print(issue models.Fix) error {
	outputJSON, err := json.Marshal(&issue)
	if err != nil {
		return err
	}

	_, err = fmt.Fprint(log.GetStdOut(), string(outputJSON))
	return err
}
