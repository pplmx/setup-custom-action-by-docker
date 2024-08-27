// internal/output/output.go

package output

import (
	"fmt"
	"log"
	"os"
)

func SetOutput(name, value string) error {
	envFile := os.Getenv("GITHUB_OUTPUT")
	f, err := os.OpenFile(envFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to open GITHUB_OUTPUT file: %w", err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Printf("Failed to close GITHUB_OUTPUT file: %v", err)
		}
	}(f)

	_, err = fmt.Fprintf(f, "%s=%s\n", name, value)
	if err != nil {
		return fmt.Errorf("failed to write to GITHUB_OUTPUT file: %w", err)
	}
	return nil
}
