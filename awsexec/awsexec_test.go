package awsexec

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"testing"
)

func TestAwsexec_parseArgs(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		envProfile  string
		wantProfile string
		wantCommand []string
		wantErr     bool
	}{
		{
			name:        "profile specified as argument",
			args:        []string{"myprofile", "--", "echo", "hello"},
			envProfile:  "",
			wantProfile: "myprofile",
			wantCommand: []string{"echo", "hello"},
			wantErr:     false,
		},
		{
			name:        "profile from environment variable",
			args:        []string{"--", "echo", "hello"},
			envProfile:  "myprofile",
			wantProfile: "myprofile",
			wantCommand: []string{"echo", "hello"},
			wantErr:     false,
		},
		{
			name:        "missing separator",
			args:        []string{"myprofile", "echo", "hello"},
			envProfile:  "",
			wantProfile: "",
			wantCommand: nil,
			wantErr:     true,
		},
		{
			name:        "no command after separator",
			args:        []string{"myprofile", "--"},
			envProfile:  "",
			wantProfile: "",
			wantCommand: nil,
			wantErr:     true,
		},
		{
			name:        "no profile and no AWS_PROFILE",
			args:        []string{"--", "echo", "hello"},
			envProfile:  "",
			wantProfile: "",
			wantCommand: nil,
			wantErr:     true,
		},
		{
			name:        "command with multiple arguments",
			args:        []string{"myprofile", "--", "sh", "-c", "echo hello"},
			envProfile:  "",
			wantProfile: "myprofile",
			wantCommand: []string{"sh", "-c", "echo hello"},
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up environment
			if tt.envProfile != "" {
				os.Setenv("AWS_PROFILE", tt.envProfile)
				defer os.Unsetenv("AWS_PROFILE")
			} else {
				os.Unsetenv("AWS_PROFILE")
			}

			a := New(context.Background())
			gotProfile, gotCommand, err := a.parseArgs(tt.args)

			if (err != nil) != tt.wantErr {
				t.Errorf("parseArgs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if gotProfile != tt.wantProfile {
					t.Errorf("parseArgs() gotProfile = %v, want %v", gotProfile, tt.wantProfile)
				}
				if !reflect.DeepEqual(gotCommand, tt.wantCommand) {
					t.Errorf("parseArgs() gotCommand = %v, want %v", gotCommand, tt.wantCommand)
				}
			}
		})
	}
}

// mockCredentialExporter is a mock implementation of CredentialExporter for testing.
type mockCredentialExporter struct {
	credentials *Credentials
	err         error
}

func (m *mockCredentialExporter) ExportCredentials(_ context.Context, _ string) (*Credentials, error) {
	return m.credentials, m.err
}

func TestAwsexec_Run(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		envProfile  string
		exporter    CredentialExporter
		wantErr     bool
		errContains string
	}{
		{
			name:       "export credentials error",
			args:       []string{"myprofile", "--", "echo", "hello"},
			envProfile: "",
			exporter: &mockCredentialExporter{
				credentials: nil,
				err:         errors.New("export failed"),
			},
			wantErr:     true,
			errContains: "failed to export credentials",
		},
		{
			name:        "missing separator",
			args:        []string{"myprofile", "echo", "hello"},
			envProfile:  "",
			exporter:    &mockCredentialExporter{},
			wantErr:     true,
			errContains: "failed to parse arguments",
		},
		{
			name:        "no command after separator",
			args:        []string{"myprofile", "--"},
			envProfile:  "",
			exporter:    &mockCredentialExporter{},
			wantErr:     true,
			errContains: "failed to parse arguments",
		},
		{
			name:       "command not found",
			args:       []string{"myprofile", "--", "/nonexistent/command", "arg"},
			envProfile: "",
			exporter: &mockCredentialExporter{
				credentials: &Credentials{
					AccessKeyID:     "AKIAIOSFODNN7EXAMPLE",
					SecretAccessKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
					SessionToken:    "test-session-token",
				},
				err: nil,
			},
			wantErr:     true,
			errContains: "command not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up environment
			if tt.envProfile != "" {
				os.Setenv("AWS_PROFILE", tt.envProfile)
				defer os.Unsetenv("AWS_PROFILE")
			} else {
				os.Unsetenv("AWS_PROFILE")
			}

			a := NewWithExporter(tt.exporter)
			err := a.Run(context.Background(), tt.args)

			if (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errContains != "" {
				if err == nil || !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Run() error = %v, want error containing %q", err, tt.errContains)
				}
			}
		})
	}
}

func TestAWSCommandExporter_ExportCredentials(t *testing.T) {
	// This test requires AWS CLI to be installed and configured.
	// Skip if not available.
	if _, err := exec.Command("aws", "--version").Output(); err != nil {
		t.Skip("AWS CLI not available, skipping integration test")
	}

	exporter := NewAWSCommandExporter()
	ctx := context.Background()

	// Test with empty profile (uses default)
	creds, err := exporter.ExportCredentials(ctx, "")
	if err != nil {
		t.Logf("ExportCredentials with empty profile failed (this is expected if AWS is not configured): %v", err)
		// This is expected if AWS is not configured, so we don't fail the test
		return
	}

	if creds == nil {
		t.Error("ExportCredentials() returned nil credentials")
		return
	}

	if creds.AccessKeyID == "" {
		t.Error("ExportCredentials() returned empty AccessKeyID")
	}
	if creds.SecretAccessKey == "" {
		t.Error("ExportCredentials() returned empty SecretAccessKey")
	}
}
