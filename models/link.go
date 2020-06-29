package models

import (
	"net/url"
)

// Link is single link found in a file. This structure holds more information about the link so that it can be later
// found, fixed/replaced and otherwise adjusted based on the users preferences.
type Link struct {
	Path  string   `json:"file"`
	Line  int      `json:"line"`
	Index int      `json:"index"`
	URL   *url.URL `json:"url"`
}
