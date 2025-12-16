package llm

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

type Result struct {
	Summary     string   `json:"summary"`
	ActionItems []string `json:"action_items"`
}

func ParseGPTResponse(content string) (*Result, error) {
	content = strings.TrimSpace(content)
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimPrefix(content, "```")
	content = strings.TrimSuffix(content, "```")
	content = strings.TrimSpace(content)

	var result Result
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		log.Printf("parse failed: %v | content: %s", err, content)
		return nil, fmt.Errorf("parse error: %w", err)
	}

	if result.ActionItems == nil {
		result.ActionItems = []string{}
	}

	return &result, nil
}
