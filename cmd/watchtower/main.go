package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
)

const (
	exitSuccess = 0
	exitFailure = 1
	exitUsage   = 2
)

var version = "dev"

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("watchtower", flag.ContinueOnError)
	flags.SetOutput(stderr)

	showVersion := flags.Bool("version", false, "print version information")
	var usageErr error
	flags.Usage = func() {
		_, usageErr = fmt.Fprintln(stderr, "Usage: watchtower [--version]")
	}

	if err := flags.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			if usageErr != nil {
				return exitFailure
			}

			return exitSuccess
		}

		return exitUsage
	}

	if flags.NArg() != 0 {
		if _, err := fmt.Fprintf(stderr, "watchtower: unexpected arguments: %v\n", flags.Args()); err != nil {
			return exitFailure
		}

		flags.Usage()
		if usageErr != nil {
			return exitFailure
		}

		return exitUsage
	}

	if *showVersion {
		if _, err := fmt.Fprintf(stdout, "watchtower %s\n", version); err != nil {
			return exitFailure
		}

		return exitSuccess
	}

	flags.Usage()
	if usageErr != nil {
		return exitFailure
	}

	return exitSuccess
}
