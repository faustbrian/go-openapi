package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadFixtureConfinesReadsToAuthorizedRoot(t *testing.T) {
	t.Parallel()

	base := t.TempDir()
	allowed := filepath.Join(base, "allowed")
	if err := os.Mkdir(allowed, 0o700); err != nil {
		t.Fatalf("create authorized root: %v", err)
	}
	if err := os.WriteFile(filepath.Join(allowed, "valid.json"), []byte("{}"), 0o600); err != nil {
		t.Fatalf("write valid fixture: %v", err)
	}
	outside := filepath.Join(base, "outside.json")
	if err := os.WriteFile(outside, []byte("{}"), 0o600); err != nil {
		t.Fatalf("write outside fixture: %v", err)
	}
	if err := os.Symlink(outside, filepath.Join(allowed, "escape.json")); err != nil {
		t.Fatalf("create escaping symlink: %v", err)
	}

	root, err := os.OpenRoot(allowed)
	if err != nil {
		t.Fatalf("open authorized root: %v", err)
	}
	t.Cleanup(func() {
		if err := root.Close(); err != nil {
			t.Errorf("close authorized root: %v", err)
		}
	})

	contents, err := readFixture(root, "valid.json", 2)
	if err != nil {
		t.Fatalf("read valid fixture: %v", err)
	}
	if string(contents) != "{}" {
		t.Fatalf("fixture contents = %q, want %q", contents, "{}")
	}

	for name, path := range map[string]string{
		"parent traversal": "../outside.json",
		"absolute path":    outside,
		"symlink escape":   "escape.json",
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if _, err := readFixture(root, path, 2); err == nil {
				t.Fatal("read succeeded outside the authorized root")
			}
		})
	}
}

func TestReadFixtureRejectsOversizedInput(t *testing.T) {
	t.Parallel()

	directory := t.TempDir()
	if err := os.WriteFile(filepath.Join(directory, "large.json"), []byte("1234"), 0o600); err != nil {
		t.Fatalf("write oversized fixture: %v", err)
	}
	root, err := os.OpenRoot(directory)
	if err != nil {
		t.Fatalf("open authorized root: %v", err)
	}
	t.Cleanup(func() {
		if err := root.Close(); err != nil {
			t.Errorf("close authorized root: %v", err)
		}
	})

	if _, err := readFixture(root, "large.json", 3); err == nil {
		t.Fatal("read succeeded above the fixture byte limit")
	}
}
