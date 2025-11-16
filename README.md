# codetodocx

A smart Go package that exports your code to Microsoft Word documents. codetodocx automatically detects if you're in a git repository and exports only changed files, or exports the entire project if it's not a git repository.

## Features

- 🚀 **Smart Detection**: First-time export includes all tracked files, incremental exports only changed files
- 📝 **Word Document Export**: Creates professional Word documents with your code
- 🔢 **Line Numbers**: Automatically adds line numbers to all code
- 🎯 **Git Integration**: Exports files with git status M (modified), A (added), or U (untracked)
- 📁 **File Filtering**: Automatically skips binary files and detects text files
- 💪 **Improved Formatting**: Bold headings with full file paths, horizontal separators between files
- ⚙️ **Flexible Options**: Command-line flags for full or changed-only exports

## Installation

```bash
go get github.com/Karthikeya-Akhandam/codetodocx
```

## Usage

### Basic Usage

```go
package main

import (
    "fmt"
    "github.com/Karthikeya-Akhandam/codetodocx"
)

func main() {
    // One function call - automatically does the smart thing!
    err := codetodocx.ExportProject("./myproject", "mycode.docx")
    if err != nil {
        fmt.Println("Error:", err)
    }
}
```

### How It Works

The `ExportProject()` function is intelligent and automatically:

1. **First-Time Git Export**: Exports ALL git-tracked files
2. **Incremental Git Export**: Exports only changed files (M/A/U status)
3. **Regular Folder**: Exports entire project
4. **Document Creation**: Always creates a fresh Word document (overwrites if exists)

### Example Output

**For Git Repositories (First Time):**
- Document Title: "Code Export from [project] (All Tracked Files - First Export)"

**For Git Repositories (Incremental):**
- Document Title: "Code Export from [project] (Changed Files Only - Incremental)"

**For Regular Folders:**
- Document Title: "Code Export from [project] (Full Export)"

## Git Integration

codetodocx automatically detects git repositories and exports files with the following status:

- **M** - Modified files
- **A** - Added files  
- **U** - Untracked files

If you're not in a git repository, it exports all text files in the project.

## File Filtering

The package automatically:

- ✅ Includes text files (source code, config files, etc.)
- ❌ Skips binary files (.exe, .dll, .so, .dylib, .bin)
- ❌ Skips directories
- ❌ Skips files with null bytes (binary detection)

## CLI Usage

Build and use as a command-line tool:

```bash
# Build the CLI
go build -o codetodocx ./cmd

# First-time export (exports all tracked files in git repo)
./codetodocx -project ./myproject -output mycode.docx

# Incremental export (exports only changed files)
./codetodocx -project ./myproject -output mycode.docx

# Force full export (export all files even on subsequent runs)
./codetodocx -project ./myproject -output mycode.docx -full

# Force changed-only export (export only changed files even on first run)
./codetodocx -project ./myproject -output mycode.docx -changed-only

# Show help
./codetodocx -help
```

## Requirements

- Go 1.24.4 or later
- Git (optional, for git integration)
- Microsoft Word or compatible software to view the output

## Dependencies

- `github.com/fumiama/go-docx` - Open-source library for Word document generation (AGPL-3.0)

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Author

Created by [Karthikeya Akhandam](https://github.com/Karthikeya-Akhandam)

---

**codetodocx** - Smart code export to Word documents! 📄✨