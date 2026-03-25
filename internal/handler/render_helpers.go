package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/artmuc/fatyai/internal/views/pages"
	"github.com/artmuc/fatyai/internal/views/partials"
)

// RenderScanPage renders the scan page for the inline route closure in routes.go.
func RenderScanPage(c *gin.Context) {
	renderTempl(c, pages.ScanPage())
}

// RenderMealForm renders the manual meal entry form partial.
func RenderMealForm(c *gin.Context) {
	renderTempl(c, partials.MealForm())
}

// RenderActivityForm renders the activity entry form partial.
func RenderActivityForm(c *gin.Context) {
	renderTempl(c, partials.ActivityForm())
}
