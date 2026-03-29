package cli

import (
	"testing"
)

func TestNewDownloadCmd(t *testing.T) {
	cmd := newDownloadCmd()
	if cmd.Use != "download" {
		t.Errorf("expected Use 'download', got %s", cmd.Use)
	}

	urlFlag := cmd.Flags().Lookup("url")
	if urlFlag == nil {
		t.Fatal("expected --url flag")
	}

	// Should fail without --url
	cmd.SetArgs([]string{})
	err := cmd.Execute()
	if err == nil {
		t.Error("expected error without --url flag")
	}
}

func TestNewLoadCmd(t *testing.T) {
	cmd := newLoadCmd()
	if cmd.Use != "load" {
		t.Errorf("expected Use 'load', got %s", cmd.Use)
	}

	fileFlag := cmd.Flags().Lookup("file")
	if fileFlag == nil {
		t.Fatal("expected --file flag")
	}

	// Should fail without --file
	cmd.SetArgs([]string{})
	err := cmd.Execute()
	if err == nil {
		t.Error("expected error without --file flag")
	}
}

func TestNewBuildCmd(t *testing.T) {
	cmd := newBuildCmd()
	if cmd.Use != "build" {
		t.Errorf("expected Use 'build', got %s", cmd.Use)
	}

	outputFlag := cmd.Flags().Lookup("output")
	if outputFlag == nil {
		t.Fatal("expected --output flag")
	}

	daysFlag := cmd.Flags().Lookup("days")
	if daysFlag == nil {
		t.Fatal("expected --days flag")
	}
	if daysFlag.DefValue != "120" {
		t.Errorf("expected days default 120, got %s", daysFlag.DefValue)
	}

	// Should fail without --output
	cmd.SetArgs([]string{})
	err := cmd.Execute()
	if err == nil {
		t.Error("expected error without --output flag")
	}
}

func TestNewExportCmd(t *testing.T) {
	cmd := newExportCmd()
	if cmd.Use != "export" {
		t.Errorf("expected Use 'export', got %s", cmd.Use)
	}

	outputFlag := cmd.Flags().Lookup("output")
	if outputFlag == nil {
		t.Fatal("expected --output flag")
	}

	// Should fail without --output
	cmd.SetArgs([]string{})
	err := cmd.Execute()
	if err == nil {
		t.Error("expected error without --output flag")
	}
}

func TestNewServeCmd(t *testing.T) {
	cmd := newServeCmd()
	if cmd.Use != "serve" {
		t.Errorf("expected Use 'serve', got %s", cmd.Use)
	}

	addrFlag := cmd.Flags().Lookup("addr")
	if addrFlag == nil {
		t.Fatal("expected --addr flag")
	}
	if addrFlag.DefValue != ":8080" {
		t.Errorf("expected default addr :8080, got %s", addrFlag.DefValue)
	}
}

func TestNewUpdateCmd(t *testing.T) {
	cmd := newUpdateCmd()
	if cmd.Use != "update" {
		t.Errorf("expected Use 'update', got %s", cmd.Use)
	}

	sourcesFlag := cmd.Flags().Lookup("sources")
	if sourcesFlag == nil {
		t.Fatal("expected --sources flag")
	}

	fullFlag := cmd.Flags().Lookup("full")
	if fullFlag == nil {
		t.Fatal("expected --full flag")
	}
}

func TestNewStatusCmd(t *testing.T) {
	cmd := newStatusCmd()
	if cmd.Use != "status" {
		t.Errorf("expected Use 'status', got %s", cmd.Use)
	}
}

func TestNewPushCmd(t *testing.T) {
	cmd := newPushCmd()
	if cmd.Use != "push <oci-reference>" {
		t.Errorf("unexpected Use: %s", cmd.Use)
	}

	// Should require exactly 1 arg
	cmd.SetArgs([]string{})
	err := cmd.Execute()
	if err == nil {
		t.Error("expected error without args")
	}
}

func TestNewPullCmd(t *testing.T) {
	cmd := newPullCmd()
	if cmd.Use != "pull <oci-reference>" {
		t.Errorf("unexpected Use: %s", cmd.Use)
	}

	// Should require exactly 1 arg
	cmd.SetArgs([]string{})
	err := cmd.Execute()
	if err == nil {
		t.Error("expected error without args")
	}
}
