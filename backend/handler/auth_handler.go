package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/onoderaryou/smart-store-admin/backend/config"
	"github.com/onoderaryou/smart-store-admin/backend/models"
	"github.com/onoderaryou/smart-store-admin/backend/repository"
	"github.com/onoderaryou/smart-store-admin/backend/utils/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthHandler struct {
	userRepo     repository.UserRepository
	authConfig   *config.AuthConfig
	oauth2Config *oauth2.Config
}

func NewAuthHandler(userRepo repository.UserRepository, authConfig *config.AuthConfig) *AuthHandler {
	oauth2Config := &oauth2.Config{
		ClientID:     authConfig.GoogleClientID,
		ClientSecret: authConfig.GoogleClientSecret,
		RedirectURL:  "http://localhost:3000/api/auth/callback/google", // TODO: 環境変数化
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return &AuthHandler{
		userRepo:     userRepo,
		authConfig:   authConfig,
		oauth2Config: oauth2Config,
	}
}

type GoogleUserInfo struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

func (h *AuthHandler) HandleGoogleCallback(c echo.Context) error {
	code := c.QueryParam("code")
	if code == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "missing code parameter")
	}

	oauthToken, err := h.oauth2Config.Exchange(context.Background(), code)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to exchange token")
	}

	userInfo, err := h.getGoogleUserInfo(oauthToken.AccessToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get user info")
	}

	user, err := h.userRepo.FindByGoogleID(context.Background(), userInfo.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to find user")
	}

	if user == nil {
		// 新規ユーザー作成
		user = &models.User{
			Email:    userInfo.Email,
			Name:     userInfo.Name,
			Picture:  userInfo.Picture,
			GoogleID: userInfo.ID,
			Role:     models.RoleViewer, // デフォルトロール
		}
		if err := h.userRepo.Create(context.Background(), user); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to create user")
		}
	} else {
		// 既存ユーザー更新
		user.LastLoginAt = time.Now()
		if err := h.userRepo.Update(context.Background(), user); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to update user")
		}
	}

	// JWTトークン生成
	jwtToken, err := h.generateAuthToken(user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to generate token")
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": jwtToken,
		"user": models.UserResponse{
			ID:      user.ID,
			Email:   user.Email,
			Name:    user.Name,
			Picture: user.Picture,
			Role:    user.Role,
		},
	})
}

func (h *AuthHandler) getGoogleUserInfo(accessToken string) (*GoogleUserInfo, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + accessToken)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userInfo GoogleUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}

func (h *AuthHandler) generateAuthToken(user *models.User) (string, error) {
	return jwt.GenerateToken(user.ID, string(user.Role), h.authConfig.JWTSecret)
}

// GetGoogleAuthURL Googleログイン用のURLを生成
func (h *AuthHandler) GetGoogleAuthURL() string {
	return h.oauth2Config.AuthCodeURL("state")
}

// GetCurrentUser 現在のユーザー情報を取得
func (h *AuthHandler) GetCurrentUser(c echo.Context) error {
	userID := c.Get("user_id").(primitive.ObjectID)

	user, err := h.userRepo.FindByID(context.Background(), userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get user")
	}
	if user == nil {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}

	return c.JSON(http.StatusOK, models.UserResponse{
		ID:      user.ID,
		Email:   user.Email,
		Name:    user.Name,
		Picture: user.Picture,
		Role:    user.Role,
	})
}
