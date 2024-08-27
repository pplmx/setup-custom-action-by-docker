package main

import (
	"context"
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

func checkAPIReachability(ctx context.Context, apiURL string) error {
	// 创建带有上下文的 HTTP GET 请求
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{Timeout: httpTimeout}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make API request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Failed to close response body: %v", err)
		}
	}(resp.Body)

	// 检查响应状态码是否在 200-299 范围内
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("API is not reachable, status code: %d", resp.StatusCode)
	}

	return nil
}

func setOutput(name, value string) error {
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

func main() {
	configPath := os.Getenv("INPUT_CONFIG_PATH")
	config, err := loadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// 1. 文本处理
	processedText, wordCount := processText(config.InputText, config.FindWord, config.ReplaceWord)
	if err := setOutput("processed_text", processedText); err != nil {
		log.Fatalf("Failed to set output: %v", err)
	}
	if err := setOutput("word_count", strconv.Itoa(wordCount)); err != nil {
		log.Fatalf("Failed to set output: %v", err)
	}

	// 2. 列表处理
	sum, average := calculateSumAndAverage(config.NumberList)
	if err := setOutput("sum", strconv.FormatFloat(sum, 'f', -1, 64)); err != nil {
		log.Fatalf("Failed to set output: %v", err)
	}
	if err := setOutput("average", strconv.FormatFloat(average, 'f', -1, 64)); err != nil {
		log.Fatalf("Failed to set output: %v", err)
	}

	// 3. 文件处理
	if err := readAndAppendToFile(config.InputFile, config.OutputFile, config.AppendText); err != nil {
		log.Fatalf("File processing error: %v", err)
	}

	// 4. API 请求处理
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	err = checkAPIReachability(ctx, config.ApiURL)
	if err != nil {
		log.Fatalf("API request error: %v", err)
	}
	if err := setOutput("response_field", "API Reachable"); err != nil {
		log.Fatalf("Failed to set output: %v", err)
	}
}
