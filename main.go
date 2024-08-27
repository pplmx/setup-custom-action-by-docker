package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

type Config struct {
	InputText   string    `toml:"input_text"`
	FindWord    string    `toml:"find_word"`
	ReplaceWord string    `toml:"replace_word"`
	NumberList  []float64 `toml:"number_list"`
	InputFile   string    `toml:"input_file"`
	OutputFile  string    `toml:"output_file"`
	AppendText  string    `toml:"append_text"`
	ApiURL      string    `toml:"api_url"`
}

const (
	defaultConfigPath = ".github/configs/setup-my-action.toml"
	httpTimeout       = 10 * time.Second
	contextTimeout    = 15 * time.Second
)

func loadConfig(path string) (*Config, error) {
	if path == "" {
		path = defaultConfigPath
	}

	var config Config
	if _, err := toml.DecodeFile(path, &config); err != nil {
		return nil, fmt.Errorf("failed to decode TOML file: %w", err)
	}

	// Validate config
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}

func validateConfig(config *Config) error {
	if config.InputText == "" {
		return fmt.Errorf("input_text is required")
	}
	if config.InputFile == "" || config.OutputFile == "" {
		return fmt.Errorf("input_file and output_file are required")
	}
	if config.ApiURL == "" {
		return fmt.Errorf("api_url is required")
	}
	return nil
}

func processText(inputText, findWord, replaceWord string) (string, int) {
	processedText := strings.ReplaceAll(inputText, findWord, replaceWord)
	wordCount := len(strings.Fields(inputText))
	return processedText, wordCount
}

func calculateSumAndAverage(numbers []float64) (float64, float64) {
	sum := 0.0
	for _, num := range numbers {
		sum += num
	}
	average := sum / float64(len(numbers))
	return sum, average
}

func readAndAppendToFile(inputFile, outputFile, appendText string) error {
	content, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	modifiedContent := string(content) + "\n" + appendText

	return os.WriteFile(outputFile, []byte(modifiedContent), 0644)
}

func callAPIAndExtractField(ctx context.Context, apiURL string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{Timeout: httpTimeout}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make API request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Failed to close response body: %v", err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read API response: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to parse JSON response: %w", err)
	}

	dataField, ok := result["data"].(string)
	if !ok {
		return "", fmt.Errorf("field 'data' not found in API response or is not a string")
	}

	return dataField, nil
}

func setOutput(name, value string) {
	fmt.Printf("::set-output name=%s::%s\n", name, value)
}

func main() {
	configPath := os.Getenv("INPUT_CONFIG_PATH")
	config, err := loadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// 1. 文本处理
	processedText, wordCount := processText(config.InputText, config.FindWord, config.ReplaceWord)
	setOutput("processed_text", processedText)
	setOutput("word_count", strconv.Itoa(wordCount))

	// 2. 列表处理
	sum, average := calculateSumAndAverage(config.NumberList)
	setOutput("sum", strconv.FormatFloat(sum, 'f', -1, 64))
	setOutput("average", strconv.FormatFloat(average, 'f', -1, 64))

	// 3. 文件处理
	if err := readAndAppendToFile(config.InputFile, config.OutputFile, config.AppendText); err != nil {
		log.Fatalf("File processing error: %v", err)
	}

	// 4. API 请求处理
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	responseField, err := callAPIAndExtractField(ctx, config.ApiURL)
	if err != nil {
		log.Fatalf("API request error: %v", err)
	}
	setOutput("response_field", responseField)
}
