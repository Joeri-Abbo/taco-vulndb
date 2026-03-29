package cli

import (
	"testing"
)

func TestNewRootCmd(t *testing.T) {
	cmd := NewRootCmd()

	if cmd.Use != "taco-vulndb" {
		t.Errorf("expected Use 'taco-vulndb', got %s", cmd.Use)
	}

	// Check all subcommands are registered
	expectedCmds := map[string]bool{
		"update": false, "download": false, "load": false,
		"build": false, "export": false, "serve": false,
		"status": false, "push": false, "pull": false,
	}

	for _, sub := range cmd.Commands() {
		if _, ok := expectedCmds[sub.Name()]; ok {
			expectedCmds[sub.Name()] = true
		}
	}

	for name, found := range expectedCmds {
		if !found {
			t.Errorf("expected subcommand %q not found", name)
		}
	}
}

func TestRootCmd_GlobalFlags(t *testing.T) {
	cmd := NewRootCmd()

	debugFlag := cmd.PersistentFlags().Lookup("debug")
	if debugFlag == nil {
		t.Error("expected --debug flag")
	}

	quietFlag := cmd.PersistentFlags().Lookup("quiet")
	if quietFlag == nil {
		t.Error("expected --quiet flag")
	}
}

func TestRootCmd_SilenceUsageAndErrors(t *testing.T) {
	cmd := NewRootCmd()
	if !cmd.SilenceUsage {
		t.Error("expected SilenceUsage to be true")
	}
	if !cmd.SilenceErrors {
		t.Error("expected SilenceErrors to be true")
	}
}
