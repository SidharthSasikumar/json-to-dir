package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOutputConfinement(t *testing.T) {
	for _, key := range []string{"../escaped", "/absolute", "a/../../escaped", `..\escaped`} {
		t.Run(key, func(t *testing.T) {
			if err := createDirectoriesAndFiles(map[string]interface{}{key: "bad"}, t.TempDir()); err == nil {
				t.Fatal("unsafe key accepted")
			}
		})
	}
	outside := t.TempDir()
	root := t.TempDir()
	if err := os.Symlink(outside, filepath.Join(root, "link")); err != nil {
		t.Fatal(err)
	}
	if err := createDirectoriesAndFiles(map[string]interface{}{"link/escaped": "bad"}, root); err == nil {
		t.Fatal("escaped symlink")
	}
	if _, err := os.Stat(filepath.Join(outside, "escaped.json")); !os.IsNotExist(err) {
		t.Fatal("outside path touched")
	}
}
func TestCreatesNestedOutput(t *testing.T) {
	root := filepath.Join(t.TempDir(), "new")
	if err := createDirectoriesAndFiles(map[string]interface{}{"nested/file": "ok"}, root); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(filepath.Join(root, "nested/file.json"))
	if err != nil || string(b) != `"ok"` {
		t.Fatalf("%q %v", b, err)
	}
}
