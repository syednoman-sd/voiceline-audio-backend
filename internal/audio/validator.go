package audio

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"voiceline-audio-backend/internal/common"
)

type Validator struct {
	maxSizeMB      int64
	allowedFormats []string
}

func NewValidator(maxSizeMB int64) *Validator {
	return &Validator{
		maxSizeMB: maxSizeMB,
		allowedFormats: []string{
			".mp3", ".wav", ".m4a", ".webm", ".ogg",
			".flac", ".aac", ".mp4", ".mpeg", ".mpga",
		},
	}
}

func (v *Validator) ValidateFile(file *multipart.FileHeader) error {
	maxBytes := v.maxSizeMB * 1024 * 1024
	if file.Size > maxBytes {
		return common.NewBadRequestError(
			fmt.Sprintf("file too large (max %dMB)", v.maxSizeMB),
		)
	}

	if file.Size == 0 {
		return common.NewBadRequestError("empty file")
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !v.isAllowedFormat(ext) {
		return common.NewBadRequestError(
			fmt.Sprintf("unsupported format: %s", ext),
		)
	}

	return nil
}

func (v *Validator) isAllowedFormat(ext string) bool {
	for _, format := range v.allowedFormats {
		if ext == format {
			return true
		}
	}
	return false
}
