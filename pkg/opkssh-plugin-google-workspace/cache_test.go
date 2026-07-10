package opksshplugingoogleworkspace

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"testing"
	"time"
)

type fakeGroupMembersFetcher struct{}

func (fakeGroupMembersFetcher) GroupMembers(context.Context, *slog.Logger, string) ([]*Member, error) {
	return []*Member{{Email: "user@example.com"}}, nil
}

func TestCacheWritesPrivateFileAndDirectory(t *testing.T) {
	t.Parallel()

	cachePath := filepath.Join(t.TempDir(), "cache", "cache.json")
	duration := time.Minute
	cache := NewCacheFetcher(ConfigCache{Path: &cachePath, Duration: &duration}, "C1", fakeGroupMembersFetcher{})

	if _, err := cache.GroupMembers(context.Background(), slog.Default(), "group@example.com"); err != nil {
		t.Fatalf("GroupMembers() error = %v", err)
	}

	cacheDirInfo, err := filepath.EvalSymlinks(filepath.Dir(cachePath))
	if err != nil {
		t.Fatalf("EvalSymlinks(cache dir) error = %v", err)
	}
	assertPathMode(t, cacheDirInfo, 0700)
	assertPathMode(t, cachePath, 0600)
}

func assertPathMode(t *testing.T, path string, want uint32) {
	t.Helper()
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat(%s) error = %v", path, err)
	}
	if got := uint32(info.Mode().Perm()); got != want {
		t.Fatalf("%s mode = %04o, want %04o", path, got, want)
	}
}
