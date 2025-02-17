# file-mapper

**file-mapper** is a versatile CLI tool that combines the power of the Linux **`tree`** command (for hierarchical directory listings) and the **`cat`** command (for file content viewing)—all rolled into one convenient package. Recursively map out your project structure, filter files, and optionally display file contents (with features like line numbering). It's ideal for developers who want a comprehensive project overview or need to share both file structure and content with others (e.g., in code reviews or AI tools).

---

## Features

1. **Tree or Flat**
    - By default, `file-mapper` displays a **hierarchical tree** of your files and directories (like `tree`).
    - Use the `--flat` flag for a **flat** listing instead.

2. **Git-Tracked-Only**
    - Restrict the output to only files that are **tracked by Git** (`--git`).

3. **Include / Exclude Patterns**
    - Filter specific file types (e.g. `--include="*.go,*.md"`)
    - Exclude directories/files (e.g. `--exclude=".git,node_modules"`).

4. **Hidden & Binary Skips**
    - Hidden files/directories (those starting with `.`) are **ignored by default**.
    - Uses a naive approach to skip **binary** files.

5. **Content Viewing**
    - **Inline** content right below each file (similar to `cat`, but recursive) (`--content`).
    - Or **separate**: list all files, then dump their contents afterward.
    - Enable **line numbers** (`--line-numbers`) for quick reference.
    - Show or hide content headers (`----- CONTENT START -----` / `----- CONTENT END -----`).

6. **Output to File**
    - Save your entire listing and contents to a file (`--output=out.txt`).

---

## Why Use file-mapper?

- **One-Stop Command**: Instead of running `tree` to see the structure and then `cat` or `less` for file contents, `file-mapper` does it **all in one go**, recursively.
- **Share With AI or Peers**: Perfect for pasting a project’s structure and file contents into ChatGPT or sending it to teammates.
- **Saves Time**: No manual copy-pasting or multiple commands—filter, track only Git files, and see everything at a glance.

---

## Installation

Make sure you have Go **1.18+** installed, then run:

```bash
go install github.com/sky93/file-mapper@latest
```

This will download, build, and install the `file-mapper` binary (the main CLI for `file-mapper`) into your Go bin directory (e.g., `~/go/bin`).

> **Tip**: Ensure your `~/go/bin` (or the equivalent on your system) is in your `$PATH` so that you can run `file-mapper` from anywhere.

---

## Usage

Once installed:

```bash
file-mapper [flags...]
```

Alternatively, if you cloned the repo locally and want to run without installing:

```bash
git clone https://github.com/sky93/file-mapper.git
cd file-mapper/cmd/projectMapper
go run main.go [flags...]
```

---

### CLI Flags

| Flag                | Alias | Default | Description                                                                                                       |
|---------------------|-------|---------|-------------------------------------------------------------------------------------------------------------------|
| `--path`            | `-p`  | `.`     | Root path to scan                                                                                                 |
| `--include`         | `-i`  |         | Comma-separated file patterns to include (e.g. `--include="*.go,*.md"`)                                           |
| `--exclude`         | `-e`  |         | Comma-separated directories/files to exclude (e.g. `--exclude=".git,.idea,.env"`)                                |
| `--git`             | `-g`  | `false` | Only list files tracked by Git                                                                                    |
| `--content`         | `-c`  | `false` | Show file content (for text files)                                                                                |
| `--separate-content`|       | `false` | Print the file listing first, then all file contents afterward (instead of inline)                                |
| `--flat`            |       | `false` | Show a flat list instead of a tree                                                                                |
| `--output`          | `-o`  |         | Output file path (if not provided, prints to stdout)                                                              |
| `--line-numbers`    |       | `false` | Show line numbers for file content                                                                                |
| `--header-footer`   |       | `true`  | Print `----- CONTENT START -----` / `----- CONTENT END -----` around file content                                 |
| `--help`            |       |         | Show help message                                                                                                 |
| `--version`         |       |         | Print version information and exit                                                                                |

---

## Examples

1. **Default Tree**
   ```bash
   file-mapper
   ```
    - Displays the folder/file structure in a tree.
    - Skips hidden directories and binary files.

2. **Display File Contents Inline**
   ```bash
   file-mapper --content
   ```
    - Combines tree + `cat` for each file, printing content right below the file.

3. **Only Git-Tracked Files**
   ```bash
   file-mapper --git
   ```
    - Recursively traverses your directories, showing **only** files tracked by Git.

4. **Tree + Separate Content**
   ```bash
   file-mapper --content --separate-content
   ```
    - Prints the tree, then appends file contents at the end.

5. **Flat Listing**
   ```bash
   file-mapper --flat
   ```
    - Shows a simple list of files/directories.

6. **Include/Exclude Patterns**
   ```bash
   file-mapper --include="*.go,*.txt" --exclude=".git,.idea,node_modules"
   ```

7. **Output to a File**
   ```bash
   file-mapper --output=map.txt
   ```
    - Saves the entire listing and optional content to `map.txt`.

8. **Line Numbers + Content**
   ```bash
   file-mapper --content --line-numbers
   ```

---

## Contributing

All contributions are welcome! If you have any feature requests, bug reports, or ideas, feel free to open an [issue](https://github.com/sky93/file-mapper/issues) or submit a pull request.

---

## License

Please see the [LICENSE](https://github.com/sky93/file-mapper/blob/main/LICENSE) file for details.

---

Happy mapping with **file-mapper**! If you find it useful, give it a star on [GitHub](https://github.com/sky93/file-mapper).