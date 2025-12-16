package sheets

import (
	"context"
	"strings"

	"voiceline-audio-backend/internal/audio"
)

type Writer struct {
	client *Client
}

func NewWriter(client *Client) *Writer {
	return &Writer{client: client}
}

func (w *Writer) WriteAudioData(ctx context.Context, data *audio.ProcessedAudio) error {
	actionItems := strings.Join(data.ActionItems, "; ")
	if actionItems == "" {
		actionItems = "None"
	}

	row := []any{
		data.Timestamp,
		data.Filename,
		data.Transcription,
		data.Summary,
		actionItems,
		"Processed",
	}

	return w.client.AppendRow(ctx, row)
}
