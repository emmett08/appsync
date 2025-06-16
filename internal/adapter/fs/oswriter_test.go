package fs_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/emmett08/1ai-pr/internal/adapter/fs"
)

func TestWrite(t *testing.T) {
	tmp := t.TempDir()
	target := filepath.Join(tmp, "a/b/c.txt")

	w := fs.New()
	if err := w.Write(target, []byte("data")); err != nil {
		t.Fatalf("write failed: %v", err)
	}
	b, _ := os.ReadFile(target)
	if string(b) != "data" {
		t.Fatalf("unexpected contents: %s", b)
	}
}
