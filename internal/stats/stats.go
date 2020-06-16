package stats

import (
	"strconv"
	"strings"
	"sync"

	"github.com/fatih/color"
)

type Stats struct {
	sync.Mutex
	InQueue   int
	InProcess int
	NOK       int
	NFailed   int
}

func New() *Stats {
	return &Stats{}
}

func (s *Stats) Queue(n int) {
	s.Lock()
	defer s.Unlock()
	s.InQueue = n
}

func (s *Stats) Process() {
	s.Lock()
	defer s.Unlock()
	s.InQueue--
	s.InProcess++
}

func (s *Stats) OK() {
	s.Lock()
	defer s.Unlock()
	s.InProcess--
	s.NOK++
}

func (s *Stats) Failed() {
	s.Lock()
	defer s.Unlock()
	s.InProcess--
	s.NFailed++
}

func (s *Stats) String() string {
	s.Lock()
	defer s.Unlock()
	var b strings.Builder
	b.WriteString(strconv.Itoa(s.InQueue))
	b.WriteString(" in queue, ")
	b.WriteString(strconv.Itoa(s.InProcess))
	b.WriteString(" being processed, ")
	b.WriteString(strconv.Itoa(s.NOK + s.NFailed))
	b.WriteString(" (")
	b.WriteString(color.New(color.FgGreen).Sprintf("%d", s.NOK))
	b.WriteString("+")
	b.WriteString(color.New(color.FgRed).Sprintf("%d", s.NFailed))
	b.WriteString(") done")
	return b.String()
}
