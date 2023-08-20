package main

import (
	"context"
	"io/fs"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

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

	// catch interrupt signals
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	generateDiff(ctx, args.dirPathA, args.dirPathB)

	// restore default behavior to interrupt signals
	stop()

	return nil
}

func generateDiff(ctx context.Context, folderA, folderB string) error {
	log.Infof("Traversing folderA: %s\n", folderA)

	err := filepath.Walk(folderA, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			log.Errorf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}

		if ctx.Err() != nil {
			log.Debugf("Signaled to stop, reached path: %s", path)
			return fs.SkipAll
		}

		log.Infof("visited file or dir: %q\n", path)

		return nil
	})

	if err != nil {
		log.Errorf("error walking the path %q: %v\n", folderA, err)
		return err
	}

	return nil
}
