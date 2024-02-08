package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
)

type InstallOptions struct {
	Interactive       bool
	Simulate          bool
	Update            bool
	CleanPackageCache bool
	CleanLists        bool
	PackageList       []string
	RequirementsFile  string
	QuietLevel        int
	AptListsDir       string
	ExtraOptions      string
}

func AptInstall(args InstallOptions) error {
	return AptInstallContext(context.TODO(), args)
}

func AptInstallContext(ctx context.Context, opts InstallOptions) (err error) {
	reqs := append([]string{}, opts.PackageList...)
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

	err = aptGetInstall(ctx, reqs, opts.Simulate, opts.Interactive, opts.QuietLevel, opts.ExtraOptions)
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

func readRequirementsFile(filename string) ([]string, error) {
	fp, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("could not read requirements file: %w", err)
	}
	defer fp.Close()

	sc := bufio.NewScanner(fp)
	sc.Split(bufio.ScanLines)
	requirements := []string{}
	for sc.Scan() {
		req, _, _ := strings.Cut(sc.Text(), "#")
		req = strings.TrimSpace(req)
		if req == "" {
			continue
		}
		requirements = append(requirements, req)
	}
	return requirements, nil
}
