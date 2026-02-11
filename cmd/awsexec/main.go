package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/ikedam/awsexec/awsexec"
)

var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	ctx := context.Background()

	var showVersion bool
	flag.BoolVar(&showVersion, "version", false, "Show version information")
	flag.BoolVar(&showVersion, "v", false, "Show version information (shorthand)")
	flag.Parse()

	if showVersion {
		fmt.Printf("awsexec version %s\n", version)
		fmt.Printf("commit: %s\n", commit)
		fmt.Printf("date: %s\n", date)
		os.Exit(0)
	}

	awsexec := awsexec.New(ctx)
	runArgs := buildRunArgs(os.Args, flag.Args())
	err := awsexec.Run(ctx, runArgs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// buildRunArgs determines which arguments to pass to Awsexec.Run.
//
// By default, returns flag.Args() as-is. However, when flag.Args() corresponds
// to the tail of os.Args and the argument immediately before that tail is "--",
// returns that "--" plus the tail so that "awsexec -- command" works when
// AWS_PROFILE is set.
//
// Example:
//
//	os.Args  = ["awsexec", "--", "echo", "hello"]
//	flagArgs = ["echo", "hello"]
//	-> argument before tail is "--", so return ["--", "echo", "hello"].
//
//	os.Args  = ["awsexec", "myprofile", "--", "echo", "hello"]
//	flagArgs = ["myprofile", "--", "echo", "hello"]
//	-> argument before tail is program name, so return flagArgs as-is.
func buildRunArgs(osArgs, flagArgs []string) []string {
	if len(osArgs) < len(flagArgs) {
		return flagArgs
	}

	start := len(osArgs) - len(flagArgs)
	// os.Args[0] is the program name; only consider indices >= 1.
	if start-1 >= 1 && osArgs[start-1] == "--" {
		return osArgs[start-1:]
	}

	return flagArgs
}
