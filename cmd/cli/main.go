package main

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func main() {
	rootCmd := NewRootCmd()
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func initLogger(verbose bool) {
	log.SetLevel(log.ErrorLevel)
	if verbose {
		log.SetLevel(log.DebugLevel)
	}

	log.SetFormatter(&log.TextFormatter{
		ForceColors: true,
		// DisableColors:   false,
		TimestampFormat: "2006-01-02 15:04:05.999",
		FullTimestamp:   true,
	})
}

func run(args programArguments) error {
	initLogger(args.enableLogs)

	log.Infof("started")

	return nil
}
