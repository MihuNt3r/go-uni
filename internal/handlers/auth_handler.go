package handlers

import (
	"net/http"
	"time"

	"go-uni/pkg/middleware"
)

type AuthHandler struct {
	username  string
	password  string
	jwtSecret string
	tokenTTL  time.Duration
}

type tokenRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type tokenResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	TokenType string    `json:"token_type"`
}

func NewAuthHandler(username, password, jwtSecret string, tokenTTL time.Duration) *AuthHandler {
	return &AuthHandler{
		username:  username,
		password:  password,
		jwtSecret: jwtSecret,
		tokenTTL:  tokenTTL,
	}
}

// Token godoc
// @Summary Issue auth token
// @Description Issues JWT token for authenticated API access.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body tokenRequest true "Credentials"
// @Success 200 {object} tokenResponse
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /auth/token [post]
func (h *AuthHandler) Token(w http.ResponseWriter, r *http.Request) {
	var req tokenRequest
	if err := readJSON(w, r, &req); err != nil {
		middleware.LogHandlerError(r, "invalid JSON body", err)
		_ = writeJSONError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	if validationErr := validatePayload(req); validationErr != nil {
		middleware.LogHandlerError(r, "token request validation failed", validationErr)
		_ = writeJSONError(w, http.StatusBadRequest, validationErr.Error())
		return
	}

	if req.Username != h.username || req.Password != h.password {
		middleware.LogHandlerError(r, "invalid credentials", nil)
		_ = writeJSONError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	token, err := middleware.GenerateJWT(h.jwtSecret, req.Username, h.tokenTTL)
	if err != nil {
		middleware.LogHandlerError(r, "failed to issue token", err)
		_ = writeJSONError(w, http.StatusInternalServerError, "failed to issue token")
		return
	}

	resp := tokenResponse{
		Token:     token,
		ExpiresAt: time.Now().Add(h.tokenTTL),
		TokenType: "Bearer",
	}

	jsonResponse(w, http.StatusOK, resp)
}
