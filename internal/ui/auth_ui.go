package ui

import (
	"encoding/base64"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tmsankram/gonotes/internal/auth"
	"github.com/tmsankram/gonotes/internal/users"
)

const (
	CSRF_COOKIE = "gonotes_csrf"
	JWT_COOKIE  = "gonotes_token"
)

type AuthUI struct {
	Users    *users.Service
	Renderer *Renderer
}

func NewAuthUI(us *users.Service, r *Renderer) *AuthUI {
	return &AuthUI{
		Users:    us,
		Renderer: r,
	}
}

func GenerateCSRF(c *gin.Context) string {
	token := uuid.New().String()
	enc := base64.StdEncoding.EncodeToString([]byte(token))
	c.SetCookie(CSRF_COOKIE, enc, 300, "/", "", false, true)
	return enc
}

func ValidateCSRF(c *gin.Context) error {
	cookie, err := c.Cookie(CSRF_COOKIE)
	if err != nil || cookie == "" {
		return errors.New("csrf cookie missing")
	}

	if err := c.Request.ParseForm(); err != nil {
		return errors.New("invalid form")
	}

	form := c.PostForm("csrf")
	if form == "" {
		return errors.New("csrf missing in form")
	}

	if form != cookie {
		return errors.New("csrf mismatch")
	}
	return nil
}

// GET /login
func (a *AuthUI) LoginPage(c *gin.Context) {
	csrf := GenerateCSRF(c)
	a.Renderer.Page(c, "auth/login.html", gin.H{
		"Title": "Login",
		"CSRF":  csrf,
	})
}

// POST /login (form)
func (a *AuthUI) LoginPost(c *gin.Context) {
	// Validate CSRF
	if err := ValidateCSRF(c); err != nil {
		a.Renderer.Page(c, "auth/login.html", gin.H{
			"Title": "Login",
			"CSRF":  GenerateCSRF(c),
			"Flash": "Invalid CSRF token",
		})
		return
	}

	email := c.PostForm("email")
	password := c.PostForm("password")

	// Basic Server side validation
	if email == "" || password == "" {
		a.Renderer.Page(c, "auth/login.html", gin.H{
			"Title": "Login",
			"CSRF":  GenerateCSRF(c),
			"Flash": "email and password are required",
		})
		return
	}

	u, err := a.Users.GetByEmail(email)
	if err != nil {
		a.Renderer.Page(c, "auth/login.html", gin.H{
			"Title": "Login",
			"CSRF":  GenerateCSRF(c),
			"Flash": "invalid credentials",
		})
		return
	}
	// If user has TOTP enabled, we should redirect them to TOTP verification.
	// For simplicity here: if totp enabled -> redirect to /totp/verify-ui (we can add later)
	if u.TOTPEnabled {
		// set a short-lived cookie to indicate stage (optional)
		// but for now, reject UI login and instruct user to use API flow that supports TOTP
		a.Renderer.Page(c, "auth/login.html", gin.H{
			"Title": "Login",
			"CSRF":  GenerateCSRF(c),
			"Flash": "TOTP enabled accounts must login via API (or implement additional UI step).",
		})
		return
	}
	// Generate JWT using existing internal/auth package
	token, err := auth.GenerateToken(u.ID)
	if err != nil {
		a.Renderer.Page(c, "auth/login.html", gin.H{
			"Title": "Login",
			"CSRF":  GenerateCSRF(c),
			"Flash": "internal error",
		})
		return
	}
	// Set secure, httpOnly cookie with JWT
	c.SetCookie(JWT_COOKIE, token, 24*3600, "/", "", false, true)

	// If this is an HTMX request, return a small fragment to redirect
	if c.GetHeader("HX-Request") == "true" {
		// HX-Redirect header causes htmx to navigate
		c.Header("HX-Redirect", "/notes")
		c.Status(http.StatusOK)
		return
	}

	// Non-HTMX: standard redirect
	c.Redirect(http.StatusFound, "/notes")
}

// GET /register
func (a *AuthUI) RegisterPage(c *gin.Context) {
	csrf := GenerateCSRF(c)
	a.Renderer.Page(c, "auth/register.html", gin.H{
		"Title": "Register",
		"CSRF":  csrf,
	})
}

// POST /register
func (a *AuthUI) RegisterPost(c *gin.Context) {
	if err := ValidateCSRF(c); err != nil {
		a.Renderer.Page(c, "auth/register.html", gin.H{
			"Title": "Register",
			"CSRF":  GenerateCSRF(c),
			"Flash": "Invalid CSRF token",
		})
		return
	}

	email := c.PostForm("email")
	password := c.PostForm("password")
	password2 := c.PostForm("password2")

	if email == "" || password == "" || password2 == "" {
		a.Renderer.Page(c, "auth/register.html", gin.H{
			"Title": "Register",
			"CSRF":  GenerateCSRF(c),
			"Flash": "All fields are required",
		})
		return
	}

	if password != password2 {
		a.Renderer.Page(c, "auth/register.html", gin.H{
			"Title": "Register",
			"CSRF":  GenerateCSRF(c),
			"Flash": "Passwords do not match",
		})
		return
	}

	// create user
	hash, err := users.HashPassword(password)
	if err != nil {
		a.Renderer.Page(c, "auth/register.html", gin.H{
			"Title": "Register",
			"CSRF":  GenerateCSRF(c),
			"Flash": "internal error",
		})
		return
	}

	u, err := a.Users.Create(users.User{
		Email:    email,
		Password: hash,
	})

	if err != nil {
		a.Renderer.Page(c, "auth/register.html", gin.H{
			"Title": "Register",
			"CSRF":  GenerateCSRF(c),
			"Flash": "could not create user (maybe email already exists)",
		})
		return
	}

	// on success, set cookie and redirect
	token, err := auth.GenerateToken(u.ID)
	if err != nil {
		a.Renderer.Page(c, "auth/register.html", gin.H{
			"Title": "Register",
			"CSRF":  GenerateCSRF(c),
			"Flash": "internal error",
		})
		return
	}

	c.SetCookie(JWT_COOKIE, token, 24*3600, "/", "", false, true)

	if c.GetHeader("HX-Request") == "true" {
		c.Header("HX-Redirect", "/notes")
		c.Status(http.StatusOK)
		return
	}

	c.Redirect(http.StatusFound, "/notes")
}
