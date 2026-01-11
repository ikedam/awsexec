package awsexec

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
)

// AWSCommandExporter implements CredentialExporter using aws configure export-credentials command.
type AWSCommandExporter struct{}

// NewAWSCommandExporter creates a new AWSCommandExporter instance.
func NewAWSCommandExporter() *AWSCommandExporter {
	return &AWSCommandExporter{}
}

// ExportCredentials executes aws configure export-credentials and returns the credentials.
func (e *AWSCommandExporter) ExportCredentials(ctx context.Context, profile string) (*Credentials, error) {
	// Build the command
	// `process` means JSON format.
	cmdArgs := []string{"configure", "export-credentials", "--format", "process"}
	if profile != "" {
		cmdArgs = append(cmdArgs, "--profile", profile)
	}

	cmd := exec.CommandContext(ctx, "aws", cmdArgs...)
	output, err := cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("aws configure export-credentials failed: %s: %w", string(exitError.Stderr), err)
		}
		return nil, fmt.Errorf("failed to execute aws configure export-credentials: %w", err)
	}

	// Parse JSON output
	var creds Credentials
	if err := json.Unmarshal(output, &creds); err != nil {
		return nil, fmt.Errorf("failed to parse credentials JSON: %w", err)
	}

	return &creds, nil
}
