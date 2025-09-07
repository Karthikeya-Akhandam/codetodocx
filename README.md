# CodeTodox

A smart Go package that exports your code to Microsoft Word documents. CodeTodox automatically detects if you're in a git repository and exports only changed files, or exports the entire project if it's not a git repository.

## Features

- 🚀 **Smart Detection**: Automatically detects git repositories and exports only changed files
- 📝 **Word Document Export**: Creates professional Word documents with your code
- 🔄 **Incremental Updates**: Appends to existing documents instead of overwriting
- 🎯 **Git Integration**: Exports files with git status M (modified), A (added), or U (untracked)
- 📁 **File Filtering**: Automatically skips binary files and detects text files
- 💪 **Bold Headings**: File names appear as bold headings in the document

## Installation

```bash
go get github.com/Karthikeya-Akhandam/codetodox
```

## Usage

### Basic Usage

```go
package main

import (
    "fmt"
    "github.com/Karthikeya-Akhandam/codetodox"
)

func main() {
    // One function call - automatically does the smart thing!
    err := codetodox.ExportProject("./myproject", "mycode.docx")
    if err != nil {
        fmt.Println("Error:", err)
    }
}
```

### How It Works

The `ExportProject()` function is intelligent and automatically:

1. **Git Repository**: Exports only changed files (M/A/U status)
2. **Regular Folder**: Exports entire project
3. **Existing Document**: Appends to existing Word document
4. **New Document**: Creates fresh Word document

### Example Output

**For Git Repositories:**
- Document Title: "Code Export - Changed Files from [project]"
- Updates: "--- Updated Export (Changed Files Only) ---"

**For Regular Folders:**
- Document Title: "Code Export from [project]"
- Updates: "--- Updated Export ---"

## Git Integration

CodeTodox automatically detects git repositories and exports files with the following status:

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

## CLI Example

You can also use it as a command-line tool:

```bash
# Build the CLI
go build ./cmd

# Run it
./cmd
```

## Requirements

- Go 1.24.4 or later
- Git (optional, for git integration)
- Microsoft Word or compatible software to view the output

## Dependencies

- `github.com/unidoc/unioffice` - For Word document generation

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Author

Created by [Karthikeya Akhandam](https://github.com/Karthikeya-Akhandam)

---

**CodeTodox** - Smart code export to Word documents! 📄✨