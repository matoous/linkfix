package check

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/matoous/linkfix/internal/clients"
	"github.com/matoous/linkfix/models"
	"github.com/matoous/linkfix/models/severity"
)

var wc *clients.Wayback

func init() {
	wc = clients.NewWayback()
}

func HTTP(ctx context.Context, link models.Link) (models.Fix, error) {
	fix := models.Fix{Link: link}
	req, err := http.NewRequestWithContext(ctx, "GET", link.URL.String(), nil)
	if err != nil {
		return fix, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return fix, err
		}
		var httpErr *url.Error
		if errors.As(err, &httpErr) {
			switch {
			case httpErr.Timeout():
				fix.Reason = fmt.Sprintf("request timed out: %s", httpErr)
				fix.Severity = severity.Error
			case httpErr.Temporary():
				fix.Reason = fmt.Sprintf("url is temporarily unavailable: %s", httpErr)
				fix.Severity = severity.Error
			default:
				fix.Reason = fmt.Sprintf("couldn't establish a connection: %s", httpErr)
				fix.Severity = severity.Error
			}
		} else {
			return fix, err
		}
	}
	if err == nil && res.StatusCode == http.StatusOK {
		_ = res.Body.Close()
		// The URL redirects, suggest using the final destination instead of the original URL
		if link.URL.String() != res.Request.URL.String() {
			fix.Reason = "URL redirects, consider using the final destination instead"
			fix.Severity = severity.Warn
			fix.Suggestion = res.Request.URL.String()
			return fix, nil
		}
		// Try upgrading from http to https
		if link.URL.Scheme == "http" {
			slink := link.URL
			slink.Scheme = "https"
			res, err := http.Get(slink.String())
			if err == nil && res.StatusCode == http.StatusOK {
				_ = res.Body.Close()
				fix.Reason = "consider using https instead of http"
				fix.Severity = severity.Warn
				fix.Suggestion = slink.String()
			}
		}
		return fix, nil
	} else if err == nil && res.StatusCode != http.StatusOK {
		fix.Reason = fmt.Sprintf("url responded with status code: %d", res.StatusCode)
		fix.Severity = severity.Error
	}
	snap, err := wc.GetSnapshot(link.URL.String())
	if err != nil {
		return fix, nil
	}
	fix.Suggestion = snap.Closest.URL
	return fix, nil
}
