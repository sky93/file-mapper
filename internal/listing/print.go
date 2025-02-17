package listing

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// TreeOutput helps hold both the tree string and
// the ordered list of files encountered (for separate content printing).
type TreeOutput struct {
	TreeString string   // the actual ASCII tree
	FileOrder  []string // the order in which files appeared in the tree
}

// buildTreeOutput creates a tree-like output from the list of entries
// and returns a TreeOutput struct. If cfg.ShowContent && !cfg.SeparateContent,
// it will inline the content under each file in the tree itself.
func buildTreeOutput(cfg *Config, root string, entries []string) *TreeOutput {
	// Build map of dir -> children
	treeMap := make(map[string][]string)
	for _, e := range entries {
		rel, _ := filepath.Rel(root, e)
		dir := filepath.Dir(rel)
		treeMap[dir] = append(treeMap[dir], rel)
	}

	// Sort for consistent output
	for k := range treeMap {
		sort.Strings(treeMap[k])
	}

	builder := &strings.Builder{}
	var fileOrder []string

	// We'll recurse from top-level (".")
	recurseTree(builder, cfg, root, ".", treeMap, 0, &fileOrder)

	return &TreeOutput{
		TreeString: builder.String(),
		FileOrder:  fileOrder,
	}
}

// recurseTree is a recursive helper to print directories/files in a tree view.
func recurseTree(
	sb *strings.Builder,
	cfg *Config,
	root string,
	dir string,
	treeMap map[string][]string,
	level int,
	fileOrder *[]string,
) {
	children, ok := treeMap[dir]
	if !ok {
		return
	}

	for i, child := range children {
		indent(sb, level)

		connector := "├──"
		if i == len(children)-1 {
			connector = "└──"
		}

		base := filepath.Base(child)
		sb.WriteString(fmt.Sprintf("%s %s\n", connector, base))

		// Is child a directory with further children?
		if hasChildren(treeMap, child) {
			// Recurse deeper
			recurseTree(sb, cfg, root, child, treeMap, level+1, fileOrder)
		} else {
			// It's a file
			fullPath := filepath.Join(root, child)
			*fileOrder = append(*fileOrder, fullPath)

			// If we should show content inline (tree + content, but NOT separate)
			if cfg.ShowContent && !cfg.SeparateContent {
				printInlineContent(sb, cfg, fullPath, level+1)
			}
		}
	}
}

// printInlineContent prints the content of a single file inline,
// under the current tree level. We handle line-numbers and header-footers here.
func printInlineContent(sb *strings.Builder, cfg *Config, filePath string, level int) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return
	}

	lines := strings.Split(string(content), "\n")

	if cfg.ShowHeaderFooters {
		// Indent a line, print "----- CONTENT START -----"
		indent(sb, level)
		sb.WriteString("----- CONTENT START -----\n")
	}

	if cfg.ShowLineNumbers {
		// Print each line with line numbers and indentation
		for i, line := range lines {
			indent(sb, level)
			sb.WriteString(fmt.Sprintf("%4d: %s\n", i+1, line))
		}
	} else {
		// Print each line with indentation, but no line numbering
		for _, line := range lines {
			indent(sb, level)
			sb.WriteString(line + "\n")
		}
	}

	if cfg.ShowHeaderFooters {
		indent(sb, level)
		sb.WriteString("----- CONTENT END -----\n")
	}
}

// hasChildren checks if there are sub-entries for the given key
func hasChildren(treeMap map[string][]string, key string) bool {
	_, ok := treeMap[key]
	return ok
}

// buildFlatListOutput returns a simple list of all entries (dirs + files)
func buildFlatListOutput(entries []string) string {
	var sb strings.Builder
	for _, e := range entries {
		sb.WriteString(e + "\n")
	}
	return sb.String()
}

// buildFlatListWithContent inlines file content after each file path
func buildFlatListWithContent(entries []string, cfg *Config) string {
	var sb strings.Builder
	for _, e := range entries {
		info, err := os.Stat(e)
		if err != nil || info.IsDir() {
			// Just print directories or skip on error
			sb.WriteString(e + "\n")
			continue
		}
		// It's a file
		sb.WriteString(e + "\n")

		content, err := os.ReadFile(e)
		if err != nil {
			continue
		}
		lines := strings.Split(string(content), "\n")

		// Optional header/footer
		if cfg.ShowHeaderFooters {
			sb.WriteString("----- CONTENT START -----\n")
		}
		if cfg.ShowLineNumbers {
			for i, line := range lines {
				sb.WriteString(fmt.Sprintf("%4d: %s\n", i+1, line))
			}
		} else {
			sb.WriteString(string(content))
			// Ensure trailing newline
			if !strings.HasSuffix(string(content), "\n") {
				sb.WriteString("\n")
			}
		}
		if cfg.ShowHeaderFooters {
			sb.WriteString("----- CONTENT END -----\n")
		}
	}
	return sb.String()
}

// buildSeparateContentSection prints content for each file (by path) in order
// e.g. "internal/listing/listing.go (60 lines):"
func buildSeparateContentSection(filePaths []string, cfg *Config) string {
	var sb strings.Builder
	for _, path := range filePaths {
		info, err := os.Stat(path)
		if err != nil || info.IsDir() {
			continue
		}
		content, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		lines := strings.Split(string(content), "\n")
		lineCount := len(lines)

		// "filename (NN lines):"
		sb.WriteString(fmt.Sprintf("%s (%d lines):\n", path, lineCount))

		if cfg.ShowHeaderFooters {
			sb.WriteString("----- CONTENT START -----\n")
		}
		if cfg.ShowLineNumbers {
			for i, line := range lines {
				sb.WriteString(fmt.Sprintf("%4d: %s\n", i+1, line))
			}
		} else {
			sb.WriteString(string(content))
			// Ensure trailing newline
			if !strings.HasSuffix(string(content), "\n") {
				sb.WriteString("\n")
			}
		}
		if cfg.ShowHeaderFooters {
			sb.WriteString("----- CONTENT END -----\n")
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

// indent writes indentation for the tree
func indent(sb *strings.Builder, level int) {
	for i := 0; i < level; i++ {
		sb.WriteString("│   ")
	}
}
