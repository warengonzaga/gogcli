package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPaths_CreateDirs(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	dir, err := EnsureDir()
	if err != nil {
		t.Fatalf("EnsureDir: %v", err)
	}

	if _, statErr := os.Stat(dir); statErr != nil {
		t.Fatalf("expected dir: %v", statErr)
	}

	if filepath.Base(dir) != AppName {
		t.Fatalf("unexpected base: %q", filepath.Base(dir))
	}

	keyringDir, err := EnsureKeyringDir()
	if err != nil {
		t.Fatalf("EnsureKeyringDir: %v", err)
	}

	if _, statErr := os.Stat(keyringDir); statErr != nil {
		t.Fatalf("expected keyring dir: %v", statErr)
	}

	downloadsDir, err := EnsureDriveDownloadsDir()
	if err != nil {
		t.Fatalf("EnsureDriveDownloadsDir: %v", err)
	}

	if _, statErr := os.Stat(downloadsDir); statErr != nil {
		t.Fatalf("expected downloads dir: %v", statErr)
	}

	attachmentsDir, err := EnsureGmailAttachmentsDir()
	if err != nil {
		t.Fatalf("EnsureGmailAttachmentsDir: %v", err)
	}

	if _, statErr := os.Stat(attachmentsDir); statErr != nil {
		t.Fatalf("expected attachments dir: %v", statErr)
	}

	watchDir, err := EnsureGmailWatchDir()
	if err != nil {
		t.Fatalf("EnsureGmailWatchDir: %v", err)
	}

	if _, statErr := os.Stat(watchDir); statErr != nil {
		t.Fatalf("expected watch dir: %v", statErr)
	}

	credsPath, err := ClientCredentialsPath()
	if err != nil {
		t.Fatalf("ClientCredentialsPath: %v", err)
	}

	if filepath.Base(credsPath) != "credentials.json" {
		t.Fatalf("unexpected creds file: %q", filepath.Base(credsPath))
	}
}
