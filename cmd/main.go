package main

import (
	"fmt"

	"github.com/Karthikeya-Akhandam/codetodox"
)

func main() {
	projectFolderPath := ".."                          // Project folder
	outputDocxPath := "kisanlink_erp_code_export.docx" // Output file

	fmt.Println("Exporting project...")
	err := codetodox.ExportProject(projectFolderPath, outputDocxPath)
	if err != nil {
		fmt.Println("Error exporting project:", err)
	} else {
		fmt.Println("✅ Word document saved to:", outputDocxPath)
	}
}
