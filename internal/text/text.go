// internal/text/text.go

package text

import "strings"

func ProcessText(inputText, findWord, replaceWord string) (string, int) {
	processedText := strings.ReplaceAll(inputText, findWord, replaceWord)
	wordCount := len(strings.Fields(inputText))
	return processedText, wordCount
}

func CalculateSumAndAverage(numbers []float64) (float64, float64) {
	sum := 0.0
	for _, num := range numbers {
		sum += num
	}
	average := sum / float64(len(numbers))
	return sum, average
}
