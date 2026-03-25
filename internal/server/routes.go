package router

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	activityapplication "github.com/artmuc/fatyai/internal/application/activity"
	mealapplication "github.com/artmuc/fatyai/internal/application/meal"
	profileapplication "github.com/artmuc/fatyai/internal/application/profile"
	weightapplication "github.com/artmuc/fatyai/internal/application/weight"
	userapplication "github.com/artmuc/fatyai/internal/application/user"
	"github.com/artmuc/fatyai/internal/eventbus"
	"github.com/artmuc/fatyai/internal/handler"
	"github.com/artmuc/fatyai/internal/middleware"
	activitypersistence "github.com/artmuc/fatyai/internal/infrastructure/persistence/activity"
	mealpersistence "github.com/artmuc/fatyai/internal/infrastructure/persistence/meal"
	profilepersistence "github.com/artmuc/fatyai/internal/infrastructure/persistence/profile"
	weightpersistence "github.com/artmuc/fatyai/internal/infrastructure/persistence/weight"
	userpersistence "github.com/artmuc/fatyai/internal/infrastructure/persistence/user"
	emailinfra "github.com/artmuc/fatyai/internal/infrastructure/email"
	openaiinfra "github.com/artmuc/fatyai/internal/infrastructure/openai"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.GET("/health", s.healthHandler)

	gormDB := s.db.GetDB()
	cfg := s.cfg

	// Static files and PWA assets.
	r.Static("/static", "./static")
	r.StaticFile("/manifest.json", "./static/manifest.json")
	r.StaticFile("/sw.js", "./static/sw.js")

	// Session middleware (cookie store — signed + encrypted).
	store := cookie.NewStore([]byte(cfg.SessionSecret))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 30, // 30 days
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
	r.Use(sessions.Sessions("kalorie_session", store))

	// Shared event bus.
	bus := eventbus.New()
	bus.RegisterDefaultHandlers()

	// ---------------------------------------------------------------
	// Repositories & services
	// ---------------------------------------------------------------
	userRepo := userpersistence.NewUserRepository(gormDB)
	emailClient := emailinfra.NewBrevoClient(cfg.BrevoAPIKey, cfg.BrevoSenderEmail, cfg.BrevoSenderName)
	userSvc := userapplication.NewService(userRepo, bus, emailClient, cfg.AppBaseURL)

	// ---------------------------------------------------------------
	// Auth routes (public)
	// ---------------------------------------------------------------
	authH := handler.NewAuthHandler(userSvc)
	r.GET("/login", authH.ShowLogin)
	r.POST("/login", authH.HandleLogin)
	r.GET("/register", authH.ShowRegister)
	r.POST("/register", authH.HandleRegister)
	r.GET("/verify", authH.HandleVerify)
	r.POST("/logout", authH.Logout)

	// ---------------------------------------------------------------
	// Kalorie AI: repositories & services
	// ---------------------------------------------------------------
	profileRepo := profilepersistence.NewProfileRepository(gormDB)
	mealRepo := mealpersistence.NewMealRepository(gormDB)
	activityRepo := activitypersistence.NewActivityRepository(gormDB)
	weightRepo := weightpersistence.NewWeightRepository(gormDB)

	visionClient := openaiinfra.NewClient(cfg.OpenAIKey, cfg.OpenAIModel)

	profileSvc := profileapplication.NewService(profileRepo, bus)
	mealSvc := mealapplication.NewService(mealRepo, activityRepo, profileRepo, visionClient, bus)
	activitySvc := activityapplication.NewService(activityRepo, bus)
	weightSvc := weightapplication.NewService(weightRepo, profileRepo, bus)

	// ---------------------------------------------------------------
	// Kalorie AI: handlers
	// ---------------------------------------------------------------
	profileH := handler.NewProfileHandler(profileSvc)
	journalH := handler.NewJournalHandler(mealSvc, profileSvc)
	scanH := handler.NewScanHandler(mealSvc)
	mealH := handler.NewMealHandler(mealSvc)
	activityH := handler.NewActivityHandler(activitySvc)
	weightH := handler.NewWeightHandler(weightSvc)
	chartsH := handler.NewChartsHandler(mealSvc, activitySvc)

	// ---------------------------------------------------------------
	// Kalorie AI: protected UI routes
	// ---------------------------------------------------------------
	protected := r.Group("/")
	protected.Use(middleware.RequireAuth())

	// Root redirect.
	protected.GET("/", func(c *gin.Context) { c.Redirect(http.StatusSeeOther, "/journal") })

	// Onboarding (F4)
	protected.GET("/onboarding", profileH.ShowOnboarding)
	protected.POST("/onboarding", profileH.SubmitOnboarding)

	// Journal (F2)
	protected.GET("/journal", journalH.ShowJournal)
	protected.GET("/journal/partials/summary", journalH.JournalSummaryPartial)

	// Scan (F1)
	protected.GET("/scan", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		handler.RenderScanPage(c)
	})

	// Meals (F1, F2)
	protected.POST("/meals/scan", scanH.Scan)
	protected.POST("/meals", mealH.LogMeal)
	protected.PUT("/meals/:id", mealH.EditMeal)
	protected.DELETE("/meals/:id", mealH.DeleteMeal)
	protected.GET("/meals/:id/edit", mealH.ShowEditForm)
	protected.GET("/meals/manual-form", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		handler.RenderMealForm(c)
	})

	// Activities (F3)
	protected.POST("/activities", activityH.LogActivity)
	protected.DELETE("/activities/:id", activityH.DeleteActivity)
	protected.GET("/activities/form", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		handler.RenderActivityForm(c)
	})

	// Weight & correction (F5, F6)
	protected.POST("/weight", weightH.LogWeight)
	protected.GET("/correction/check", weightH.CheckCorrection)
	protected.POST("/correction/apply", weightH.ApplyCorrection)

	// Charts (F6)
	protected.GET("/charts", chartsH.ShowCharts)
	protected.GET("/charts/partials/weight", weightH.GetWeightChartData)
	protected.GET("/charts/partials/calories", chartsH.GetCalorieChartData)
	protected.GET("/charts/partials/net-calories", chartsH.GetNetCalorieChartData)
	protected.GET("/charts/partials/macros", chartsH.GetMacroChartData)

	// Profile (F4)
	protected.GET("/profile", profileH.ShowProfile)
	protected.POST("/profile", profileH.UpdateProfile)

	return r
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"
	c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}
