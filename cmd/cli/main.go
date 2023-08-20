package main

import (
	"context"
	"fmt"
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

	folderAContents, err := listAllFiles(ctx, folderA)
	if err != nil {
		log.Errorf("failed to traverse folderA '%s': %s\n", folderA, err)
		return err
	}

	log.Infof("Traversing folderB: %s\n", folderB)

	folderBContents, err := listAllFiles(ctx, folderB)
	if err != nil {
		log.Errorf("failed to traverse folderA '%s': %s\n", folderB, err)
		return err
	}

	log.Infof("folderA has %d files, folderB has %d files", len(folderAContents), len(folderBContents))

	for k := range folderBContents {
		delete(folderAContents, k)
	}

	fmt.Printf("List of files from '%s' not present in '%s' (%d files)\n", folderA, folderB, len(folderAContents))
	for k := range folderAContents {
		fmt.Println(k)
	}

	return nil
}

type fileInfo struct {
	path  string
	isDir bool
	size  int64
}

func listAllFiles(ctx context.Context, rootDir string) (map[string]fileInfo, error) {
	dirContents := make(map[string]fileInfo)

	var err error
	rootDir, err = filepath.Abs(rootDir)

	if err != nil {
		log.Errorf("failed to get absolute path for rootDir '%s': %s", rootDir, err)
		return nil, err
	}

	err = filepath.Walk(rootDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			log.Errorf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}

		if ctx.Err() != nil {
			log.Debugf("Signaled to stop, reached path: %s", path)
			return fs.SkipAll
		}

		pathKey := path
		if info.IsDir() {
			pathKey += "/"
		}

		dirContents[pathKey] = fileInfo{path, info.IsDir(), info.Size()}

		return nil
	})

	// we don't need the rootDir entry, it doesn't make sense since it will always be different
	delete(dirContents, rootDir)
	delete(dirContents, rootDir+"/")

	if err != nil {
		return nil, err
	}

	return dirContents, err
}
