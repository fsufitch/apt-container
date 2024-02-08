package main

import (
	"context"
	"errors"

	"github.com/urfave/cli/v2"
)

var cliSatisfyArgs = append(
	cliCommonInstallSatisfyArgs,
	[]cli.Flag{
		&cli.PathFlag{
			Name:    "requirements",
			Aliases: []string{"r"},
			Usage:   "read dependency list from `FILE` (one per line, # is a comment); mutually exclusive with packages as arguments",
		},
	}...,
)

func parseSatisfyArgs(ctx *cli.Context) (*SatisfyOptions, error) {
	reqsFile := ctx.Path("requirements")
	requirements := ctx.Args().Slice()

	if reqsFile != "" && len(requirements) > 0 {
		return nil, errors.New("not allowed to use both a requirements file and package arguments")
	}

	args := SatisfyOptions{
		Simulate:          ctx.Bool("simulate"),
		Update:            !ctx.Bool("no-update"),
		CleanPackageCache: !ctx.Bool("keep-cache"),
		CleanLists:        !ctx.Bool("keep-lists"),
		Requirements:      requirements,
		RequirementsFile:  reqsFile,
		ExtraOptions:      ctx.String("extra-options"),
	}

	return &args, nil
}

func createSatisfyCommand(mainCtx context.Context, satisfyFunc func(context.Context, SatisfyOptions) error) *cli.Command {
	return &cli.Command{
		Name:        "satisfy",
		Aliases:     []string{"s"},
		Description: "container-friendly version of 'apt-get satisfy'; uses dependency strings (see apt-get man pages)",
		ArgsUsage:   "dependencies...",
		Flags:       cliSatisfyArgs,
		Action: func(ctx *cli.Context) error {
			args, err := parseSatisfyArgs(ctx)
			if err != nil {
				return err
			}
			return satisfyFunc(mainCtx, *args)
		},
	}
}
