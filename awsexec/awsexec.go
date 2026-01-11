package awsexec

import (
	"context"
	"fmt"
)

// AWSCommandExecutor is the interface for executing AWS commands.
type Awsexec struct {
}

// New creates a new Awsexec instance.
func New(_ context.Context) *Awsexec {
	return &Awsexec{}
}

// Run executes the command.
func (a *Awsexec) Run(_ context.Context, args []string) error {
	fmt.Println("Hello, World!", args)
	return nil
}
