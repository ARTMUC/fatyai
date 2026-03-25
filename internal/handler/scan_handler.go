package handler

import (
	"encoding/base64"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	mealapplication "github.com/artmuc/fatyai/internal/application/meal"
	"github.com/artmuc/fatyai/internal/views/partials"
)

// ScanHandler handles the AI meal scanning endpoint (F1).
type ScanHandler struct {
	mealSvc *mealapplication.Service
}

// NewScanHandler creates a new ScanHandler.
func NewScanHandler(svc *mealapplication.Service) *ScanHandler {
	return &ScanHandler{mealSvc: svc}
}

// Scan reads a photo upload, calls GPT-4o Vision, and returns the confirm form partial.
// POST /meals/scan — multipart/form-data with field "photo".
func (h *ScanHandler) Scan(c *gin.Context) {
	file, header, err := c.Request.FormFile("photo")
	if err != nil {
		renderTempl(c, partials.ErrorMsg("Could not read photo: "+err.Error()))
		return
	}
	defer file.Close()

	// Detect MIME type from the Content-Type header of the part, fall back to jpeg.
	mimeType := header.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "image/jpeg"
	}

	// Limit to 10 MB.
	const maxBytes = 10 << 20
	data, err := io.ReadAll(io.LimitReader(file, maxBytes))
	if err != nil {
		renderTempl(c, partials.ErrorMsg("Could not read file: "+err.Error()))
		return
	}

	imageBase64 := base64.StdEncoding.EncodeToString(data)
	draft, err := h.mealSvc.ScanMeal(c.Request.Context(), mealapplication.ScanMealRequest{
		UserID:      sessionUserID(c),
		ImageBase64: imageBase64,
		MimeType:    mimeType,
	})
	if err != nil {
		renderTempl(c, partials.ErrorMsg("AI analysis failed: "+err.Error()))
		return
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	if err := partials.MealConfirmForm(draft).Render(c.Request.Context(), c.Writer); err != nil {
		c.Status(http.StatusInternalServerError)
	}
}
