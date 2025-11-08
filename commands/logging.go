package commands

import (
	"github.com/osteele/gojekyll/logger"
	"github.com/osteele/gojekyll/utils"
)

// log is the logger instance used by commands
var log = logger.Default()

type bannerLogger struct{}

var bannerLog = bannerLogger{}

func (l *bannerLogger) Info(a ...interface{}) {
	log.Println(a...)
}

func (l *bannerLogger) label(label string, msg string, a ...interface{}) {
	if !quiet {
		log.Label(label, msg, a...)
	}
}

func (l *bannerLogger) path(label string, filename string) {
	//nolint:govet
	l.label(label, utils.MustAbs(filename))
}
