package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"

	"github.com/urfave/cli/v2"
)

func withGracefulSigInt(parent context.Context) context.Context {
	ctx, cancel := context.WithCancel(parent)

	sigIntCh := make(chan os.Signal, 1)
	signal.Notify(sigIntCh, os.Interrupt)
	go func() {
		<-sigIntCh
		fmt.Fprintln(os.Stderr, "SIGINT received; stopping gracefully")
		signal.Stop(sigIntCh)
		cancel()
	}()

	return ctx
}

func main() {
	ctx := context.Background()
	ctx = withGracefulSigInt(ctx)

	var cliApp = &cli.App{
		Name:  "apt-container",
		Usage: "wrapper around apt-get tools, better suited for use in container building",
		Commands: []*cli.Command{
			createInstallCommand(ctx, AptInstallContext),
			createSatisfyCommand(ctx, AptSatisfyContext),
		},
		EnableBashCompletion:   true,
		UseShortOptionHandling: true,
	}
	err := cliApp.Run(os.Args)

	if err == nil {
		return
	}
	exitCode := 1

	var aptWrappedErr AptGetError
	if errors.As(err, &aptWrappedErr) {
		exitCode = aptWrappedErr.ExitCode()
	}

	fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
	os.Exit(exitCode)
}
