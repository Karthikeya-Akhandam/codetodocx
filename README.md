# codetodocx

[![Go Reference](https://pkg.go.dev/badge/github.com/Karthikeya-Akhandam/codetodocx.svg)](https://pkg.go.dev/github.com/Karthikeya-Akhandam/codetodocx)
[![Go Report Card](https://goreportcard.com/badge/github.com/Karthikeya-Akhandam/codetodocx)](https://goreportcard.com/report/github.com/Karthikeya-Akhandam/codetodocx)
[![GitHub release](https://img.shields.io/github/release/Karthikeya-Akhandam/codetodocx.svg)](https://github.com/Karthikeya-Akhandam/codetodocx/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**Get AI code reviews from ChatGPT without expensive subscriptions!**

Export your entire codebase to a single Word document and upload it to ChatGPT for instant code reviews, bug detection, and suggestions. No more manual copy-pasting multiple files or paying for Cursor/Claude Code subscriptions.

## Why codetodocx?

### 🤖 The Problem

You want AI code reviews, but:
- ❌ Claude Code, Cursor, and GitHub Copilot cost $20-40/month
- ❌ Manually copying/pasting code to ChatGPT is tedious and time-consuming
- ❌ Multiple prompts waste tokens and lose context
- ❌ Free ChatGPT can't access your codebase directly

### ✅ The Solution

**One command. One file. Complete code review.**

```bash
codetodocx -project . -output my_code.docx
```

Upload `my_code.docx` to ChatGPT and ask:
- "Review this code for bugs and suggest improvements"
- "Find security vulnerabilities in this codebase"
- "Explain how this project works"
- "Suggest refactoring opportunities"

### 💰 Why This Matters

**For Small Developers & Students:**
- ✅ **Free**: Use ChatGPT free tier instead of $20/month AI subscriptions
- ✅ **Fast**: One upload vs. dozens of copy-pastes
- ✅ **Smart**: Automatically exports only changed files on updates
- ✅ **Organized**: Line numbers, file paths, clean formatting

## Features

- 🚀 **Smart Detection**: First-time export includes all tracked files, incremental exports only changed files
- 📝 **Word Document Export**: Creates professional Word documents with your code
- 🔢 **Line Numbers**: Automatically adds line numbers to all code
- 🎯 **Git Integration**: Exports files with git status M (modified), A (added), or U (untracked)
- 📁 **File Filtering**: Automatically skips binary files and detects text files
- 💪 **Improved Formatting**: Bold headings with full file paths, horizontal separators between files
- ⚙️ **Flexible Options**: Command-line flags for full or changed-only exports

## Use Cases

### Primary Use Case: ChatGPT Code Reviews

1. 🤖 **Get AI Code Reviews** - Upload one Word doc to ChatGPT instead of expensive AI coding subscriptions
2. 💰 **Save Money** - Free alternative to Cursor ($20/mo), Claude Code ($20/mo), GitHub Copilot ($10-40/mo)
3. ⚡ **Save Time** - Export entire project in seconds vs. manual copy-paste of every file
4. 📊 **Token Efficient** - One organized document preserves context better than scattered prompts
5. 🔄 **Incremental Updates** - Export only changed files for quick re-reviews

### Other Use Cases

- 📋 **Code Reviews**: Share formatted code with human reviewers
- 📚 **Documentation**: Create technical documentation with actual source code
- 🎓 **Education**: Prepare teaching materials with properly formatted code examples
- 💼 **Client Deliverables**: Professional code submissions for clients
- 📊 **Portfolio**: Present your projects in a readable format

## Installation

```bash
go get github.com/Karthikeya-Akhandam/codetodocx
```

Or install the latest version:

```bash
go get github.com/Karthikeya-Akhandam/codetodocx@latest
```

For more information, visit [pkg.go.dev](https://pkg.go.dev/github.com/Karthikeya-Akhandam/codetodocx).

## Quick Start: ChatGPT Code Review in 3 Steps

### Step 1: Export Your Code

```bash
# Build the CLI
go build -o codetodocx ./cmd

# Export your project
./codetodocx -project . -output my_code.docx
```

### Step 2: Upload to ChatGPT

1. Go to [ChatGPT](https://chat.openai.com)
2. Click the attachment icon (📎)
3. Upload `my_code.docx`

### Step 3: Ask for Review

**Example Prompts:**

```
Review this code and suggest improvements
```

```
Find potential bugs and security vulnerabilities in this codebase
```

```
Explain the architecture and how the components interact
```

```
Suggest refactoring opportunities to improve code quality
```

```
Check for performance issues and optimization opportunities
```

### Step 4 (Optional): Update After Changes

```bash
# Export only changed files (faster, smaller file)
./codetodocx -project . -output my_code.docx

# Upload the updated file to ChatGPT
# Ask: "Review the changes I made"
```

---

## Usage

### Basic Usage (Programmatic)

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

- Go 1.24 or later
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