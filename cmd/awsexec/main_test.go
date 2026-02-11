package main

import "testing"

func Test_buildRunArgs(t *testing.T) {
	tests := []struct {
		name     string
		osArgs   []string
		flagArgs []string
		want     []string
	}{
		{
			name:     "no flags, profile and command with separator",
			osArgs:   []string{"awsexec", "myprofile", "--", "echo", "hello"},
			flagArgs: []string{"myprofile", "--", "echo", "hello"},
			want:     []string{"myprofile", "--", "echo", "hello"},
		},
		{
			name:     "AWS_PROFILE with awsexec -- command",
			osArgs:   []string{"awsexec", "--", "echo", "hello"},
			flagArgs: []string{"echo", "hello"},
			want:     []string{"--", "echo", "hello"},
		},
		{
			name:     "extra flag before profile, no special handling",
			osArgs:   []string{"awsexec", "-x", "myprofile", "--", "echo", "hello"},
			flagArgs: []string{"myprofile", "--", "echo", "hello"},
			want:     []string{"myprofile", "--", "echo", "hello"},
		},
		{
			name:     "no separator, just command args",
			osArgs:   []string{"awsexec", "echo", "hello"},
			flagArgs: []string{"echo", "hello"},
			want:     []string{"echo", "hello"},
		},
		{
			name:     "no flag args",
			osArgs:   []string{"awsexec"},
			flagArgs: []string{},
			want:     []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildRunArgs(tt.osArgs, tt.flagArgs)
			if len(got) != len(tt.want) {
				t.Fatalf("buildRunArgs() len = %d, want %d; got=%v, want=%v", len(got), len(tt.want), got, tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Fatalf("buildRunArgs()[%d] = %q, want %q; got=%v, want=%v", i, got[i], tt.want[i], got, tt.want)
				}
			}
		})
	}
}

