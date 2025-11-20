package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tmsankram/gonotes/internal/config"
	"github.com/tmsankram/gonotes/internal/response"
	"github.com/tmsankram/gonotes/internal/users"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

type OAuthHandler struct {
	users  *users.Service
	cfg    *config.Config
	google *oauth2.Config
	github *oauth2.Config
}

func NewOauthHandler(usersSvc *users.Service, cfg *config.Config) *OAuthHandler {
	return &OAuthHandler{
		users: usersSvc,
		cfg:   cfg,
		google: &oauth2.Config{
			ClientID:     cfg.GoogleClientID,
			ClientSecret: cfg.GithubClientSecret,
			RedirectURL:  cfg.GoogleRedirectURL,
			Scopes:       []string{"email", "profile"},
			Endpoint:     google.Endpoint,
		},
		github: &oauth2.Config{
			ClientID:     cfg.GithubClientID,
			ClientSecret: cfg.GithubClientSecret,
			RedirectURL:  cfg.GithubRedirectURL,
			Scopes:       []string{"user:email"},
			Endpoint:     github.Endpoint,
		},
	}
}

func randomState() string {
	rand.Seed(time.Now().UnixNano())
	return RandString(32)
}

func RandString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func (h *OAuthHandler) GoogleLogin(c *gin.Context) {
	state := randomState()
	c.SetCookie("oauth_state", state, 300, "/", "", false, true)

	url := h.google.AuthCodeURL(state)
	c.Redirect(302, url)
}

func (h *OAuthHandler) GoogleCallback(c *gin.Context) {
	stateCookie, _ := c.Cookie("oauth_state")
	stateQuery := c.Query("state")

	if stateCookie == "" || stateCookie != stateQuery {
		response.Unauthorized(c, errors.New("invalid oauth state"))
		return
	}

	code := c.Query("code")
	token, err := h.google.Exchange(context.Background(), code)
	if err != nil {
		response.Unauthorized(c, err)
		return
	}

	client := h.google.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		response.Internal(c, err)
		return
	}
	defer resp.Body.Close()

	var data struct {
		Email string `json:"email"`
		Id    string `json:"id"`
	}
	json.NewDecoder(resp.Body).Decode(&data)

	u, err := h.users.GetByOauth("google", data.Id)

	if err != nil {
		// If not found, create a new one
		u, err = h.users.CreateOAuthUser(data.Email, "google", data.Id)
		if err != nil {
			response.Internal(c, err)
		}
	}
	jwt, err := GenerateToken(u.ID)
	if err != nil {
		response.Internal(c, err)
		return
	}
	response.Success(c, "oauth login successful", gin.H{"token": jwt})
}

func (h *OAuthHandler) GithubLogin(c *gin.Context) {
	state := randomState()
	c.SetCookie("oauth_state", state, 300, "/", "", false, true)

	url := h.github.AuthCodeURL(state)
	c.Redirect(http.StatusFound, url)
}

func (h *OAuthHandler) GithubCallback(c *gin.Context) {
	stateCookie, _ := c.Cookie("oauth_state")
	if stateCookie == "" || stateCookie != c.Query("state") {
		response.Unauthorized(c, errors.New("invalid oauth state"))
		return
	}

	code := c.Query("code")
	token, err := h.github.Exchange(context.Background(), code)
	if err != nil {
		response.Unauthorized(c, err)
		return
	}

	client := h.github.Client(context.Background(), token)

	// STEP 1: Fetch basic user info
	userResp, err := client.Get("https://api.github.com/user")
	if err != nil {
		response.Internal(c, err)
		return
	}
	defer userResp.Body.Close()

	var ghUser struct {
		ID    int64  `json:"id"`
		Email string `json:"email"`
	}

	if err := json.NewDecoder(userResp.Body).Decode(&ghUser); err != nil {
		response.Internal(c, err)
		return
	}

	// STEP 2: If GitHub didn’t send email, fetch emails list
	if ghUser.Email == "" {
		email, err := fetchGithubPrimaryEmail(client)
		if err != nil {
			response.Internal(c, err)
			return
		}
		ghUser.Email = email
	}

	// STEP 3: Login / register user in DB
	u, err := h.users.GetByOauth("github", fmt.Sprint(ghUser.ID))
	if err != nil {
		// create new OAuth user
		u, err = h.users.CreateOAuthUser(ghUser.Email, "github", fmt.Sprint(ghUser.ID))
		if err != nil {
			response.Internal(c, err)
			return
		}
	}

	// STEP 4: Issue JWT
	jwt, err := GenerateToken(u.ID)
	if err != nil {
		response.Internal(c, err)
		return
	}

	response.Success(c, "oauth login successful", gin.H{
		"token": jwt,
	})
}

func fetchGithubPrimaryEmail(client *http.Client) (string, error) {
	resp, err := client.Get("https://api.github.com/user/emails")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var emails []struct {
		Email    string `json:"email"`
		Primary  bool   `json:"primary"`
		Verified bool   `json:"verified"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		return "", err
	}

	for _, e := range emails {
		if e.Primary && e.Verified {
			return e.Email, nil
		}
	}

	if len(emails) > 0 {
		return emails[0].Email, nil
	}

	return "", errors.New("no email found")
}
