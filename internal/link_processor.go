package internal

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/apex/log"
	"golang.org/x/sync/errgroup"

	"github.com/matoous/linkfix/check"
	"github.com/matoous/linkfix/internal/callonce"
	"github.com/matoous/linkfix/internal/cfg"
	"github.com/matoous/linkfix/internal/stats"
	"github.com/matoous/linkfix/models"
)

type LinkProcessor struct {
	sync.RWMutex
	log log.Interface

	checker *check.Checker
	stats   *stats.Stats

	cache callonce.Group

	// settings
	workers int
	noStats bool
	exclude []string
}

func NewLinkProcessor(c *cfg.Config, l log.Interface) *LinkProcessor {
	return &LinkProcessor{
		log:     l,
		workers: c.Workers,
		noStats: c.NoStats,
		exclude: c.Exclude,
		stats:   stats.New(),
		checker: check.NewChecker(c.Exclude),
	}
}

func (lp *LinkProcessor) Process(ctx context.Context, links []models.Link) []models.Fix {
	lp.log.Infof("found %d links", len(links))
	lp.stats.Queue(len(links))

	toProcess := make(chan models.Link, 2048)
	processed := make(chan models.Fix)

	go func() {
		defer close(toProcess)
		for i := range links {
			toProcess <- links[i]
		}
	}()

	lp.log.Infof("starting %d workers", lp.workers)

	// This starts gorutines that process found links from the source, check whether the links is available
	// and if necessary try to find a replacement on wayback machine.
	g, ctx := errgroup.WithContext(ctx)
	for i := 0; i < lp.workers; i++ {
		g.Go(func() error {
			for {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case l, ok := <-toProcess:
					if !ok {
						return nil
					}
					val, pErr := lp.cache.Do(ctx, l.URL.String(), func(ctx context.Context) (interface{}, error) {
						return lp.checker.ProcessURL(ctx, l)
					})
					if pErr != nil {
						return fmt.Errorf("process %q in file %s: %s", l.URL, l.Path, pErr)
					}
					res := val.(models.Fix)
					if res.Reason != "" {
						lp.stats.Failed()
						lp.log.Infof("  ❌  %s", res.URL)
						processed <- res
					} else {
						lp.stats.OK()
						lp.log.Infof("  ✅  %s", res.URL)
					}
				}
			}
		})
	}

	if !lp.noStats {
		stopChan := make(chan struct{})
		go func() {
			_ = g.Wait()
			stopChan <- struct{}{}
		}()
		go lp.printStats(stopChan)
	}

	go func() {
		_ = g.Wait()
		close(processed)
	}()

	// Collect all processed links.
	var results []models.Fix
	for l := range processed {
		results = append(results, l)
	}

	if wErr := g.Wait(); wErr != nil {
		lp.log.Errorf("process links: %s", wErr)
	}

	return results
}

func (lp *LinkProcessor) printStats(stopChan chan struct{}) {
	ticker := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-ticker.C:
			lp.log.Infof(lp.stats.String())
		case <-stopChan:
			ticker.Stop()
			return
		}
	}
}
