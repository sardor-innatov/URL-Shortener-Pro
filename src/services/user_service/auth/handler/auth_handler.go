package auth_handler

import (
	"url_shortener_pro/src/common/middleware"
	auth_dto "url_shortener_pro/src/services/user_service/auth/dto"

	"fmt"
	"log"
	"net/http"
	auth_service "url_shortener_pro/src/services/user_service/auth/service"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type AuthHandler interface {
}

type authHandler struct {
	//	db             *gorm.DB
	service auth_service.AuthService
}

func NewAuthHandler(group *echo.Group, db *gorm.DB) {

	handler := &authHandler{
		service: auth_service.NewAuthService(db),
	}

	authGroup := group.Group("/auth")
	{
		authGroup.POST("", handler.SignUp)
		authGroup.POST("/login", handler.LogIn)
		authGroup.GET("/me", handler.GetMe,middleware.RequireAuth)
		authGroup.POST("/refresh", handler.RefreshToken,middleware.RequireAuth)
	}
}

// @Summary      sign up
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        auth   body   auth_dto.SignUpDto  true  "auth"
// @Success      200  {object}  nil
// @Failed       400  {object}  nil
// @Failed       500  {object}  nil
// @Router       /auth [POST]
func (h *authHandler) SignUp(ctx echo.Context) error {
	fmt.Println("signup")
	var dto auth_dto.SignUpDto

	err := h.service.SignUp(ctx, &dto)
	{
		if err != nil {

			return err
		}
	}

	return nil
}

// @Summary      login
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        auth   body   auth_dto.LogInDto  true  "auth"
// @Success      200  {object}  nil
// @Failed       400  {object}  nil
// @Failed       404  {object}  nil
// @Failed       500  {object}  nil
// @Router       /auth/login [POST]
func (h *authHandler) LogIn(ctx echo.Context) error {
	log.Println("login")
	var dto auth_dto.LogInDto

	err := h.service.LogIn(ctx, &dto)
	{
		if err != nil {
			return err
		}
	}

	return nil
}

// @Summary      get my profile
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200  {object}  nil
// @Failed       400  {object}  nil
// @Failed       404  {object}  nil
// @Failed       500  {object}  nil
// @Router       /auth/me [GET]
func (h *authHandler) GetMe(ctx echo.Context) error {
	log.Println("get me")

	dto, err := h.service.Me(ctx)
	{
		if err != nil {
			return err
		}
	}

	return ctx.JSON(http.StatusOK, echo.Map{"user": dto})
}

// @Summary      refresh token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        auth   body   auth_dto.RefreshTokenDto  true  "auth"
// @Success      200  {object}  nil
// @Failed       400  {object}  nil
// @Failed       404  {object}  nil
// @Failed       500  {object}  nil
// @Router       /auth/refresh [POST]
func (h *authHandler) RefreshToken(ctx echo.Context) error {
	log.Println("refresh token")
	var dto auth_dto.RefreshTokenDto

	err := h.service.RefreshToken(ctx, &dto)
	{
		if err != nil {
			return err
		}
	}

	return nil
}
