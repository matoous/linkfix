package models

import (
    "net/url"
)

type Link struct {
    Path  string   `json:"file"`
    Line  int      `json:"line"`
    Index int      `json:"index"`
    URL   *url.URL `json:"url"`
}
