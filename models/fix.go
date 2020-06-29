package models

import (
	"github.com/matoous/linkfix/models/severity"
)

// Fix represents fix for given link that is shown/suggested to the user.
type Fix struct {
	Link
	Severity   severity.Severity `json:"severity"`
	Reason     string            `json:"reason"`
	Suggestion string            `json:"suggestion"`
}
