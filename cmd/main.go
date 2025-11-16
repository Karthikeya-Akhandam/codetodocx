package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Karthikeya-Akhandam/codetodocx"
)

func main() {
	// Define command-line flags
	projectPath := flag.String("project", ".", "Path to the project folder to export")
	outputPath := flag.String("output", "code_export.docx", "Output Word document path")
	fullExport := flag.Bool("full", false, "Export all files (not just changed files)")
	changedOnly := flag.Bool("changed-only", false, "Export only changed files (even on first run)")
	help := flag.Bool("help", false, "Show usage information")

	flag.Parse()

	// Show help if requested
	if *help {
		fmt.Println("codetodocx - Export your code to Microsoft Word documents")
		fmt.Println("\nUsage:")
		flag.PrintDefaults()
		fmt.Println("\nExamples:")
		fmt.Println("  # First-time export (exports all tracked files in git repo)")
		fmt.Println("  codetodocx -project ./myproject -output mycode.docx")
		fmt.Println()
		fmt.Println("  # Incremental export (exports only changed files)")
		fmt.Println("  codetodocx -project ./myproject -output mycode.docx")
		fmt.Println()
		fmt.Println("  # Force full export even when updating existing document")
		fmt.Println("  codetodocx -project ./myproject -output mycode.docx -full")
		fmt.Println()
		fmt.Println("  # Force changed-only even on first run")
		fmt.Println("  codetodocx -project ./myproject -output mycode.docx -changed-only")
		os.Exit(0)
	}

	// Validate conflicting flags
	if *fullExport && *changedOnly {
		fmt.Fprintf(os.Stderr, "Error: Cannot use both -full and -changed-only flags together\n")
		os.Exit(1)
	}

	fmt.Printf("Exporting project from: %s\n", *projectPath)
	fmt.Printf("Output file: %s\n", *outputPath)

	err := codetodocx.ExportProjectWithOptions(*projectPath, *outputPath, *fullExport, *changedOnly)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error exporting project: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Word document saved to: %s\n", *outputPath)
}
