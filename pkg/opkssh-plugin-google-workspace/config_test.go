package opksshplugingoogleworkspace

import (
	"os"
	"testing"
)

func TestValidatePrivateReadableFileRejectsWorldAccess(t *testing.T) {
	t.Parallel()

	path := writeModeFixture(t, 0644)
	if err := validatePrivateReadableFile(path, "fixture"); err == nil {
		t.Fatal("validatePrivateReadableFile() error = nil, want rejection")
	}
}

func TestValidatePrivateReadableFileAllowsGroupRead(t *testing.T) {
	t.Parallel()

	path := writeModeFixture(t, 0640)
	if err := validatePrivateReadableFile(path, "fixture"); err != nil {
		t.Fatalf("validatePrivateReadableFile() error = %v", err)
	}
}

func writeModeFixture(t *testing.T, mode uint32) string {
	t.Helper()
	path := t.TempDir() + "/fixture"
	if err := os.WriteFile(path, []byte("fixture"), 0600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	if err := os.Chmod(path, os.FileMode(mode)); err != nil {
		t.Fatalf("Chmod() error = %v", err)
	}
	return path
}
