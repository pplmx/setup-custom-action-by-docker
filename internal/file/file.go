// internal/file/file.go

package file

import (
	"fmt"
	"os"
)

func ReadAndAppendToFile(inputFile, outputFile, appendText string) error {
	content, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	modifiedContent := string(content) + "\n" + appendText

	return os.WriteFile(outputFile, []byte(modifiedContent), 0644)
}
