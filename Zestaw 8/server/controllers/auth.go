package controllers

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"zadanie8/database"
	"zadanie8/models"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/gorm"
)

const sessionTTL = 7 * 24 * time.Hour

type authRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authResponse struct {
	Token string           `json:"token"`
	User  authUserResponse `json:"user"`
}

type authUserResponse struct {
	ID        uint     `json:"id"`
	Email     string   `json:"email"`
	Name      string   `json:"name"`
	Providers []string `json:"providers"`
}

type oauthProfile struct {
	ProviderUserID string
	Email          string
	Name           string
}

func Register(c echo.Context) error {
	var input authRequest
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	email := normalizeEmail(input.Email)
	if email == "" || len(input.Password) < 8 {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "email and password with at least 8 characters are required"})
	}

	var existing models.User
	if err := database.DB.Where("email = ?", email).First(&existing).Error; err == nil {
		return c.JSON(http.StatusConflict, echo.Map{"error": "user already exists"})
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "database error"})
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "password hashing failed"})
	}

	user := models.User{
		Email:        email,
		Name:         strings.TrimSpace(input.Name),
		PasswordHash: string(passwordHash),
	}
	if err := database.DB.Create(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "could not create user"})
	}

	return issueAuthResponse(c, user, http.StatusCreated)
}

func Login(c echo.Context) error {
	var input authRequest
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	var user models.User
	if err := database.DB.Preload("OAuthAccounts").Where("email = ?", normalizeEmail(input.Email)).First(&user).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid email or password"})
	}

	if user.PasswordHash == "" || bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)) != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid email or password"})
	}

	return issueAuthResponse(c, user, http.StatusOK)
}

func Logout(c echo.Context) error {
	session, err := sessionFromBearer(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": err.Error()})
	}

	if err := database.DB.Delete(session).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "logout failed"})
	}

	return c.NoContent(http.StatusNoContent)
}

func Me(c echo.Context) error {
	user, _, err := userFromBearer(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, userResponse(user))
}

func GoogleLogin(c echo.Context) error {
	return beginOAuth(c, "google", googleOAuthConfig())
}

func GoogleCallback(c echo.Context) error {
	return finishOAuth(c, "google", googleOAuthConfig(), fetchGoogleProfile)
}

func GitHubLogin(c echo.Context) error {
	return beginOAuth(c, "github", githubOAuthConfig())
}

func GitHubCallback(c echo.Context) error {
	return finishOAuth(c, "github", githubOAuthConfig(), fetchGitHubProfile)
}

func issueAuthResponse(c echo.Context, user models.User, status int) error {
	token, err := createSession(user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "could not create session"})
	}

	if err := database.DB.Preload("OAuthAccounts").First(&user, user.ID).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "could not load user"})
	}

	return c.JSON(status, authResponse{
		Token: token,
		User:  userResponse(&user),
	})
}

func beginOAuth(c echo.Context, provider string, cfg *oauth2.Config) error {
	if cfg.ClientID == "" || cfg.ClientSecret == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": provider + " login is not configured"})
	}

	state, err := randomToken(32)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "state generation failed"})
	}

	record := models.OAuthState{
		Provider:  provider,
		StateHash: hashToken(state),
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}
	if err := database.DB.Create(&record).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "could not store oauth state"})
	}

	authURL := cfg.AuthCodeURL(state, oauth2.AccessTypeOffline)
	return c.Redirect(http.StatusTemporaryRedirect, authURL)
}

func finishOAuth(
	c echo.Context,
	provider string,
	cfg *oauth2.Config,
	profileFetcher func(context.Context, *oauth2.Config, *oauth2.Token) (oauthProfile, error),
) error {
	if errValue := c.QueryParam("error"); errValue != "" {
		return redirectWithAuthError(c, errValue)
	}

	state := c.QueryParam("state")
	code := c.QueryParam("code")
	if state == "" || code == "" {
		return redirectWithAuthError(c, "missing oauth callback parameters")
	}

	if err := consumeOAuthState(provider, state); err != nil {
		return redirectWithAuthError(c, err.Error())
	}

	token, err := cfg.Exchange(context.Background(), code)
	if err != nil {
		return redirectWithAuthError(c, "token exchange failed")
	}

	profile, err := profileFetcher(context.Background(), cfg, token)
	if err != nil {
		return redirectWithAuthError(c, "profile fetch failed")
	}

	user, err := upsertOAuthUser(provider, profile, token)
	if err != nil {
		return redirectWithAuthError(c, "could not persist oauth user")
	}

	appToken, err := createSession(user.ID)
	if err != nil {
		return redirectWithAuthError(c, "could not create application session")
	}

	callbackURL := frontendCallbackURL()
	params := callbackURL.Query()
	params.Set("token", appToken)
	callbackURL.RawQuery = params.Encode()

	return c.Redirect(http.StatusTemporaryRedirect, callbackURL.String())
}

func fetchGoogleProfile(ctx context.Context, cfg *oauth2.Config, token *oauth2.Token) (oauthProfile, error) {
	client := cfg.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return oauthProfile{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return oauthProfile{}, fmt.Errorf("google profile status %d", resp.StatusCode)
	}

	var body struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return oauthProfile{}, err
	}

	return oauthProfile{
		ProviderUserID: body.ID,
		Email:          normalizeEmail(body.Email),
		Name:           body.Name,
	}, nil
}

func fetchGitHubProfile(ctx context.Context, cfg *oauth2.Config, token *oauth2.Token) (oauthProfile, error) {
	client := cfg.Client(ctx, token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return oauthProfile{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return oauthProfile{}, fmt.Errorf("github profile status %d", resp.StatusCode)
	}

	var body struct {
		ID    int64  `json:"id"`
		Login string `json:"login"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return oauthProfile{}, err
	}

	email := normalizeEmail(body.Email)
	if email == "" {
		email = fetchPrimaryGitHubEmail(client)
	}

	name := strings.TrimSpace(body.Name)
	if name == "" {
		name = body.Login
	}

	return oauthProfile{
		ProviderUserID: fmt.Sprintf("%d", body.ID),
		Email:          email,
		Name:           name,
	}, nil
}

func fetchPrimaryGitHubEmail(client *http.Client) string {
	resp, err := client.Get("https://api.github.com/user/emails")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return ""
	}

	var emails []struct {
		Email    string `json:"email"`
		Primary  bool   `json:"primary"`
		Verified bool   `json:"verified"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		return ""
	}

	for _, item := range emails {
		if item.Primary && item.Verified {
			return normalizeEmail(item.Email)
		}
	}
	return ""
}

func upsertOAuthUser(provider string, profile oauthProfile, token *oauth2.Token) (*models.User, error) {
	if profile.ProviderUserID == "" {
		return nil, errors.New("missing provider user id")
	}

	var account models.OAuthAccount
	err := database.DB.Where("provider = ? AND provider_user_id = ?", provider, profile.ProviderUserID).First(&account).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	var user models.User
	if errors.Is(err, gorm.ErrRecordNotFound) {
		email := profile.Email
		if email != "" {
			if lookupErr := database.DB.Where("email = ?", email).First(&user).Error; lookupErr != nil && !errors.Is(lookupErr, gorm.ErrRecordNotFound) {
				return nil, lookupErr
			}
		}

		if user.ID == 0 {
			if email == "" {
				email = fmt.Sprintf("%s-%s@oauth.local", provider, profile.ProviderUserID)
			}
			user = models.User{
				Email: email,
				Name:  profile.Name,
			}
			if err := database.DB.Create(&user).Error; err != nil {
				return nil, err
			}
		}

		account = models.OAuthAccount{
			UserID:         user.ID,
			Provider:       provider,
			ProviderUserID: profile.ProviderUserID,
		}
	} else if err := database.DB.First(&user, account.UserID).Error; err != nil {
		return nil, err
	}

	account.Email = profile.Email
	account.Name = profile.Name
	account.AccessToken = token.AccessToken
	account.RefreshToken = token.RefreshToken
	account.TokenType = token.TokenType
	account.Expiry = token.Expiry
	if err := database.DB.Save(&account).Error; err != nil {
		return nil, err
	}

	updates := map[string]any{}
	if user.Name == "" && profile.Name != "" {
		updates["name"] = profile.Name
	}
	if user.Email == "" && profile.Email != "" {
		updates["email"] = profile.Email
	}
	if len(updates) > 0 {
		if err := database.DB.Model(&user).Updates(updates).Error; err != nil {
			return nil, err
		}
	}

	if err := database.DB.Preload("OAuthAccounts").First(&user, user.ID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func createSession(userID uint) (string, error) {
	token, err := randomToken(32)
	if err != nil {
		return "", err
	}

	session := models.AuthSession{
		UserID:    userID,
		TokenHash: hashToken(token),
		ExpiresAt: time.Now().Add(sessionTTL),
	}
	if err := database.DB.Create(&session).Error; err != nil {
		return "", err
	}

	return token, nil
}

func userFromBearer(c echo.Context) (*models.User, *models.AuthSession, error) {
	session, err := sessionFromBearer(c)
	if err != nil {
		return nil, nil, err
	}

	var user models.User
	if err := database.DB.Preload("OAuthAccounts").First(&user, session.UserID).Error; err != nil {
		return nil, nil, errors.New("user not found")
	}

	return &user, session, nil
}

func sessionFromBearer(c echo.Context) (*models.AuthSession, error) {
	authHeader := c.Request().Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return nil, errors.New("missing bearer token")
	}

	token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	if token == "" {
		return nil, errors.New("missing bearer token")
	}

	var session models.AuthSession
	if err := database.DB.Where("token_hash = ?", hashToken(token)).First(&session).Error; err != nil {
		return nil, errors.New("invalid bearer token")
	}

	if time.Now().After(session.ExpiresAt) {
		database.DB.Delete(&session)
		return nil, errors.New("expired bearer token")
	}

	return &session, nil
}

func consumeOAuthState(provider string, state string) error {
	var record models.OAuthState
	err := database.DB.Where("provider = ? AND state_hash = ?", provider, hashToken(state)).First(&record).Error
	if err != nil {
		return errors.New("invalid oauth state")
	}

	database.DB.Delete(&record)
	if time.Now().After(record.ExpiresAt) {
		return errors.New("expired oauth state")
	}

	return nil
}

func redirectWithAuthError(c echo.Context, message string) error {
	callbackURL := frontendCallbackURL()
	params := callbackURL.Query()
	params.Set("error", message)
	callbackURL.RawQuery = params.Encode()
	return c.Redirect(http.StatusTemporaryRedirect, callbackURL.String())
}

func userResponse(user *models.User) authUserResponse {
	providers := make([]string, 0, len(user.OAuthAccounts))
	for _, account := range user.OAuthAccounts {
		providers = append(providers, account.Provider)
	}

	return authUserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Providers: providers,
	}
}

func googleOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  serverBaseURL() + "/auth/google/callback",
		Scopes:       []string{"openid", "email", "profile"},
		Endpoint:     google.Endpoint,
	}
}

func githubOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		RedirectURL:  serverBaseURL() + "/auth/github/callback",
		Scopes:       []string{"read:user", "user:email"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://github.com/login/oauth/authorize",
			TokenURL: "https://github.com/login/oauth/access_token",
		},
	}
}

func frontendCallbackURL() url.URL {
	base := strings.TrimRight(env("CLIENT_BASE_URL", "http://localhost:5173"), "/")
	callbackURL, err := url.Parse(base + "/auth/callback")
	if err != nil {
		return url.URL{Scheme: "http", Host: "localhost:5173", Path: "/auth/callback"}
	}
	return *callbackURL
}

func serverBaseURL() string {
	return strings.TrimRight(env("SERVER_BASE_URL", "http://localhost:8080"), "/")
}

func env(key string, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func randomToken(size int) (string, error) {
	bytes := make([]byte, size)
	if _, err := io.ReadFull(rand.Reader, bytes); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
