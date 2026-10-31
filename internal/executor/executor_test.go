package executor

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"kook/internal/config"

	"github.com/spf13/cobra"
)

func TestIncludedCommandEnv(t *testing.T) {
	callerDir := t.TempDir()
	includedDir := t.TempDir()
	out := filepath.Join(t.TempDir(), "out")

	prevDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(callerDir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chdir(prevDir) })

	cfg := &config.Config{Dir: includedDir, IsIncluded: true}
	cmd := config.Command{
		Name:   "env",
		Silent: true,
		Script: `printf '%s\n%s\n%s' "$PWD" "$KOOK_CALLER_DIR" "$KOOK" > ` + out,
	}

	if err := Execute(cfg, cmd, &cobra.Command{}); err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}
	lines := strings.Split(string(data), "\n")

	if lines[0] != includedDir {
		t.Errorf("Expected working dir %s, got %s", includedDir, lines[0])
	}
	if lines[1] != callerDir {
		t.Errorf("Expected KOOK_CALLER_DIR %s, got %s", callerDir, lines[1])
	}
	if lines[2] == "" {
		t.Error("Expected KOOK to be set")
	}
}
