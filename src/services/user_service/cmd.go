package user_cmd

import (
	auth_handler "url_shortener_pro/src/services/user_service/auth/handler"
	role_handler "url_shortener_pro/src/services/user_service/role/handler"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func Cmd(router *echo.Echo, db *gorm.DB) {

	routerGroup := router.Group("/api/v1")
	{
		auth_handler.NewAuthHandler(routerGroup, db)
		role_handler.NewRoleHandler(routerGroup, db)
	}
}
