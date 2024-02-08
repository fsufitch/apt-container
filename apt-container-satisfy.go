package main

import (
	"context"
	"fmt"
	"os"
)

type SatisfyOptions struct {
	Interactive       bool
	Simulate          bool
	Update            bool
	CleanPackageCache bool
	CleanLists        bool
	Requirements      []string
	RequirementsFile  string
	QuietLevel        int
	AptListsDir       string
	ExtraOptions      string
}

func AptSatisfy(opts SatisfyOptions) error {
	return AptSatisfyContext(context.TODO(), opts)
}

func AptSatisfyContext(ctx context.Context, opts SatisfyOptions) (err error) {
	reqs := append([]string{}, opts.Requirements...)
	if opts.RequirementsFile != "" {
		var reqsFromFile []string
		reqsFromFile, err = readRequirementsFile(opts.RequirementsFile)
		if err == nil {
			reqs = append(reqs, reqsFromFile...)
		}
	}
	if err != nil {
		return
	}

	if len(reqs) < 1 {
		return ErrNoInputPackages
	}

	if opts.Update {
		if opts.Simulate {
			fmt.Fprintln(os.Stderr, "NOTE: simulated mode; not updating lists")
		} else {
			err = aptGetUpdate(ctx, opts.Interactive, opts.QuietLevel)
		}
	}
	if err != nil {
		return
	}

	err = aptGetSatisfy(ctx, reqs, opts.Simulate, opts.Interactive, opts.QuietLevel, opts.ExtraOptions)
	if err != nil {
		return
	}

	if opts.CleanPackageCache {
		aptGetClean(ctx, opts.Simulate, opts.Interactive, opts.QuietLevel)
	}
	if err != nil {
		return
	}

	if opts.CleanLists {
		if opts.Simulate {
			fmt.Fprintln(os.Stderr, "NOTE: simlated mode; not cleaning apt lists")
		} else {
			err = cleanAptLists(opts.AptListsDir, opts.QuietLevel)
		}
	}
	if err != nil {
		return
	}

	return nil

}
