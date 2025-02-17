package listing

import (
	"bufio"
	"bytes"
	"fmt"
	"io/fs"
	"os/exec"
	"path/filepath"
	"strings"
)

// Run is the main entry point for listing logic.
// It returns a string containing the final output (tree or flat + optional content).
func Run(cfg *Config) (string, error) {
	includePatterns := splitPatterns(cfg.Include)
	excludePatterns := splitPatterns(cfg.Exclude)

	// Build a set of Git-tracked files if needed
	var trackedFiles map[string]bool
	var err error
	if cfg.GitTrackedOnly {
		trackedFiles, err = getGitTrackedFiles(cfg.RootPath)
		if err != nil {
			return "", fmt.Errorf("failed to get Git-tracked files: %v", err)
		}
	}

	// We'll store all "accepted" paths
	// We'll also keep a separate slice of "files only" for potential separate content printing
	var entries []string
	var fileEntries []string

	// Walk the root directory
	err = filepath.Walk(cfg.RootPath, func(path string, info fs.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		// Skip the root path in listing output, but still descend
		if path == cfg.RootPath {
			return nil
		}

		// If hidden (e.g. ".git"), skip
		if isHidden(path, cfg.RootPath) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// If user-specified excludes match, skip
		if shouldExclude(path, info, excludePatterns, cfg.RootPath) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// **Key Fix**: Handle directories separately so we always descend.
		if info.IsDir() {
			// We can list the directory if we want it to appear in the final tree,
			// or skip it if we prefer only to show files.
			entries = append(entries, path)
			return nil
		}

		// Now it's a file:
		// If Git-tracked-only, skip files that aren't tracked.
		if cfg.GitTrackedOnly {
			rel := strings.TrimPrefix(filepath.ToSlash(path), filepath.ToSlash(cfg.RootPath))
			rel = strings.TrimPrefix(rel, "/")
			if !trackedFiles[rel] {
				return nil
			}
		}

		// Check if it matches the include patterns
		if len(includePatterns) > 0 && !matchesAnyPattern(info.Name(), includePatterns) {
			return nil
		}

		// Skip binary
		if isBinary(path) {
			return nil
		}

		// If we've made it this far, it's an accepted file
		entries = append(entries, path)
		fileEntries = append(fileEntries, path)

		return nil
	})
	if err != nil {
		return "", err
	}

	// Build up the output
	var outputBuilder strings.Builder

	if cfg.ShowTree {
		// Build tree structure
		tOut := buildTreeOutput(cfg, cfg.RootPath, entries)
		outputBuilder.WriteString(tOut.TreeString)

		// Optionally print file contents separately after the tree
		if cfg.ShowContent && cfg.SeparateContent && len(tOut.FileOrder) > 0 {
			outputBuilder.WriteString("\n")
			outputBuilder.WriteString(buildSeparateContentSection(tOut.FileOrder, cfg))
		}
		// If cfg.ShowContent && !cfg.SeparateContent, the content
		// is already handled inline in buildTreeOutput.
	} else {
		// Flat listing
		fOut := buildFlatListOutput(entries)

		// If no content or separate content, just print the file listing
		if !cfg.ShowContent || cfg.SeparateContent {
			outputBuilder.WriteString(fOut)

			if cfg.ShowContent && cfg.SeparateContent && len(fileEntries) > 0 {
				outputBuilder.WriteString("\n")
				outputBuilder.WriteString(buildSeparateContentSection(fileEntries, cfg))
			}
		} else {
			// We want content inlined with the flat listing
			outputBuilder.WriteString(buildFlatListWithContent(entries, cfg))
		}
	}

	return outputBuilder.String(), nil
}

// splitPatterns takes a comma-separated string of patterns and splits them
func splitPatterns(patterns string) []string {
	if patterns == "" {
		return nil
	}
	parts := strings.Split(patterns, ",")
	for i := 0; i < len(parts); i++ {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

// getGitTrackedFiles uses "git ls-files" to list tracked files
func getGitTrackedFiles(root string) (map[string]bool, error) {
	cmd := exec.Command("git", "ls-files")
	cmd.Dir = root

	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	tracked := make(map[string]bool)
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		tracked[scanner.Text()] = true
	}

	return tracked, nil
}
