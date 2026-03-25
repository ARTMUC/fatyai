package handler

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	activityapplication "github.com/artmuc/fatyai/internal/application/activity"
	mealapplication "github.com/artmuc/fatyai/internal/application/meal"
	"github.com/artmuc/fatyai/internal/views/pages"
)

// ChartsHandler serves the charts page and its JSON data endpoints (F6).
type ChartsHandler struct {
	mealSvc     *mealapplication.Service
	activitySvc *activityapplication.Service
}

// NewChartsHandler creates a new ChartsHandler.
func NewChartsHandler(mealSvc *mealapplication.Service, activitySvc *activityapplication.Service) *ChartsHandler {
	return &ChartsHandler{mealSvc: mealSvc, activitySvc: activitySvc}
}

// ShowCharts renders the charts page (GET /charts).
func (h *ChartsHandler) ShowCharts(c *gin.Context) {
	renderTempl(c, pages.ChartsPage())
}

// GetCalorieChartData returns last-30-days calorie totals as JSON.
func (h *ChartsHandler) GetCalorieChartData(c *gin.Context) {
	labels := make([]string, 0, 30)
	data := make([]float64, 0, 30)

	for i := 29; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i)
		journal, err := h.mealSvc.GetDailyJournal(c.Request.Context(), sessionUserID(c), date)
		if err != nil {
			slog.Error("get daily journal for chart failed", "date", date, "error", err)
			continue
		}
		labels = append(labels, date.Format("01/02"))
		data = append(data, journal.TotalCalories)
	}

	c.JSON(http.StatusOK, gin.H{"labels": labels, "data": data})
}

// GetNetCalorieChartData returns last-30-days net calories (food minus exercise) as JSON.
func (h *ChartsHandler) GetNetCalorieChartData(c *gin.Context) {
	labels := make([]string, 0, 30)
	data := make([]float64, 0, 30)

	for i := 29; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i)
		journal, err := h.mealSvc.GetDailyJournal(c.Request.Context(), sessionUserID(c), date)
		if err != nil {
			slog.Error("get daily journal for net chart failed", "date", date, "error", err)
			continue
		}
		burned, err := h.activitySvc.GetDailyBurnedCalories(c.Request.Context(), sessionUserID(c), date)
		if err != nil {
			slog.Error("get burned calories for net chart failed", "date", date, "error", err)
			burned = 0
		}
		labels = append(labels, date.Format("01/02"))
		data = append(data, journal.TotalCalories-burned)
	}

	c.JSON(http.StatusOK, gin.H{"labels": labels, "data": data})
}

// GetMacroChartData returns today's macro breakdown as JSON.
func (h *ChartsHandler) GetMacroChartData(c *gin.Context) {
	journal, err := h.mealSvc.GetDailyJournal(c.Request.Context(), sessionUserID(c), time.Now())
	if err != nil {
		slog.Error("get macro chart data failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"protein": journal.TotalProteinG,
		"fat":     journal.TotalFatG,
		"carbs":   journal.TotalCarbsG,
	})
}
