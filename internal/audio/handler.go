package audio

import (
	"net/http"

	"voiceline-audio-backend/internal/common"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) UploadAudio(c *gin.Context) {
	file, err := c.FormFile("audio")
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			http.StatusBadRequest,
			"no audio file in request",
		))
		return
	}

	result, err := h.service.ProcessAudio(c.Request.Context(), file)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, UploadResponse{
		Status:  "success",
		Message: "processed",
		Data:    result,
	})
}

func (h *Handler) handleError(c *gin.Context, err error) {
	if appErr, ok := err.(*common.AppError); ok {
		c.JSON(appErr.Code, common.NewErrorResponse(appErr.Code, appErr.Message))
		return
	}

	c.JSON(http.StatusInternalServerError, common.NewErrorResponse(
		http.StatusInternalServerError,
		"something went wrong",
	))
}
