//go:build agentic

package agentic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const defaultOllamaBase = "http://127.0.0.1:11434"

func PingOllama(ctx context.Context) (OllamaHealth, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, defaultOllamaBase+"/api/tags", nil)
	if err != nil {
		return OllamaHealth{OK: false, Error: err.Error()}, nil
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return OllamaHealth{OK: false, Error: err.Error()}, nil
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return OllamaHealth{OK: false, Error: fmt.Sprintf("ollama status %d", resp.StatusCode)}, nil
	}
	var body struct {
		Models []struct {
			Name string `json:"name"`
		} `json:"models"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return OllamaHealth{OK: false, Error: err.Error()}, nil
	}
	names := make([]string, 0, len(body.Models))
	for _, m := range body.Models {
		names = append(names, m.Name)
	}
	return OllamaHealth{OK: true, Models: names}, nil
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatRequest struct {
	Model    string          `json:"model"`
	Messages []chatMessage   `json:"messages"`
	Stream   bool            `json:"stream"`
	Format   json.RawMessage `json:"format,omitempty"`
}

type chatResponse struct {
	Message chatMessage `json:"message"`
}

func ollamaChat(ctx context.Context, model, systemPrompt, userPrompt string) (string, error) {
	// ponytail: JSON mode often yields empty/tiny replies on some local models; plain text + parseToolTurn instead
	return ollamaChatFormatted(ctx, model, systemPrompt, userPrompt, false)
}

func ollamaChatFormatted(ctx context.Context, model, systemPrompt, userPrompt string, jsonMode bool) (string, error) {
	if model == "" {
		model = "llama3.2"
	}
	payload := chatRequest{
		Model: model,
		Messages: []chatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		Stream: false,
	}
	if jsonMode {
		payload.Format = json.RawMessage(`"json"`)
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, defaultOllamaBase+"/api/chat", bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 12 * time.Minute}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ollama chat %d: %s", resp.StatusCode, string(raw))
	}
	var out chatResponse
	if err := json.Unmarshal(raw, &out); err != nil {
		return "", err
	}
	return out.Message.Content, nil
}
