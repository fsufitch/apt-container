package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/urfave/cli/v2"
)

var cliCommonInstallSatisfyArgs = []cli.Flag{
	&cli.BoolFlag{
		Name:  "interactive",
		Usage: "run interactively (include -y in apt-get commands)",
	},
	&cli.BoolFlag{
		Name:    "simulate",
		Aliases: []string{"s"},
		Usage:   "run in simulated mode (include -s in apt-get commands)",
	},
	&cli.BoolFlag{
		Name:    "no-update",
		Aliases: []string{"U"},
		Usage:   "skip running 'apt-get update' before installing",
	},
	&cli.BoolFlag{
		Name:    "keep-cache",
		Aliases: []string{"C"},
		Usage:   "keep package cache after install; don't run 'apt-get clean'",
	},
	&cli.BoolFlag{
		Name:    "keep-lists",
		Aliases: []string{"L"},
		Usage:   "keep package lists used for install",
	},
	&cli.StringFlag{
		Name:  "apt-lists-dir",
		Usage: "directory that apt keeps its stuff",
		Value: "/var/lib/apt/lists",
	},
	&cli.StringFlag{
		Name:    "extra-options",
		Aliases: []string{"o"},
		Usage:   "options to passthrough to the apt-get command",
	},
	&cli.IntFlag{
		Name:    "quiet",
		Aliases: []string{"q"},
		Value:   1,
		Usage:   "how many '-q' to pass to apt-get",
	},
}

var cliInstallArgs = append(
	cliCommonInstallSatisfyArgs,
	[]cli.Flag{
		&cli.PathFlag{
			Name:    "requirements",
			Aliases: []string{"r"},
			Usage:   "read package list from `FILE` (one per line, # is a comment); mutually exclusive with packages as arguments",
		},
	}...,
)

func parseInstallArgs(ctx *cli.Context) (*InstallOptions, error) {
	reqsFile := ctx.Path("requirements")
	packageList := ctx.Args().Slice()

	if reqsFile != "" && len(packageList) > 0 {
		return nil, fmt.Errorf("%w; cannot specify packages through both file and arguments", ErrInvalidOptions)
	}

	args := InstallOptions{
		Simulate:          ctx.Bool("simulate"),
		Update:            !ctx.Bool("no-update"),
		CleanPackageCache: !ctx.Bool("keep-cache"),
		CleanLists:        !ctx.Bool("keep-lists"),
		PackageList:       packageList,
		RequirementsFile:  reqsFile,
		ExtraOptions:      ctx.String("extra-options"),
		QuietLevel:        ctx.Int("quiet"),
	}

	return &args, nil
}

func createInstallCommand(mainCtx context.Context, installFunc func(context.Context, InstallOptions) error) *cli.Command {
	return &cli.Command{
		Name:        "install",
		Aliases:     []string{"i"},
		Description: "container-friendly version of 'apt-get install' (see apt-get man pages)",
		ArgsUsage:   "packages...",
		Flags:       cliInstallArgs,
		Action: func(ctx *cli.Context) error {
			args, err := parseInstallArgs(ctx)

			if errors.Is(err, ErrInvalidOptions) {
				cli.ShowSubcommandHelp(ctx)
			}
			if err != nil {
				return err
			}

			err = installFunc(mainCtx, *args)
			if errors.Is(err, ErrInvalidOptions) {
				cli.ShowSubcommandHelp(ctx)
			}
			return err
		},
	}
}
