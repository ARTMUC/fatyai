package handler

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	weightapplication "github.com/artmuc/fatyai/internal/application/weight"
	"github.com/artmuc/fatyai/internal/views/partials"
)

// WeightHandler handles weight logging and correction endpoints (F5, F6).
type WeightHandler struct {
	svc *weightapplication.Service
}

// NewWeightHandler creates a new WeightHandler.
func NewWeightHandler(svc *weightapplication.Service) *WeightHandler {
	return &WeightHandler{svc: svc}
}

// LogWeight handles the morning weight form (POST /weight).
func (h *WeightHandler) LogWeight(c *gin.Context) {
	weightKg, _ := strconv.ParseFloat(c.PostForm("weight_kg"), 64)

	entry, err := h.svc.LogWeight(c.Request.Context(), weightapplication.LogWeightRequest{
		UserID:     sessionUserID(c),
		WeightKg:   weightKg,
		MeasuredAt: time.Now(),
	})
	if err != nil {
		slog.Error("log weight failed", "error", err)
		renderTempl(c, partials.ErrorMsg(err.Error()))
		return
	}

	renderTempl(c, partials.WeightRow(entry))
}

// CheckCorrection checks if a calorie correction is needed (GET /correction/check).
// Returns an empty response or the correction banner HTML.
func (h *WeightHandler) CheckCorrection(c *gin.Context) {
	dto, err := h.svc.CheckCorrection(c.Request.Context(), sessionUserID(c))
	if err != nil {
		slog.Error("check correction failed", "error", err)
		c.Status(http.StatusOK)
		return
	}
	if !dto.ShouldAdjust {
		c.Status(http.StatusOK)
		return
	}
	renderTempl(c, partials.CorrectionBanner(dto))
}

// ApplyCorrection accepts a correction and updates the calorie target (POST /correction/apply).
func (h *WeightHandler) ApplyCorrection(c *gin.Context) {
	// The correction service already knows the right delta — re-check and apply.
	dto, err := h.svc.CheckCorrection(c.Request.Context(), sessionUserID(c))
	if err != nil || !dto.ShouldAdjust {
		c.Status(http.StatusOK)
		return
	}

	// Delegate to profile service via the weight service's reference to profileRepo.
	// We expose an ApplyCorrection method on the weight service for this.
	if err := h.svc.ApplyCorrection(c.Request.Context(), sessionUserID(c), dto.SuggestedDelta); err != nil {
		slog.Error("apply correction failed", "error", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	// Return empty to clear the banner.
	c.Status(http.StatusOK)
}

// GetWeightChartData returns JSON data for the Chart.js weight chart.
func (h *WeightHandler) GetWeightChartData(c *gin.Context) {
	history, err := h.svc.GetHistory(c.Request.Context(), sessionUserID(c), 90)
	if err != nil {
		slog.Error("get weight history failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	labels := make([]string, 0, len(history.Entries))
	data := make([]float64, 0, len(history.Entries))
	for _, e := range history.Entries {
		labels = append(labels, e.MeasuredAt.Format("2006-01-02"))
		data = append(data, e.WeightKg)
	}

	c.JSON(http.StatusOK, gin.H{
		"labels":       labels,
		"data":         data,
		"trend_slope":  history.TrendSlope,
		"forecast_days": history.ForecastDays,
		"forecast_date": history.ForecastDate.Format("2006-01-02"),
	})
}
