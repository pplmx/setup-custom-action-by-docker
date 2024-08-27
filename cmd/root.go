package cmd

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/pplmx/setup-my-action/internal/api"
	"github.com/pplmx/setup-my-action/internal/config"
	"github.com/pplmx/setup-my-action/internal/file"
	"github.com/pplmx/setup-my-action/internal/output"
	"github.com/pplmx/setup-my-action/internal/text"
)

const contextTimeout = 15 * time.Second

// Execute is the main entry point for the cmd package, orchestrating the main program flow.
func Execute() {
	configPath := os.Getenv("INPUT_CONFIG_PATH")
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Process text
	if err := processText(cfg); err != nil {
		log.Fatalf("Text processing error: %v", err)
	}

	// Process numbers
	if err := processNumbers(cfg); err != nil {
		log.Fatalf("Number processing error: %v", err)
	}

	// Process files
	if err := processFiles(cfg); err != nil {
		log.Fatalf("File processing error: %v", err)
	}

	// Check API reachability
	if err := checkAPI(cfg); err != nil {
		log.Fatalf("API check error: %v", err)
	}
}

func processText(cfg *config.Config) error {
	processedText, wordCount := text.ProcessText(cfg.InputText, cfg.FindWord, cfg.ReplaceWord)
	if err := output.SetOutput("processed_text", processedText); err != nil {
		return err
	}
	return output.SetOutput("word_count", strconv.Itoa(wordCount))
}

func processNumbers(cfg *config.Config) error {
	sum, average := text.CalculateSumAndAverage(cfg.NumberList)
	if err := output.SetOutput("sum", strconv.FormatFloat(sum, 'f', -1, 64)); err != nil {
		return err
	}
	return output.SetOutput("average", strconv.FormatFloat(average, 'f', -1, 64))
}

func processFiles(cfg *config.Config) error {
	return file.ReadAndAppendToFile(cfg.InputFile, cfg.OutputFile, cfg.AppendText)
}

func checkAPI(cfg *config.Config) error {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	if err := api.CheckAPIReachability(ctx, cfg.ApiURL); err != nil {
		return err
	}
	return output.SetOutput("response_field", "API Reachable")
}
