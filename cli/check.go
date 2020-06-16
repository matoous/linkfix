package cli

import (
    "context"
    "fmt"
    "time"

    "github.com/spf13/cobra"
    "golang.org/x/sync/errgroup"

    "github.com/matoous/linkfix/check"
    "github.com/matoous/linkfix/internal/log"
    "github.com/matoous/linkfix/internal/printers"
    "github.com/matoous/linkfix/internal/stats"
    "github.com/matoous/linkfix/models"
    "github.com/matoous/linkfix/source"
)

func checkCmd() *cobra.Command {
    var verbose, noStats bool
    var workers int
    var format string
    // exclude is for files, ignore is for URLs
    var exclude, ignore []string

    cmd := &cobra.Command{
        Use:   "check",
        Short: "Checks the links but doesn't apply or suggest any changes",
        Run: func(cmd *cobra.Command, args []string) {

            logger := log.New(verbose)
            stat := stats.New()

            logger.Infof("🔗 starting linkcheck %s", time.Now().Format(time.RFC3339))

            logger.Info("looking for links")
            src := source.Filesystem(args[0], ignore)
            links, err := src.List()
            if err != nil {
                logger.Errorf("listing links from the source: %s", err)
                return
            }
            logger.Infof("found %d links", len(links))
            stat.Queue(len(links))

            checker := check.NewChecker(exclude)
            toProcess := make(chan models.Link, 2048)
            processed := make(chan models.Fix)
            stopChan := make(chan struct{})

            logger.Infof("starting %d workers", workers)
            g, ctx := errgroup.WithContext(context.Background())
            // This gorutine lists the links from source.
            g.Go(func() error {
                defer close(toProcess)
                for i := range links {
                    toProcess <- links[i]
                }
                return nil
            })
            // This starts gorutines that process found links from the source, check whether the links is available
            // and if necessary try to find a replacement on wayback machine.
            for i := 0; i < workers; i++ {
                g.Go(func() error {
                    for {
                        select {
                        case <-ctx.Done():
                            return ctx.Err()
                        case l, ok := <-toProcess:
                            if !ok {
                                return nil
                            }
                            stat.Process()
                            res, err := checker.ProcessUrl(l)
                            if err != nil {
                                return fmt.Errorf("process %q in file %s: %s", l.URL, l.Path, err)
                            }
                            if res.Reason != "" {
                                stat.Failed()
                                logger.Infof("  ❌  %s", res.URL)
                                processed <- res
                            } else {
                                stat.OK()
                                logger.Infof("  ✅  %s", res.URL)
                            }
                        }
                    }
                })
            }
            // This gorutine makes sure that when all workers are done the channel with processed links is closed
            // and stop is submitted for the gorutine that reports statistics.
            go func() {
                _ = g.Wait()
                close(processed)
                stopChan <- struct{}{}
            }()

            // This gorutine periodically prints statistics.
            if !noStats {
                go func() {
                    ticker := time.Tick(time.Second * 5)
                    for {
                        select {
                        case <-ticker:
                            logger.Infof(stat.String())
                        case <-stopChan:
                            return
                        }
                    }
                }()
            }

            // Collect all processed links.
            var results []models.Fix
            for l := range processed {
                results = append(results, l)
            }

            if err := g.Wait(); err != nil {
                logger.Errorf("process links: %s", err)
            }

            p := printers.GetPrinter(format)
            err = p.PrintAll(results)
            if err != nil {
                logger.Errorf("print results: %s", err)
            }
        },
    }

    cmd.Flags().IntVarP(&workers, "workers", "w", 10, "number of workers for processing the links")
    cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "run in verbose mode")
    cmd.Flags().BoolVarP(&noStats, "no-stats", "n", false, "disable runtime statistics")
    cmd.Flags().StringSliceVarP(&exclude, "exclude-links", "e", []string{}, "websites to exclude from the link checking")
    cmd.Flags().StringSliceVarP(&ignore, "ignore-paths", "i", []string{}, "file path patterns to ignore")
    cmd.Flags().StringVarP(&format, "format", "f", "text", "output format, one of: text|json|github")

    return cmd
}
