package log

import (
	"github.com/apex/log"
	"github.com/apex/log/handlers/text"
)

func New(verbose bool) log.Interface {
	log.SetHandler(text.Default)
	log.SetLevel(log.ErrorLevel)
	if verbose {
		log.SetLevel(log.InfoLevel)
	}
	return log.Log
}
