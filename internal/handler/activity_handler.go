package handler

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	activityapplication "github.com/artmuc/fatyai/internal/application/activity"
	"github.com/artmuc/fatyai/internal/views/partials"
)

// ActivityHandler handles activity logging (F3).
type ActivityHandler struct {
	svc *activityapplication.Service
}

// NewActivityHandler creates a new ActivityHandler.
func NewActivityHandler(svc *activityapplication.Service) *ActivityHandler {
	return &ActivityHandler{svc: svc}
}

// LogActivity handles activity form submission (POST /activities).
func (h *ActivityHandler) LogActivity(c *gin.Context) {
	durationMin, _ := strconv.Atoi(c.PostForm("duration_min"))
	calories, _ := strconv.ParseFloat(c.PostForm("calories_burned"), 64)

	var loggedAt time.Time
	if raw := c.PostForm("logged_at"); raw != "" {
		loggedAt, _ = time.Parse("2006-01-02T15:04", raw)
	}
	if loggedAt.IsZero() {
		loggedAt = time.Now()
	}

	req := activityapplication.LogActivityRequest{
		UserID:         sessionUserID(c),
		ActivityType:   c.PostForm("activity_type"),
		DurationMin:    durationMin,
		Intensity:      c.PostForm("intensity"),
		CaloriesBurned: calories,
		LoggedAt:       loggedAt,
	}

	act, err := h.svc.LogActivity(c.Request.Context(), req)
	if err != nil {
		slog.Error("log activity failed", "error", err)
		renderTempl(c, partials.ErrorMsg(err.Error()))
		return
	}

	renderTempl(c, partials.ActivityRow(act))
}

// DeleteActivity handles activity deletion (DELETE /activities/:id).
func (h *ActivityHandler) DeleteActivity(c *gin.Context) {
	if err := h.svc.DeleteActivity(c.Request.Context(), c.Param("id")); err != nil {
		slog.Error("delete activity failed", "activity_id", c.Param("id"), "error", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}
