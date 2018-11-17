package main

import (
	"context"
	"os"
	"os/exec"
	"time"
)

type subcommand struct {
	name string
	path string
}

func (s *subcommand) usage() (string, error) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelFn()

	cmd := exec.CommandContext(ctx, s.path, "--short-help")
	out, err := cmd.Output()

	return string(out), err
}

func (s *subcommand) exec(args []string) error {
	cmd := exec.Command(s.path, args...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
