package listing

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildTreeOutput(t *testing.T) {
	tmp := t.TempDir()
	sub := filepath.Join(tmp, "dir")
	os.MkdirAll(sub, 0755)

	file1 := filepath.Join(tmp, "file1.txt")
	os.WriteFile(file1, []byte("file1 content"), 0644)
	file2 := filepath.Join(sub, "file2.md")
	os.WriteFile(file2, []byte("file2 content"), 0644)

	// We'll emulate a small slice of accepted paths
	entries := []string{file1, sub, file2}

	cfg := &Config{
		RootPath:          tmp,
		ShowTree:          true,
		ShowContent:       true,
		SeparateContent:   false,
		ShowLineNumbers:   false,
		ShowHeaderFooters: true,
	}

	treeOut := buildTreeOutput(cfg, tmp, entries)
	if !strings.Contains(treeOut.TreeString, "file1.txt") {
		t.Error("Expected file1.txt in tree output")
	}
	if !strings.Contains(treeOut.TreeString, "file2.md") {
		t.Error("Expected file2.md in tree output")
	}
	if !strings.Contains(treeOut.TreeString, "----- CONTENT START -----") {
		t.Error("Expected inline content markers in tree output")
	}
	if len(treeOut.FileOrder) != 2 {
		t.Errorf("Expected 2 files in FileOrder, got %d", len(treeOut.FileOrder))
	}
}

func TestBuildFlatListOutput(t *testing.T) {
	entries := []string{"file1.txt", "dir", "file2.md"}
	out := buildFlatListOutput(entries)
	if !strings.Contains(out, "file1.txt") {
		t.Error("Expected file1.txt in flat output")
	}
	if !strings.Contains(out, "dir") {
		t.Error("Expected dir in flat output")
	}
	if !strings.Contains(out, "file2.md") {
		t.Error("Expected file2.md in flat output")
	}
}

func TestBuildFlatListWithContent(t *testing.T) {
	tmp := t.TempDir()
	file1 := filepath.Join(tmp, "file1.txt")
	os.WriteFile(file1, []byte("hello"), 0644)
	dir1 := filepath.Join(tmp, "sub")
	os.MkdirAll(dir1, 0755)

	entries := []string{file1, dir1}
	cfg := &Config{
		ShowLineNumbers:   false,
		ShowHeaderFooters: true,
	}

	out := buildFlatListWithContent(entries, cfg)
	if !strings.Contains(out, "file1.txt") {
		t.Error("Expected file1.txt in output")
	}
	if !strings.Contains(out, "----- CONTENT START -----") {
		t.Error("Expected content markers")
	}
	// directory "sub" should appear, but no content
	if !strings.Contains(out, "sub") {
		t.Error("Expected 'sub' directory name in output")
	}
	if strings.Contains(out, "hello") == false {
		t.Error("Expected file1 content 'hello' in output")
	}
}

func TestBuildSeparateContentSection(t *testing.T) {
	tmp := t.TempDir()
	file1 := filepath.Join(tmp, "file1.txt")
	os.WriteFile(file1, []byte("line1\nline2"), 0644)
	file2 := filepath.Join(tmp, "file2.md")
	os.WriteFile(file2, []byte("## doc"), 0644)

	cfg := &Config{
		ShowLineNumbers:   true,
		ShowHeaderFooters: true,
	}

	out := buildSeparateContentSection([]string{file1, file2}, cfg)
	if !strings.Contains(out, "line1") || !strings.Contains(out, "line2") {
		t.Error("Expected file1 lines in separate content")
	}
	if !strings.Contains(out, "## doc") {
		t.Error("Expected file2 content in separate content")
	}
	if !strings.Contains(out, "----- CONTENT START -----") {
		t.Error("Expected header-footer markers")
	}
	if !strings.Contains(out, "   1: line1") {
		t.Error("Expected line numbering in output ( '   1: line1' )")
	}
}
