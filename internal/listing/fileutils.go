package listing

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// isBinary does a naive check if a file is binary by scanning the first 8KB for null bytes.
func isBinary(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return true // if we can't open, assume binary
	}
	defer f.Close()

	buf := make([]byte, 8000)
	n, err := f.Read(buf)
	if err != nil && err != io.EOF {
		return true
	}
	return isBinaryData(buf[:n])
}

func isBinaryData(buf []byte) bool {
	if bytes.Contains(buf, []byte{0}) {
		return true
	}
	return false
}

// isHidden checks if the file/dir starts with a dot relative to the root
func isHidden(path, root string) bool {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return false // fallback
	}
	parts := strings.Split(rel, string(filepath.Separator))
	for _, p := range parts {
		if strings.HasPrefix(p, ".") {
			return true
		}
	}
	return false
}

// shouldExclude checks if the path or directory name matches the exclude list
func shouldExclude(path string, info os.FileInfo, excludePatterns []string, root string) bool {
	base := info.Name()
	for _, pattern := range excludePatterns {
		matched, _ := filepath.Match(pattern, base)
		if matched || base == pattern {
			return true
		}
	}
	return false
}

// matchesAnyPattern checks if filename matches any of the include patterns
func matchesAnyPattern(name string, patterns []string) bool {
	for _, p := range patterns {
		matched, err := filepath.Match(p, name)
		if err == nil && matched {
			return true
		}
	}
	return false
}
