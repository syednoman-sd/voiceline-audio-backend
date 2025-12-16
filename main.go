package main

import (
	"context"
	"log"

	"voiceline-audio-backend/internal/audio"
	"voiceline-audio-backend/internal/common"
	"voiceline-audio-backend/internal/llm"
	"voiceline-audio-backend/internal/middleware"
	"voiceline-audio-backend/internal/sheets"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := common.LoadConfig(); err != nil {
		log.Fatalf("config load failed: %v", err)
	}

	gin.SetMode(common.AppConfig.GinMode)
	r := gin.Default()

	r.Use(middleware.CORS())
	r.Use(middleware.Logger())

	llmClient := llm.NewClient(common.AppConfig.OpenAIAPIKey)

	var sheetsWriter audio.SheetsWriter
	if common.AppConfig.GoogleCredsFile != "" && common.AppConfig.GoogleSheetID != "" {
		sheetsClient, err := sheets.NewClient(
			common.AppConfig.GoogleCredsFile,
			common.AppConfig.GoogleSheetID,
		)
		if err != nil {
			log.Printf("sheets init failed: %v", err)
		} else {
			ctx := context.Background()
			sheetsClient.CreateHeaderRow(ctx)
			sheetsWriter = sheets.NewWriter(sheetsClient)
		}
	}

	validator := audio.NewValidator(common.AppConfig.MaxAudioSizeMB)
	service := audio.NewService(llmClient, validator, sheetsWriter)
	handler := audio.NewHandler(service)

	api := r.Group("/api")
	{
		audioRoutes := api.Group("/audio")
		{
			audioRoutes.POST("/upload", handler.UploadAudio)
		}
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	port := ":" + common.AppConfig.Port
	log.Printf("server starting on port %s", common.AppConfig.Port)

	if err := r.Run(port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
