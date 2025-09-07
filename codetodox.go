package codetodox

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

	"github.com/unidoc/unioffice/document"
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

// addCodeToDoc adds file content to the Word document
func addCodeToDoc(doc *document.Document, path string) error {
	// Add file name as bold heading
	headingPara := doc.AddParagraph()
	headingRun := headingPara.AddRun()
	headingRun.Properties().SetBold(true)
	headingRun.AddText(filepath.Base(path))
	para := doc.AddParagraph()

	file, err := os.Open(path)
	if err != nil {
		para.AddRun().AddText("⚠️ Error reading file: " + err.Error())
		return err
	}
	defer file.Close()

	// Read full file content
	content, err := io.ReadAll(file)
	if err != nil {
		para.AddRun().AddText("⚠️ Error reading file: " + err.Error())
		return err
	}

	if bytes.IndexByte(content, 0) != -1 {
		para.AddRun().AddText("⚠️ Binary file - content not displayed")
		return nil
	}

	// Add content as plain text (no syntax highlighting)
	scanner := bufio.NewScanner(bytes.NewReader(content))
	for scanner.Scan() {
		para.AddRun().AddText(scanner.Text())
		para = doc.AddParagraph() // New paragraph for each line to preserve line breaks
	}

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

// ExportProject exports the project to a Word document
// Automatically detects if it's a git repo and exports only changed files, otherwise exports entire project
// Usage: codetodox.ExportProject("path/to/project", "output.docx")
func ExportProject(projectFolder, outputDocx string) error {
	var doc *document.Document
	var isGitRepo bool
	var changedFiles map[string]bool
	var err error

	// Check if we're in a git repository and get changed files
	changedFiles, err = getGitChangedFiles(projectFolder)
	if err == nil && len(changedFiles) > 0 {
		isGitRepo = true
	}

	// Check if the document already exists
	if _, err := os.Stat(outputDocx); err == nil {
		// File exists, open it
		doc, err = document.Open(outputDocx)
		if err != nil {
			return fmt.Errorf("failed to open existing document: %v", err)
		}
		if isGitRepo {
			doc.AddParagraph().AddRun().AddText("--- Updated Export (Changed Files Only) ---")
		} else {
			doc.AddParagraph().AddRun().AddText("--- Updated Export ---")
		}
	} else {
		// File doesn't exist, create new document
		doc = document.New()
		if isGitRepo {
			doc.AddParagraph().AddRun().AddText("Code Export - Changed Files from " + filepath.Base(projectFolder))
		} else {
			doc.AddParagraph().AddRun().AddText("Code Export from " + filepath.Base(projectFolder))
		}
	}

	if isGitRepo {
		// Process only the changed files
		for filePath := range changedFiles {
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
		// Process entire project
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

	return doc.SaveToFile(outputDocx)
}
