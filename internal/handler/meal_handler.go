package handler

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	mealapplication "github.com/artmuc/fatyai/internal/application/meal"
	"github.com/artmuc/fatyai/internal/views/partials"
)

// MealHandler handles CRUD endpoints for meals (F1, F2).
type MealHandler struct {
	svc *mealapplication.Service
}

// NewMealHandler creates a new MealHandler.
func NewMealHandler(svc *mealapplication.Service) *MealHandler {
	return &MealHandler{svc: svc}
}

// LogMeal handles the confirmed meal form submission (POST /meals).
func (h *MealHandler) LogMeal(c *gin.Context) {
	calories, _ := strconv.ParseFloat(c.PostForm("calories_kcal"), 64)
	protein, _ := strconv.ParseFloat(c.PostForm("protein_g"), 64)
	fat, _ := strconv.ParseFloat(c.PostForm("fat_g"), 64)
	carbs, _ := strconv.ParseFloat(c.PostForm("carbs_g"), 64)
	weightG, _ := strconv.ParseFloat(c.PostForm("weight_g"), 64)
	estimated := c.PostForm("estimated") == "true"

	var eatenAt time.Time
	if raw := c.PostForm("eaten_at"); raw != "" {
		eatenAt, _ = time.Parse("2006-01-02T15:04", raw)
	}
	if eatenAt.IsZero() {
		eatenAt = time.Now()
	}

	req := mealapplication.LogMealRequest{
		UserID:       sessionUserID(c),
		Name:         c.PostForm("name"),
		CaloriesKcal: calories,
		ProteinG:     protein,
		FatG:         fat,
		CarbsG:       carbs,
		WeightG:      weightG,
		Estimated:    estimated,
		Source:       c.PostForm("source"),
		EatenAt:      eatenAt,
	}

	meal, err := h.svc.LogMeal(c.Request.Context(), req)
	if err != nil {
		slog.Error("log meal failed", "error", err)
		renderTempl(c, partials.ErrorMsg(err.Error()))
		return
	}

	renderTempl(c, partials.MealRow(meal))
}

// EditMeal handles inline meal editing (PUT /meals/:id).
func (h *MealHandler) EditMeal(c *gin.Context) {
	mealID := c.Param("id")
	calories, _ := strconv.ParseFloat(c.PostForm("calories_kcal"), 64)
	protein, _ := strconv.ParseFloat(c.PostForm("protein_g"), 64)
	fat, _ := strconv.ParseFloat(c.PostForm("fat_g"), 64)
	carbs, _ := strconv.ParseFloat(c.PostForm("carbs_g"), 64)
	weightG, _ := strconv.ParseFloat(c.PostForm("weight_g"), 64)

	var eatenAt time.Time
	if raw := c.PostForm("eaten_at"); raw != "" {
		eatenAt, _ = time.Parse("2006-01-02T15:04", raw)
	}
	if eatenAt.IsZero() {
		eatenAt = time.Now()
	}

	meal, err := h.svc.EditMeal(c.Request.Context(), mealID, mealapplication.EditMealRequest{
		Name:         c.PostForm("name"),
		CaloriesKcal: calories,
		ProteinG:     protein,
		FatG:         fat,
		CarbsG:       carbs,
		WeightG:      weightG,
		EatenAt:      eatenAt,
	})
	if err != nil {
		slog.Error("edit meal failed", "meal_id", mealID, "error", err)
		renderTempl(c, partials.ErrorMsg(err.Error()))
		return
	}

	renderTempl(c, partials.MealRow(meal))
}

// ShowEditForm returns the inline edit form for a meal (GET /meals/:id/edit).
func (h *MealHandler) ShowEditForm(c *gin.Context) {
	mealID := c.Param("id")
	// We need to load the meal to pre-fill the form.
	// The meal service doesn't expose GetByID publicly yet — use EditMeal with same values.
	// For now, return the edit form with the meal ID; the form will POST back.
	renderTempl(c, partials.MealEditForm(mealID))
}

// DeleteMeal handles meal deletion (DELETE /meals/:id).
func (h *MealHandler) DeleteMeal(c *gin.Context) {
	mealID := c.Param("id")
	if err := h.svc.DeleteMeal(c.Request.Context(), mealID); err != nil {
		slog.Error("delete meal failed", "meal_id", mealID, "error", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	// Return empty — HTMX will remove the element via hx-swap="outerHTML".
	c.Status(http.StatusOK)
}
