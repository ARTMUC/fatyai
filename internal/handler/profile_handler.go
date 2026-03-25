package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	profileapplication "github.com/artmuc/fatyai/internal/application/profile"
	domainprofile "github.com/artmuc/fatyai/internal/domain/profile"
	"github.com/artmuc/fatyai/internal/views/pages"
	"github.com/artmuc/fatyai/internal/views/partials"
)

// ProfileHandler handles onboarding and profile management routes.
type ProfileHandler struct {
	svc *profileapplication.Service
}

// NewProfileHandler creates a new ProfileHandler.
func NewProfileHandler(svc *profileapplication.Service) *ProfileHandler {
	return &ProfileHandler{svc: svc}
}

// ShowOnboarding renders the onboarding page (GET /onboarding).
func (h *ProfileHandler) ShowOnboarding(c *gin.Context) {
	renderTempl(c, pages.OnboardingPage())
}

// SubmitOnboarding handles the onboarding form (POST /onboarding).
func (h *ProfileHandler) SubmitOnboarding(c *gin.Context) {
	birthYear, _ := strconv.Atoi(c.PostForm("birth_year"))
	heightCm, _ := strconv.ParseFloat(c.PostForm("height_cm"), 64)
	weightKg, _ := strconv.ParseFloat(c.PostForm("weight_kg"), 64)
	goalKgPerWeek, _ := strconv.ParseFloat(c.PostForm("goal_kg_per_week"), 64)

	req := profileapplication.OnboardRequest{
		UserID:        sessionUserID(c),
		Gender:        c.PostForm("gender"),
		BirthYear:     birthYear,
		HeightCm:      heightCm,
		WeightKg:      weightKg,
		ActivityLevel: c.PostForm("activity_level"),
		GoalKgPerWeek: goalKgPerWeek,
	}

	_, err := h.svc.Onboard(c.Request.Context(), req)
	if err != nil {
		slog.Error("onboard failed", "error", err)
		if isHTMX(c) {
			renderTempl(c, partials.ErrorMsg(err.Error()))
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if isHTMX(c) {
		c.Header("HX-Redirect", "/journal")
		c.Status(http.StatusOK)
		return
	}
	c.Redirect(http.StatusSeeOther, "/journal")
}

// ShowProfile renders the profile settings page (GET /profile).
func (h *ProfileHandler) ShowProfile(c *gin.Context) {
	p, err := h.svc.GetByUserID(c.Request.Context(), sessionUserID(c))
	if err != nil {
		if errors.Is(err, domainprofile.ErrNotFound) {
			c.Redirect(http.StatusSeeOther, "/onboarding")
			return
		}
		slog.Error("get profile failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	renderTempl(c, pages.ProfilePage(p))
}

// UpdateProfile handles profile update form submission (POST /profile).
func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	heightCm, _ := strconv.ParseFloat(c.PostForm("height_cm"), 64)
	weightKg, _ := strconv.ParseFloat(c.PostForm("weight_kg"), 64)
	goalKgPerWeek, _ := strconv.ParseFloat(c.PostForm("goal_kg_per_week"), 64)
	birthYear, _ := strconv.Atoi(c.PostForm("birth_year"))

	req := profileapplication.OnboardRequest{
		UserID:        sessionUserID(c),
		Gender:        c.PostForm("gender"),
		BirthYear:     birthYear,
		HeightCm:      heightCm,
		WeightKg:      weightKg,
		ActivityLevel: c.PostForm("activity_level"),
		GoalKgPerWeek: goalKgPerWeek,
	}

	p, err := h.svc.Onboard(c.Request.Context(), req)
	if err != nil {
		slog.Error("update profile failed", "error", err)
		renderTempl(c, partials.ErrorMsg(err.Error()))
		return
	}
	renderTempl(c, pages.ProfilePage(p))
}
