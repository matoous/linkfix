package printers

import (
	"encoding/json"
	"fmt"

	"github.com/matoous/linkfix/internal/log"
    "github.com/matoous/linkfix/models"
)

type JSON struct {
}

func NewJSON() *JSON {
	return &JSON{}
}

type JSONResult struct {
	Issues []models.Fix `json:"issues"`
}

func (p JSON) PrintAll(issues []models.Fix) error {
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

func (p JSON) Print(issue models.Fix) error {
	outputJSON, err := json.Marshal(&issue)
	if err != nil {
		return err
	}

	_, err = fmt.Fprint(log.GetStdOut(), string(outputJSON))
	return err
}
