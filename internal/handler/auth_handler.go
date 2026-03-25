package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	userapplication "github.com/artmuc/fatyai/internal/application/user"
	domainuser "github.com/artmuc/fatyai/internal/domain/user"
	"github.com/artmuc/fatyai/internal/middleware"
	"github.com/artmuc/fatyai/internal/views/pages"
)


// AuthHandler handles login, registration, and logout.
type AuthHandler struct {
	svc *userapplication.Service
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(svc *userapplication.Service) *AuthHandler {
	return &AuthHandler{svc: svc}
}

// ShowLogin renders the login page (GET /login).
func (h *AuthHandler) ShowLogin(c *gin.Context) {
	renderTempl(c, pages.LoginPage(""))
}

// HandleLogin processes the login form (POST /login).
func (h *AuthHandler) HandleLogin(c *gin.Context) {
	var req userapplication.LoginRequest
	if err := c.ShouldBind(&req); err != nil {
		renderTempl(c, pages.LoginPage("Wypełnij wszystkie pola."))
		return
	}

	dto, err := h.svc.Login(c.Request.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, domainuser.ErrNotActive):
			renderTempl(c, pages.LoginPage("Potwierdź swój adres email przed logowaniem."))
		case errors.Is(err, domainuser.ErrInvalidCredentials):
			renderTempl(c, pages.LoginPage("Nieprawidłowy email lub hasło."))
		default:
			slog.Error("login failed", "error", err)
			renderTempl(c, pages.LoginPage("Błąd serwera. Spróbuj ponownie."))
		}
		return
	}

	session := sessions.Default(c)
	session.Set(middleware.SessionUserKey, dto.ID)
	_ = session.Save()

	c.Redirect(http.StatusSeeOther, "/journal")
}

// ShowRegister renders the registration page (GET /register).
func (h *AuthHandler) ShowRegister(c *gin.Context) {
	renderTempl(c, pages.RegisterPage(""))
}

// HandleRegister processes the registration form (POST /register).
func (h *AuthHandler) HandleRegister(c *gin.Context) {
	var req userapplication.RegisterRequest
	if err := c.ShouldBind(&req); err != nil {
		renderTempl(c, pages.RegisterPage("Wypełnij wszystkie pola (hasło min. 8 znaków)."))
		return
	}

	dto, err := h.svc.Register(c.Request.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, domainuser.ErrEmailTaken):
			renderTempl(c, pages.RegisterPage("Ten email jest już zajęty."))
		default:
			slog.Error("register failed", "error", err)
			renderTempl(c, pages.RegisterPage("Nie udało się wysłać maila weryfikacyjnego. Sprawdź adres email lub spróbuj ponownie."))
		}
		return
	}

	renderTempl(c, pages.CheckEmailPage(dto.Email))
}

// HandleVerify activates a user account via the token from the email link (GET /verify).
func (h *AuthHandler) HandleVerify(c *gin.Context) {
	token := c.Query("token")
	err := h.svc.VerifyEmail(c.Request.Context(), token)
	if err != nil {
		slog.Error("verify email failed", "error", err)
		renderTempl(c, pages.VerifyPage(false, "Link weryfikacyjny jest nieprawidłowy lub wygasł."))
		return
	}
	renderTempl(c, pages.VerifyPage(true, "Twoje konto zostało aktywowane. Możesz się teraz zalogować."))
}

// Logout clears the session and redirects to /login (POST /logout).
func (h *AuthHandler) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	_ = session.Save()
	c.Redirect(http.StatusSeeOther, "/login")
}
