package link_cmd

import (
	click_service "url_shortener_pro/src/services/link_service/click/service"
	link_handler "url_shortener_pro/src/services/link_service/link/handler"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func Cmd(router *echo.Echo, db *gorm.DB, clickSvc click_service.ClickService) {

	routerGroup := router.Group("/api/v1")
	{
		link_handler.NewLinkHandler(router, routerGroup, db, clickSvc)
	}
}
