// internal/api/api.go

package api

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const httpTimeout = 10 * time.Second

func CheckAPIReachability(ctx context.Context, apiURL string) error {
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

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("API is not reachable, status code: %d", resp.StatusCode)
	}

	return nil
}
