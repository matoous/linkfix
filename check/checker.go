package check

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/matoous/linkfix/models"
)

// Lister is interface to sources that can list links.
type Lister interface {
	// List lists all links from given source.
	List() ([]models.Link, error)
}

// Checker checks links for availability and finds replacements if necessary.
type Checker struct {
	client  http.Client
	exclude []string
}

// NewChecker creates new checker for link checking.
func NewChecker(exclude []string) *Checker {
	return &Checker{
		exclude: exclude,
		client: http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// ProcessURL processes given URL, checking whether it works and in case it doesn't finding its snapshot
// on Wayback machine.
func (c *Checker) ProcessURL(ctx context.Context, link models.Link) (models.Fix, error) {
	fix := models.Fix{Link: link}
	if !c.shouldCheck(link.URL.String()) {
		return fix, nil
	}
	switch link.URL.Scheme {
	case "http", "https":
		return HTTP(ctx, link)
	case "mailto":
		return MailTo(link)
	case "ftp":
		return FTP(link)
	case "file":
		return File(link)
	default:
		return models.Fix{}, fmt.Errorf("unknown URL schema: %s", link.URL.Scheme)
	}
}

// shouldCheck returns true if the url should be checked, false otherwise. This is determined based on excluded urls.
func (c *Checker) shouldCheck(url string) bool {
	for _, e := range c.exclude {
		if strings.Contains(url, e) {
			return false
		}
	}
	return true
}
