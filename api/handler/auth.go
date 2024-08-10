package handler

import (
	t "auth/api/handler/token"
	pb "auth/generated/auth"
	"auth/model"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// @Summary Register user
// @Description Create new user
// @Tags Auth
// @Param info body auth.RegistrRequest true "User info"
// @Success 200 {object} auth.RegistrResponse
// @Failure 400 {object} string "Invalid data"
// @Failure 500 {object} string "Server error"
// @Router /auth/register [post]
func (h *Handler) RegisterHandler(ctx *gin.Context) {
	var req pb.RegistrRequest
	if err := ctx.ShouldBind(&req); err != nil {
		fmt.Println("saldjfaoidsjfoiasdfasd")
		h.Log.Error("Error binding request", "err", err)
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	if err != nil {
		h.Log.Error("Error hashing password", "err", err)
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	req.Password = string(passwordHash)

	resp, err := h.User.Registr(ctx, &req)
	if err != nil {
		h.Log.Error("Error registering", "err", err)
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	h.Log.Info("Registered")
	ctx.JSON(200, resp)
}

// @Summary Login user
// @Discription it generates new acces and refresh tokens
// @Tags Auth
// @Param LoginRequest body auth.LoginRequest true "email and password"
// @Success 200 {object} string
// @Failure 400 {object} map[string]interface{} "Failed request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/login [post]
func (h *Handler) LoginHandler(ctx *gin.Context) {
	var req pb.LoginRequest
	if err := ctx.ShouldBind(&req); err != nil {
		h.Log.Error("Error binding request", "err", err)
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.User.Login(ctx, &req)
	if err != nil {
		fmt.Println("sdfsmaksdk")
		h.Log.Error("Error logging in", "err", err)
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(resp.Password), []byte(req.Password))

	if err != nil {
		h.Log.Error("Invalid credentials")
		ctx.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	claims := pb.UserClaims{
		Id:          resp.Id,
		Email:       resp.Email,
		PhoneNumber: resp.PhoneNumber,
		FullName:    resp.FullName,
		Role:        resp.Role,
	}

	tokens, err := t.GenerateJwt(&claims)
	if err != nil {
		fmt.Println(err)
		h.Log.Error("Error generating JWT tokens", "err", err)
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	err = h.User.WreateRefreshToken(model.RefreshTokens{
		User_id: resp.Id,
		Token:   tokens.Refreshtoken,
		Exp:     time.Now().Add(time.Hour * 72).Unix(),
	})

	if err != nil {
		h.Log.Error("Error creating refresh token", "err", err)
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.SetCookie("access_token", tokens.Accestoken, 3600, "/", "", false, true)

	h.Log.Info("Logged in")
	ctx.JSON(200, tokens)

}

// LogOut godoc
// @Summary Log out a user
// @Description Log out a user by invalidating their refresh token and clearing cookies
// @Tags Authentication
// @Accept json
// @Produce json
// @Success 200 {string} string "Success"
// @Failure 400 {object} string "Invalid token or Error while finding user"
// @Router /auth/logout [post]
func (h *Handler) LogoutHandler(ctx *gin.Context) {
	h.Log.Info("Logged out")
	refresh, err := ctx.Cookie("access_token")
	if err != nil || refresh == "" {
		h.Log.Error("Access token")
		ctx.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	_, err = t.ExstractClaims(refresh)

	if err != nil {
		h.Log.Error("Error extracting claims", "err", err)
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	_, err = h.User.LogOut(ctx, &pb.TokenRequest{RefreshToken: refresh})
	if err != nil {
		h.Log.Error("Error logging out", "err", err)
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.SetCookie("access_token", "", -1, "/", "", false, true)
	ctx.SetCookie("refresh_token", "", -1, "/", "", false, true)

	h.Log.Info("Refresh token")
	ctx.JSON(200, gin.H{"message": "Logged out"})

}

// RefreshToken godoc
// @Summary Refresh user tokens
// @Description Refresh the user's access and refresh tokens using the refresh token from cookies
// @Tags Authentication
// @Accept json
// @Produce json
// @Param Authorization header string true "Access Token"
// @Success 200 {object} auth.TokenResponce
// @Failure 400 {object} string "Error while getting refresh token, Invalid token, or Error while refreshing"
// @Router /auth/refreshtoken [post]
func (h *Handler) RefreshToken(ctx *gin.Context) {

	h.Log.Info("Refreshing token")
	refresh := ctx.GetHeader("Authorization")
	if refresh == "" {
		h.Log.Error("Refresh token")
		ctx.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	_, err := t.ExstractClaims(refresh)

	if err != nil {
		h.Log.Error("Error extracting claims", "err", err)
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	tokens, err := h.User.RefreshToken(ctx, &pb.RefreshTokenRequest{RefreshToken: refresh})

	if err != nil {
		h.Log.Error("Error refreshing token", "err", err)
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	h.Log.Info("Refreshing token")
	ctx.JSON(200, tokens)
}
