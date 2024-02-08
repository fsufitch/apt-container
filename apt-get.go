package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/kballard/go-shellquote"
)

type AptGet struct {
	ExecPath string
	Cwd      string
	Args     []string
}

func (aptget AptGet) Executable() string {
	if aptget.ExecPath != "" {
		return aptget.ExecPath
	}
	return "apt-get"
}

func (aptget AptGet) String() string {
	return shellquote.Join(append([]string{aptget.Executable()}, aptget.Args...)...)
}

func (aptget AptGet) Run() {
	aptget.RunContext(context.TODO())
}

func (aptget AptGet) RunContext(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, aptget.Executable(), aptget.Args...)
	cmd.Dir = aptget.Cwd
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	if errors.Is(err, exec.ErrNotFound) {
		return fmt.Errorf("%w: %w", ErrAptGetNotFound, err)
	}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		err = AptGetError{exitErr}
	}

	return err
}

func aptGetArgs(simulate bool, interactive bool, quietLevel int) []string {
	args := []string{}
	if simulate {
		args = append(args, "-s")
	}
	if !interactive {
		args = append(args, "-y")
	}
	qs := strings.Join(make([]string, quietLevel+1), "q")
	if qs != "" {
		args = append(args, "-"+qs)
	}

	return args
}

func printFIfNotSuperQuiet(quietLevel int, format string, vals ...interface{}) {
	if quietLevel < 2 {
		fmt.Printf(format, vals...)
	}
}

func aptGetClean(ctx context.Context, simulate bool, interactive bool, quietLevel int) error {
	cmd := AptGet{Args: []string{"clean"}}
	cmd.Args = append(cmd.Args, aptGetArgs(simulate, interactive, quietLevel)...)
	printFIfNotSuperQuiet(quietLevel, " + %s\n", cmd.String())

	if err := cmd.RunContext(ctx); err != nil {
		return fmt.Errorf("failed running apt-get %s: %w", cmd.Args[0], err)
	}
	return nil
}

func aptGetInstall(ctx context.Context, packages []string, simulate bool, interactive bool, quietLevel int, extraOptions string) error {
	cmd := AptGet{Args: []string{"install"}}
	cmd.Args = append(cmd.Args, aptGetArgs(simulate, interactive, quietLevel)...)
	extraOptionsArr, err := shellquote.Split(extraOptions)
	if err != nil {
		return fmt.Errorf("failed parsing extra options: %w", err)
	}
	cmd.Args = append(cmd.Args, extraOptionsArr...)
	cmd.Args = append(cmd.Args, packages...)
	printFIfNotSuperQuiet(quietLevel, " + %s\n", cmd.String())

	if err := cmd.RunContext(ctx); err != nil {
		return fmt.Errorf("failed running apt-get %s: %w", cmd.Args[0], err)
	}
	return nil
}

func aptGetSatisfy(ctx context.Context, requirements []string, simulate bool, interactive bool, quietLevel int, extraOptions string) error {
	cmd := AptGet{Args: []string{"satisfy"}}
	cmd.Args = append(cmd.Args, aptGetArgs(simulate, interactive, quietLevel)...)
	extraOptionsArr, err := shellquote.Split(extraOptions)
	if err != nil {
		return fmt.Errorf("failed parsing extra options: %w", err)
	}
	cmd.Args = append(cmd.Args, extraOptionsArr...)
	cmd.Args = append(cmd.Args, requirements...)
	printFIfNotSuperQuiet(quietLevel, " + %s\n", cmd.String())

	if err := cmd.RunContext(ctx); err != nil {
		return fmt.Errorf("failed running apt-get %s: %w", cmd.Args[0], err)
	}
	return nil
}

func aptGetUpdate(ctx context.Context, interactive bool, quietLevel int) error {
	cmd := AptGet{Args: []string{"update"}}
	cmd.Args = append(cmd.Args, aptGetArgs(false, interactive, quietLevel)...)
	printFIfNotSuperQuiet(quietLevel, " + %s\n", cmd.String())

	if err := cmd.RunContext(ctx); err != nil {
		return fmt.Errorf("failed running apt-get %s: %w", cmd.Args[0], err)
	}
	return nil
}

func cleanAptLists(dir string, quietLevel int) error {
	if dir == "" {
		dir = "/var/lib/apt/lists"
	}
	printFIfNotSuperQuiet(quietLevel, " + removing lists: %s\n", dir)
	err := os.RemoveAll(dir)
	if err != nil {
		return fmt.Errorf("failed removing lists (%s): %w", dir, err)
	}
	return nil
}
