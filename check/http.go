package check

import (
    "errors"
    "fmt"
    "net/http"
    "net/url"

    "github.com/matoous/linkfix/internal/clients"
    "github.com/matoous/linkfix/models"
)

var wc *clients.Wayback

func init() {
    wc = clients.NewWayback()
}

func HTTP(link models.Link) (models.Fix, error) {
    fix := models.Fix{Link: link}
    res, err := http.Get(link.URL.String())
    if err == nil && res.StatusCode == http.StatusOK {
        _ = res.Body.Close()
        // The URL redirects, suggest using the final destination instead of the original URL
        if link.URL.String() != res.Request.URL.String() {
            fix.Reason = "URL redirects, consider using the final destination instead"
            fix.Severity = "warn"
            fix.Suggestion = res.Request.URL.String()
            return fix, nil
        }
        // Try upgrading from http to https
        if link.URL.Scheme == "http" {
            link.URL.Scheme = "https"
            if res, err := http.Get(link.URL.String()); err == nil && res.StatusCode == http.StatusOK {
                fix.Reason = "consider using https instead of http"
                fix.Severity = "warn"
                fix.Suggestion = link.URL.String()
            }
            return fix, nil
        }
    }
    if err != nil {
        var httpErr *url.Error
        if errors.As(err, &httpErr) {
            if httpErr.Temporary() {
                fix.Reason = fmt.Sprintf("url is temporarily unavailable: %s", httpErr)
                fix.Severity = "error"
            } else {
                fix.Reason = fmt.Sprintf("couldn't establish a connection: %s", httpErr)
                fix.Severity = "error"
            }
        }
    } else if res.StatusCode != http.StatusOK {
        fix.Reason = fmt.Sprintf("url responded with status code: %d", res.StatusCode)
        fix.Severity = "error"
    }
    snap, err := wc.GetSnapshot(link.URL.String())
    if err != nil {
        return fix, err
    }
    fix.Suggestion = snap.Closest.URL
    return fix, nil
}
