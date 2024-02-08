package main

import (
	"errors"
	"fmt"
	"os/exec"
)

type AptGetError struct {
	*exec.ExitError
}

func (err AptGetError) Unwrap() error {
	return err.ExitError
}

var ErrInvalidOptions = errors.New("invalid options")
var ErrNoInputPackages = fmt.Errorf("%w: received no packages/requirements", ErrInvalidOptions)
var ErrAptGetNotFound = errors.New("apt-get command not found")
var ErrSubprocessCanceled = errors.New("subprocess was canceled")
