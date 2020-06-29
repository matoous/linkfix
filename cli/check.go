package cli

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/matoous/linkfix/internal"
	"github.com/matoous/linkfix/internal/cfg"
	"github.com/matoous/linkfix/internal/log"
	"github.com/matoous/linkfix/internal/printers"
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
			ctx, cancel := context.WithCancel(context.Background())
			go func() {
				c := make(chan os.Signal, 1)
				signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
				<-c
				cancel()
			}()

			logger := log.New(verbose)
			logger.Infof("🔗 starting linkcheck %s", time.Now().Format(time.RFC3339))
			src := source.Filesystem(args[0], ignore)

			logger.Info("looking for links")
			links, err := src.List()
			if err != nil {
				logger.Errorf("listing links from the source: %s", err)
				return
			}

			processor := internal.NewLinkProcessor(&cfg.Config{
				Workers: workers,
				NoStats: noStats,
				Exclude: exclude,
			}, logger)
			results := processor.Process(ctx, links)

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
