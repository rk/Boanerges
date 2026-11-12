//go:build agentic

package agentic

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const rawRunLogName = "run.raw.log"
const rawRunLogMax = 8 << 20 // ponytail: 8MB then truncate on next read tail API

func appendRawRunLog(root, role, body string) error {
	if root == "" || body == "" {
		return nil
	}
	path := filepath.Join(root, rawRunLogName)
	if err := os.MkdirAll(root, 0o755); err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	ts := time.Now().UTC().Format(time.RFC3339)
	_, err = fmt.Fprintf(f, "\n===== %s [%s] =====\n%s\n", ts, role, body)
	return err
}

func readRawRunLogTail(root string, maxBytes int) (string, error) {
	if root == "" {
		return "", fmt.Errorf("folder required")
	}
	path := filepath.Join(root, rawRunLogName)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	if maxBytes <= 0 {
		maxBytes = 200_000
	}
	if len(data) > maxBytes {
		data = data[len(data)-maxBytes:]
		return "(…truncated…)\n" + string(data), nil
	}
	return string(data), nil
}

func truncateRaw(s string, max int) string {
	s = strings.TrimSpace(s)
	if len(s) <= max {
		return s
	}
	return s[:max] + "\n…"
}
