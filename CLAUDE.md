# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

codetodocx is a Go package that exports code to Microsoft Word documents (.docx). The core functionality intelligently detects git repositories and implements smart export behavior:
- **First-time export**: Exports ALL git-tracked files
- **Incremental export**: Exports only changed/modified files
- **Non-git directories**: Exports entire project

This solves the problem of missing the full codebase on first export while still being efficient on subsequent exports.

## Key Architecture

### Core Package Structure

- **codetodox.go**: Main library with two public functions:
  - `ExportProject(projectFolder, outputDocx string)` - Simple API with smart defaults
  - `ExportProjectWithOptions(projectFolder, outputDocx, forceFullExport, forceChangedOnly bool)` - Advanced API with manual control
- **cmd/main.go**: CLI wrapper with flags for `-project`, `-output`, `-full`, and `-changed-only`

### Smart Export Logic (codetodox.go)

The package uses a detection-based workflow with intelligent file selection:

1. **Git Detection** (TWO functions for different modes):
   - `getGitChangedFiles()`: Runs `git status --porcelain` to get M/A/?? files
   - `getAllGitTrackedFiles()`: Runs `git ls-files` to get all tracked files
   - Falls back to full directory walk if not a git repo

2. **Export Mode Selection** (in `ExportProjectWithOptions`):
   - Checks if output .docx exists to determine first-time vs incremental
   - **First-time + git repo**: Uses `getAllGitTrackedFiles()` to export full codebase
   - **Incremental + git repo**: Uses `getGitChangedFiles()` to export only changes
   - **Non-git**: Always uses directory walk
   - **Flags override**: `-full` forces all files, `-changed-only` forces changed files

3. **File Filtering** (`isTextFile`):
   - Reads first 1KB of each file to detect binary content (null bytes)
   - Rejects files with <90% printable characters
   - Skips binary extensions (.exe, .dll, .so, .dylib, .bin)
   - **Empty files** treated as text (prevents divide-by-zero bug)

4. **Document Generation** (`addCodeToDoc`) - ENHANCED FORMATTING:
   - Uses go-docx library to create Word documents (always new, never appends)
   - Horizontal separator line between files (═══════)
   - **Full file path** as bold heading with larger font (14pt)
   - **Line numbers** added to each code line (format: "   1 | code")
   - Smaller font size for code content (10pt)
   - Empty lines for spacing between sections
   - Document saved using `doc.WriteTo(file)`

## Development Commands

### Build
```bash
go build -o codetodocx ./cmd
```

### Run CLI
```bash
# Run with default options (current directory to code_export.docx)
go run ./cmd/main.go

# First-time export (exports all tracked files)
go run ./cmd/main.go -project ./myproject -output mycode.docx

# Incremental export (exports only changed files - run after first export)
go run ./cmd/main.go -project ./myproject -output mycode.docx

# Force full export even on subsequent runs
go run ./cmd/main.go -project ./myproject -output mycode.docx -full

# Force changed-only even on first run
go run ./cmd/main.go -project ./myproject -output mycode.docx -changed-only

# Show help
go run ./cmd/main.go -help
```

### Install Dependencies
```bash
go mod download
```

### Test the Package
```bash
go test ./...
```

### Format Code
```bash
go fmt ./...
```

## Dependencies

- **github.com/fumiama/go-docx**: Open-source library for Word document creation and manipulation (AGPL-3.0 license)
  - Replaced proprietary unioffice library which required a commercial license
  - Provides similar functionality for creating/parsing .docx files
- **Requires Go 1.24 or later**

## Important Implementation Details

When modifying the export logic, remember:

### Export Behavior
- The package has no tests currently - manual testing required
- Two public APIs: `ExportProject()` (smart defaults) and `ExportProjectWithOptions()` (manual control)
- Git detection failure silently falls back to full directory walk (by design)
- **Document always replaced, never appended** (go-docx limitation with parsing existing files)
- First-time detection based on whether output .docx file exists

### File Processing
- Binary file detection happens twice: once by extension check, once by content scanning
- Empty files are treated as text files (prevents divide-by-zero in `isTextFile` at line 49)
- Text file heuristic: >90% printable characters in first 1KB

### CLI and Flags
- CLI uses Go's flag package for argument parsing
- Validates that `-full` and `-changed-only` aren't used together
- Default project path is "." (current directory)
- Default output is "code_export.docx"

### Formatting Features
- Line numbers formatted as 4-digit padded with pipe separator: "   1 | "
- Font sizes: 16pt (title), 14pt (file headings), 10pt (code content)
- Horizontal separator uses box-drawing characters (═)
- Bold applied via `.Bold()` method on text runs
- Size applied via `.Size("10")` method (string parameter, not int)

### Known Limitations
- Cannot append to existing .docx files (go-docx Parse() has issues)
- No syntax highlighting (plain text only)
- No monospace font setting (go-docx Font() method requires 4 string parameters - not currently implemented)
