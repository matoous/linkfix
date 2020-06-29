package stats

import (
	"strconv"
	"strings"
	"sync"

	"github.com/fatih/color"
)

// Stats handles linting statistics during the run.
// It can be used to periodically inform the user on progress of the linting.
type Stats struct {
	sync.Mutex
	// InQueue is number of links waiting in the queue to be checked.
	InQueue int
	// NOK is number of links that were checked and are OK.
	NOK int
	// NFailed is number of links that have some issue of severity warn or error.
	NFailed int
}

//
func New() *Stats {
	return &Stats{}
}

// Queue adds `n` links into the processing queue.
func (s *Stats) Queue(n int) {
	s.Lock()
	defer s.Unlock()
	s.InQueue = n
}

// OK finishes link statistic by marking the link with ok status.
func (s *Stats) OK() {
	s.Lock()
	defer s.Unlock()
	s.InQueue--
	s.NOK++
}

// Failed finishes link statistic by marking the link with fail status.
func (s *Stats) Failed() {
	s.Lock()
	defer s.Unlock()
	s.InQueue--
	s.NFailed++
}

// String returns string representation of the stats.
func (s *Stats) String() string {
	s.Lock()
	defer s.Unlock()
	var b strings.Builder
	b.WriteString(strconv.Itoa(s.InQueue))
	b.WriteString(" in queue, ")
	b.WriteString(strconv.Itoa(s.NOK + s.NFailed))
	b.WriteString(" (")
	b.WriteString(color.New(color.FgGreen).Sprintf("%d", s.NOK))
	b.WriteString("+")
	b.WriteString(color.New(color.FgRed).Sprintf("%d", s.NFailed))
	b.WriteString(") done")
	return b.String()
}
