package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const (
	waybackURL = "https://archive.org/wayback"
)

// Wayback provides access to the Wayback machine API which can be used to obtain last
// working snapshot of the page.
type Wayback struct {
	c http.Client
}

// NewWayback creates new wayback client.
func NewWayback() *Wayback {
	return &Wayback{c: http.Client{
		Timeout: 15 * time.Second,
	}}
}

// ArchivedSnapshots represents the snapshots available through the Wayback machine.
// Currently we care only about the closest available snapshot.
type ArchivedSnapshots struct {
	Closest struct {
		Available bool   `json:"available"`
		URL       string `json:"url"`
		Timestamp string `json:"timestamp"`
		Status    string `json:"status"`
	} `json:"closest"`
}

// WaybackResponse is the response from Wayback machine API.
type WaybackResponse struct {
	ArchivedSnapshots ArchivedSnapshots `json:"archived_snapshots"`
}

// GetSnapshot obtains latest snapshot for given URL.
func (wc *Wayback) GetSnapshot(uri string) (*ArchivedSnapshots, error) {
	resp, err := wc.c.Get(fmt.Sprintf("%s/available?url=%s", waybackURL, url.QueryEscape(uri)))
	if err != nil {
		return nil, fmt.Errorf("get latest snapshot: %s", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode == http.StatusOK {
		var wbr WaybackResponse
		err := json.NewDecoder(resp.Body).Decode(&wbr)
		if err != nil {
			return nil, fmt.Errorf("decode snapshot: %s", err)
		}
		return &wbr.ArchivedSnapshots, nil
	}
	return nil, fmt.Errorf("archive responded with status code: %d", resp.StatusCode)
}
