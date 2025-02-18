package listing

import (
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestIsBinaryData(t *testing.T) {
	text := []byte("Hello world")
	bin := []byte{0x00, 0x01, 0x02}

	if isBinaryData(text) {
		t.Error("Expected plain text to not be recognized as binary.")
	}
	if !isBinaryData(bin) {
		t.Error("Expected null-containing bytes to be recognized as binary.")
	}
}

func TestIsBinary(t *testing.T) {
	// Create temp text file
	txtFile, err := ioutil.TempFile("", "testfile-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(txtFile.Name())
	_, _ = txtFile.Write([]byte("hello"))
	txtFile.Close()

	if isBinary(txtFile.Name()) {
		t.Error("Expected text file to not be recognized as binary.")
	}

	// Create temp binary file
	binFile, err := ioutil.TempFile("", "testfile-*.bin")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(binFile.Name())
	_, _ = binFile.Write([]byte{0x00, 0x01, 0x02})
	binFile.Close()

	if !isBinary(binFile.Name()) {
		t.Error("Expected binary file to be recognized as binary.")
	}
}

func TestIsHidden(t *testing.T) {
	root := "/fake/root"
	// hidden file
	if !isHidden("/fake/root/.secret", root) {
		t.Error("Expected .secret to be hidden")
	}
	// nested
	if !isHidden("/fake/root/folder/.git/config", root) {
		t.Error("Expected .git/config to be hidden")
	}
	// not hidden
	if isHidden("/fake/root/folder/main.go", root) {
		t.Error("Expected main.go NOT to be hidden")
	}
}

func TestShouldExclude(t *testing.T) {
	root := "/fake/root"
	exPatterns := []string{".git", "node_modules", "*.env"}

	finfoFile := mockFileInfo{name: "main.env", isDirVal: false}
	finfoDir := mockFileInfo{name: "node_modules", isDirVal: true}

	// matches *.env
	if !shouldExclude("/fake/root/main.env", finfoFile, exPatterns, root) {
		t.Error("Expected main.env to be excluded by *.env.")
	}
	// node_modules
	if !shouldExclude("/fake/root/node_modules", finfoDir, exPatterns, root) {
		t.Error("Expected node_modules to be excluded by direct match.")
	}
	// not excluded
	finfoNormal := mockFileInfo{name: "main.go", isDirVal: false}
	if shouldExclude("/fake/root/main.go", finfoNormal, exPatterns, root) {
		t.Error("Expected main.go NOT to be excluded.")
	}
}

func TestMatchesAnyPattern(t *testing.T) {
	patterns := []string{"*.go", "*.md"}
	if !matchesAnyPattern("main.go", patterns) {
		t.Error("Expected main.go to match *.go")
	}
	if matchesAnyPattern("main.py", patterns) {
		t.Error("Expected main.py NOT to match *.go or *.md")
	}
}

// mockFileInfo is a tiny test helper
type mockFileInfo struct {
	name     string
	isDirVal bool
}

func (m mockFileInfo) Name() string      { return m.name }
func (m mockFileInfo) Size() int64       { return 0 }
func (m mockFileInfo) Mode() os.FileMode { return 0644 }

// Return a real time.Time instead of a custom type:
func (m mockFileInfo) ModTime() time.Time { return time.Now() }
func (m mockFileInfo) IsDir() bool        { return m.isDirVal }
func (m mockFileInfo) Sys() interface{}   { return nil }
