package listing

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestSplitPatterns ensures splitPatterns works as expected
func TestSplitPatterns(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{"", nil},
		{"*.go", []string{"*.go"}},
		{"*.go, *.md", []string{"*.go", "*.md"}},
	}
	for _, c := range cases {
		got := splitPatterns(c.input)
		if len(got) != len(c.expected) {
			t.Errorf("splitPatterns(%q) length = %d; want %d", c.input, len(got), len(c.expected))
			continue
		}
		for i := range got {
			if got[i] != c.expected[i] {
				t.Errorf("splitPatterns(%q)[%d] = %q; want %q", c.input, i, got[i], c.expected[i])
			}
		}
	}
}

// TestRunBasic uses a small in-memory structure to test listing logic
func TestRunBasic(t *testing.T) {
	tmp := t.TempDir()
	// make subdir
	sub := filepath.Join(tmp, "sub")
	if err := os.Mkdir(sub, 0755); err != nil {
		t.Fatal(err)
	}

	// create a text file
	txt1 := filepath.Join(tmp, "hello.txt")
	_ = ioutil.WriteFile(txt1, []byte("hello world"), 0644)

	// create a binary file
	bin1 := filepath.Join(tmp, "binary.bin")
	_ = ioutil.WriteFile(bin1, []byte{0x00, 0x01, 0x02}, 0644)

	// create hidden file
	hidden := filepath.Join(tmp, ".gitignore")
	_ = ioutil.WriteFile(hidden, []byte("ignore stuff"), 0644)

	cfg := &Config{
		RootPath:          tmp,
		Include:           "", // no pattern
		Exclude:           "", // no exclude
		GitTrackedOnly:    false,
		ShowTree:          true,
		ShowContent:       true,
		SeparateContent:   false,
		Output:            "",
		ShowLineNumbers:   false,
		ShowHeaderFooters: true,
	}

	out, err := Run(cfg)
	if err != nil {
		t.Fatalf("Run error: %v", err)
	}

	// Expect "hello.txt" content, not binary.bin, not .gitignore
	if !strings.Contains(out, "hello.txt") {
		t.Error("Expected hello.txt to appear in output.")
	}
	if strings.Contains(out, "binary.bin") {
		t.Error("Expected binary.bin to be skipped.")
	}
	if strings.Contains(out, ".gitignore") {
		t.Error("Expected hidden file to be skipped.")
	}
}

// If you want to test getGitTrackedFiles, you can do so by spinning up a git repo,
// but that can be more complex. For now, you can rely on the coverage from above.
