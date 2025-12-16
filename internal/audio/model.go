package audio

import "time"

type ProcessedAudio struct {
	Transcription string   `json:"transcription"`
	Summary       string   `json:"summary"`
	ActionItems   []string `json:"action_items"`
	Timestamp     string   `json:"timestamp"`
	Filename      string   `json:"filename"`
}

type UploadResponse struct {
	Status  string          `json:"status"`
	Message string          `json:"message"`
	Data    *ProcessedAudio `json:"data,omitempty"`
}

func NewProcessedAudio(transcription, summary string, actionItems []string, filename string) *ProcessedAudio {
	return &ProcessedAudio{
		Transcription: transcription,
		Summary:       summary,
		ActionItems:   actionItems,
		Timestamp:     time.Now().Format(time.RFC3339),
		Filename:      filename,
	}
}
