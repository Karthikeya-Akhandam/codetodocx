package codetodocx

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fumiama/go-docx"
)

// isTextFile tries to read the first 1KB of a file and determines if it contains binary content (null bytes).
func isTextFile(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()

	buf := make([]byte, 1024)
	n, err := f.Read(buf)
	if err != nil && err != io.EOF {
		return false
	}
	buf = buf[:n]

	// Empty files are considered text files
	if len(buf) == 0 {
		return true
	}

	// Check for null bytes
	if bytes.IndexByte(buf, 0) != -1 {
		return false
	}

	// Simple heuristic: if it contains mostly printable characters, consider it text
	printableCount := 0
	for _, b := range buf {
		if (b >= 32 && b <= 126) || b == '\n' || b == '\r' || b == '\t' {
			printableCount++
		}
	}
	return printableCount*100/len(buf) > 90
}

// addCodeToDoc adds file content to the Word document with improved formatting
func addCodeToDoc(doc *docx.Docx, path string) error {
	// Add horizontal separator
	separatorPara := doc.AddParagraph()
	separatorPara.AddText("═══════════════════════════════════════════════════════════════════════════════")

	// Add file path as bold heading (full path, not just basename)
	headingPara := doc.AddParagraph()
	headingRun := headingPara.AddText(path)
	headingRun.Bold()
	headingRun.Size("14") // Slightly larger for heading

	// Add empty line for spacing
	doc.AddParagraph()

	file, err := os.Open(path)
	if err != nil {
		para := doc.AddParagraph()
		para.AddText("⚠️ Error reading file: " + err.Error())
		return err
	}
	defer file.Close()

	// Read full file content
	content, err := io.ReadAll(file)
	if err != nil {
		para := doc.AddParagraph()
		para.AddText("⚠️ Error reading file: " + err.Error())
		return err
	}

	if bytes.IndexByte(content, 0) != -1 {
		para := doc.AddParagraph()
		para.AddText("⚠️ Binary file - content not displayed")
		return nil
	}

	// Add content with line numbers
	scanner := bufio.NewScanner(bytes.NewReader(content))
	lineNum := 1
	for scanner.Scan() {
		para := doc.AddParagraph()
		// Add line with line number
		codeText := fmt.Sprintf("%4d | %s", lineNum, scanner.Text())
		codeRun := para.AddText(codeText)
		codeRun.Size("10")

		lineNum++
	}

	// Add empty line after file content
	doc.AddParagraph()

	return nil
}

// getGitChangedFiles returns a map of file paths that have git status M (modified), A (added), or U (untracked)
func getGitChangedFiles(projectFolder string) (map[string]bool, error) {
	changedFiles := make(map[string]bool)

	// Check if we're in a git repository
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = projectFolder
	output, err := cmd.Output()
	if err != nil {
		// If git command fails, return empty map (not a git repo or git not available)
		return changedFiles, nil
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if len(line) >= 3 {
			status := line[:2]
			filePath := strings.TrimSpace(line[2:])

			// Check for Modified (M), Added (A), or Untracked (??)
			if strings.Contains(status, "M") || strings.Contains(status, "A") || status == "??" {
				// Convert to absolute path
				fullPath := filepath.Join(projectFolder, filePath)
				changedFiles[fullPath] = true
			}
		}
	}

	return changedFiles, nil
}

// getAllGitTrackedFiles returns a map of all files tracked by git
func getAllGitTrackedFiles(projectFolder string) (map[string]bool, error) {
	trackedFiles := make(map[string]bool)

	// Check if we're in a git repository and get all tracked files
	cmd := exec.Command("git", "ls-files")
	cmd.Dir = projectFolder
	output, err := cmd.Output()
	if err != nil {
		// If git command fails, return empty map (not a git repo or git not available)
		return trackedFiles, nil
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			// Convert to absolute path
			fullPath := filepath.Join(projectFolder, line)
			trackedFiles[fullPath] = true
		}
	}

	return trackedFiles, nil
}

// ExportProject exports the project to a Word document
// Automatically detects if it's a git repo and exports appropriately
// First-time export: exports all git-tracked files
// Incremental export: exports only changed files
// Usage: codetodocx.ExportProject("path/to/project", "output.docx")
func ExportProject(projectFolder, outputDocx string) error {
	return ExportProjectWithOptions(projectFolder, outputDocx, false, false)
}

// ExportProjectWithOptions exports the project with custom options
// forceFullExport: if true, export all files even on incremental update
// forceChangedOnly: if true, export only changed files even on first run
func ExportProjectWithOptions(projectFolder, outputDocx string, forceFullExport, forceChangedOnly bool) error {
	var doc *docx.Docx
	var isGitRepo bool
	var filesToExport map[string]bool
	var err error
	var isFirstTime bool

	// Check if the document already exists
	_, statErr := os.Stat(outputDocx)
	isFirstTime = statErr != nil

	// Determine which files to export
	changedFiles, err := getGitChangedFiles(projectFolder)
	allTrackedFiles, err2 := getAllGitTrackedFiles(projectFolder)

	// Determine if we're in a git repo
	if (err == nil && len(changedFiles) > 0) || (err2 == nil && len(allTrackedFiles) > 0) {
		isGitRepo = true
	}

	// Decide which files to export based on flags and first-time status
	if forceFullExport {
		// Force full export requested
		if isGitRepo {
			filesToExport = allTrackedFiles
		} else {
			filesToExport = nil // Will use WalkDir for non-git repos
		}
	} else if forceChangedOnly {
		// Force changed-only export requested
		filesToExport = changedFiles
	} else {
		// Smart mode: first-time = all files, incremental = changed files
		if isFirstTime && isGitRepo {
			filesToExport = allTrackedFiles
		} else if isGitRepo {
			filesToExport = changedFiles
		} else {
			filesToExport = nil // Will use WalkDir for non-git repos
		}
	}

	// Create new document (Note: append to existing document not supported by go-docx library)
	doc = docx.New().WithDefaultTheme()

	// Determine export mode for title
	exportMode := "Full Export"
	if forceChangedOnly {
		exportMode = "Changed Files Only"
	} else if forceFullExport {
		exportMode = "All Files"
	} else if isFirstTime && isGitRepo {
		exportMode = "All Tracked Files (First Export)"
	} else if isGitRepo {
		exportMode = "Changed Files Only (Incremental)"
	}

	titlePara := doc.AddParagraph()
	titleRun := titlePara.AddText(fmt.Sprintf("Code Export from %s (%s)", filepath.Base(projectFolder), exportMode))
	titleRun.Bold()
	titleRun.Size("16")

	// Process files based on mode
	if filesToExport != nil {
		// Process specific files from the map (git mode)
		for filePath := range filesToExport {
			// Skip if it's a directory
			if info, err := os.Stat(filePath); err != nil || info.IsDir() {
				continue
			}

			// Skip binary extensions
			ext := strings.ToLower(filepath.Ext(filePath))
			if ext == ".exe" || ext == ".dll" || ext == ".so" || ext == ".dylib" || ext == ".bin" {
				continue
			}

			if isTextFile(filePath) {
				err := addCodeToDoc(doc, filePath)
				if err != nil {
					fmt.Printf("Failed to add file %s: %v\n", filePath, err)
				}
			}
		}
	} else {
		// Process entire project using directory walk (non-git mode)
		err := filepath.WalkDir(projectFolder, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}

			// Skip binary extensions
			ext := strings.ToLower(filepath.Ext(path))
			if ext == ".exe" || ext == ".dll" || ext == ".so" || ext == ".dylib" || ext == ".bin" {
				return nil
			}

			if isTextFile(path) {
				err := addCodeToDoc(doc, path)
				if err != nil {
					fmt.Printf("Failed to add file %s: %v\n", path, err)
				}
			}
			return nil
		})

		if err != nil {
			return err
		}
	}

	// Save the document
	f, err := os.Create(outputDocx)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer f.Close()

	_, err = doc.WriteTo(f)
	return err
}
