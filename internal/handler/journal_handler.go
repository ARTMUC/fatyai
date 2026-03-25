package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	mealapplication "github.com/artmuc/fatyai/internal/application/meal"
	profileapplication "github.com/artmuc/fatyai/internal/application/profile"
	domainprofile "github.com/artmuc/fatyai/internal/domain/profile"
	"github.com/artmuc/fatyai/internal/views/pages"
	"github.com/artmuc/fatyai/internal/views/partials"
)

// JournalHandler renders the daily journal view (F2).
type JournalHandler struct {
	mealSvc    *mealapplication.Service
	profileSvc *profileapplication.Service
}

// NewJournalHandler creates a new JournalHandler.
func NewJournalHandler(mealSvc *mealapplication.Service, profileSvc *profileapplication.Service) *JournalHandler {
	return &JournalHandler{mealSvc: mealSvc, profileSvc: profileSvc}
}

// ShowJournal renders the full journal page (GET /journal).
func (h *JournalHandler) ShowJournal(c *gin.Context) {
	date := dateFromQuery(c)
	journal, err := h.mealSvc.GetDailyJournal(c.Request.Context(), sessionUserID(c), date)
	if err != nil {
		slog.Error("get daily journal failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_, profileErr := h.profileSvc.GetByUserID(c.Request.Context(), sessionUserID(c))
	if profileErr != nil {
		if errors.Is(profileErr, domainprofile.ErrNotFound) {
			c.Redirect(http.StatusSeeOther, "/onboarding")
			return
		}
		slog.Error("get profile for journal failed", "error", profileErr)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	renderTempl(c, pages.JournalPage(journal, date))
}

// JournalSummaryPartial returns only the summary bar (HTMX swap target).
func (h *JournalHandler) JournalSummaryPartial(c *gin.Context) {
	date := dateFromQuery(c)
	journal, err := h.mealSvc.GetDailyJournal(c.Request.Context(), sessionUserID(c), date)
	if err != nil {
		slog.Error("get journal summary failed", "error", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	renderTempl(c, partials.JournalSummary(journal))
}
