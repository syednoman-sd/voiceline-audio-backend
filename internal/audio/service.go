package audio

import (
	"context"
	"log"
	"mime/multipart"

	"voiceline-audio-backend/internal/llm"
)

type SheetsWriter interface {
	WriteAudioData(ctx context.Context, data *ProcessedAudio) error
}

type Service struct {
	llmClient    *llm.Client
	validator    *Validator
	sheetsWriter SheetsWriter
}

func NewService(llmClient *llm.Client, validator *Validator, sheetsWriter SheetsWriter) *Service {
	return &Service{
		llmClient:    llmClient,
		validator:    validator,
		sheetsWriter: sheetsWriter,
	}
}

func (s *Service) ProcessAudio(ctx context.Context, fileHeader *multipart.FileHeader) (*ProcessedAudio, error) {
	if err := s.validator.ValidateFile(fileHeader); err != nil {
		return nil, err
	}

	file, err := fileHeader.Open()
	if err != nil {
		log.Printf("couldn't open file: %v", err)
		return nil, err
	}
	defer file.Close()

	transcription, err := s.llmClient.TranscribeAudio(ctx, file, fileHeader.Filename)
	if err != nil {
		return nil, err
	}

	result, err := s.llmClient.ProcessTranscription(ctx, transcription)
	if err != nil {
		return nil, err
	}

	processedAudio := NewProcessedAudio(
		transcription,
		result.Summary,
		result.ActionItems,
		fileHeader.Filename,
	)

	if s.sheetsWriter != nil {
		s.sheetsWriter.WriteAudioData(ctx, processedAudio)
	}

	return processedAudio, nil
}
