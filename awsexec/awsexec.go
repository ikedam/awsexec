package awsexec

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// Credentials represents AWS credentials exported from aws configure export-credentials
type Credentials struct {
	Version         int    `json:"Version"`
	AccessKeyID     string `json:"AccessKeyId"`
	SecretAccessKey string `json:"SecretAccessKey"`
	SessionToken    string `json:"SessionToken"`
	Expiration      string `json:"Expiration"`
}

// CredentialExporter is an interface for exporting AWS credentials.
type CredentialExporter interface {
	ExportCredentials(ctx context.Context, profile string) (*Credentials, error)
}

// Awsexec is the main executor for AWS commands with exported credentials
type Awsexec struct {
	exporter CredentialExporter
}

// New creates a new Awsexec instance with the default credential exporter.
func New(_ context.Context) *Awsexec {
	return &Awsexec{
		exporter: NewAWSCommandExporter(),
	}
}

// NewWithExporter creates a new Awsexec instance with the specified credential exporter.
func NewWithExporter(exporter CredentialExporter) *Awsexec {
	return &Awsexec{
		exporter: exporter,
	}
}

// Run executes the command with AWS credentials exported as environment variables.
func (a *Awsexec) Run(ctx context.Context, args []string) error {
	// Parse arguments to extract profile and command
	profile, command, err := a.parseArgs(args)
	if err != nil {
		return fmt.Errorf("failed to parse arguments: %w", err)
	}

	// Export credentials using credential exporter
	creds, err := a.exporter.ExportCredentials(ctx, profile)
	if err != nil {
		return fmt.Errorf("failed to export credentials: %w", err)
	}

	// Execute the command with credentials set as environment variables
	return a.executeCommand(ctx, creds, command)
}

// parseArgs parses command line arguments to extract profile and command.
// Format: awsexec [profile] -- command [args...]
// Or: AWS_PROFILE=profile awsexec -- command [args...]
func (a *Awsexec) parseArgs(args []string) (string, []string, error) {
	// Find the -- separator
	separatorIndex := -1
	for i, arg := range args {
		if arg == "--" {
			separatorIndex = i
			break
		}
	}

	if separatorIndex == -1 {
		return "", nil, fmt.Errorf("missing '--' separator")
	}

	// Extract profile from arguments or environment variable
	var profile string
	if separatorIndex > 0 {
		// Profile is specified as an argument
		profile = args[0]
	} else {
		// Profile should be in AWS_PROFILE environment variable
		profile = os.Getenv("AWS_PROFILE")
		if profile == "" {
			return "", nil, fmt.Errorf("AWS_PROFILE environment variable is required when profile is not specified as an argument")
		}
	}

	// Extract command after --
	if separatorIndex+1 >= len(args) {
		return "", nil, fmt.Errorf("no command specified after '--'")
	}
	command := args[separatorIndex+1:]

	return profile, command, nil
}

// executeCommand executes the specified command with AWS credentials set as environment variables.
// This uses syscall.Exec to replace the current process with the command, so it does not return on success.
func (a *Awsexec) executeCommand(_ context.Context, creds *Credentials, command []string) error {
	if len(command) == 0 {
		return fmt.Errorf("command is empty")
	}

	// Resolve the command path
	cmdPath, err := exec.LookPath(command[0])
	if err != nil {
		return fmt.Errorf("command not found: %w", err)
	}

	// Prepare environment variables
	env := os.Environ()
	env = append(env, fmt.Sprintf("AWS_ACCESS_KEY_ID=%s", creds.AccessKeyID))
	env = append(env, fmt.Sprintf("AWS_SECRET_ACCESS_KEY=%s", creds.SecretAccessKey))
	if creds.SessionToken != "" {
		env = append(env, fmt.Sprintf("AWS_SESSION_TOKEN=%s", creds.SessionToken))
	}
	if creds.Expiration != "" {
		env = append(env, fmt.Sprintf("AWS_CREDENTIAL_EXPIRATION=%s", creds.Expiration))
	}

	// Prepare arguments (first argument should be the command name)
	args := command

	// Execute the command using syscall.Exec
	// This replaces the current process, so it does not return on success.
	// On error, it returns an error.
	err = syscall.Exec(cmdPath, args, env)
	if err != nil {
		return fmt.Errorf("failed to execute command: %w", err)
	}

	// This line should never be reached, but included for completeness
	return nil
}
