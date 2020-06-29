package check

import (
	"errors"
	"os"

	"github.com/matoous/linkfix/models"
	"github.com/matoous/linkfix/models/severity"
)

func File(link models.Link) (models.Fix, error) {
	fix := models.Fix{Link: link}
	_, err := os.Stat(link.URL.Path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fix.Reason = "file doesn't exist"
			fix.Severity = severity.Error
			return fix, nil
		}
		return fix, err
	}
	return fix, nil
}
