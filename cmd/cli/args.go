package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type programArguments struct {
	dirPathA   string
	dirPathB   string
	enableLogs bool
}

const (
	version    = "1.0.0"
	binaryName = "fdiff"
)

var args programArguments

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     fmt.Sprintf("%s <folderA> <folderB>", binaryName),
		Version: version,
		Long:    fmt.Sprintf("%s - Recursively traverse folderA and folderB and prints the files that are in folderA and not in folderB (set difference folderA - folderB)", binaryName),
		Short:   fmt.Sprintf("%s - print folderA - folderB file difference", binaryName),
		Args:    cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, cmdlineArgs []string) {

			args.dirPathA = cmdlineArgs[0]
			args.dirPathB = cmdlineArgs[1]

			err := run(args)
			if err != nil {
				os.Exit(1)
			}
		},
	}

	rootCmd.Flags().BoolVarP(&args.enableLogs, "verbose", "V", false, "enable verbose logging")

	return rootCmd
}
