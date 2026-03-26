package handler

import (
	"url_shortener_pro/src/common/config"
	"url_shortener_pro/src/common/middleware"
	click_service "url_shortener_pro/src/services/link_service/click/service"
	analystics_service "url_shortener_pro/src/services/link_service/analystics/service"
	"url_shortener_pro/src/services/link_service/link/service"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type LinkHandler interface {
}

type linkHandler struct {
	service service.LinkService
	stats analystics_service.AnalyticsService
}

func NewLinkHandler(router *echo.Echo, group *echo.Group, db *gorm.DB, clickSvc click_service.ClickService) {
	handler := &linkHandler{
		service: service.NewLinkService(db, clickSvc),
		stats: analystics_service.NewAnalysticsService(db),
	}

	redisClient := config.GetRedis()

	linkGroup := router.Group("api/v1/link")
	{
		linkGroup.POST("/shorten", handler.Create, middleware.RequireAuth, middleware.CheckPermission(db))
		linkGroup.DELETE("/delete/:id", handler.Delete,middleware.RequireAuth, middleware.CheckPermission(db))
		linkGroup.GET("/my", handler.GetMyLinks,middleware.RequireAuth, middleware.CheckPermission(db))
		linkGroup.GET("/:id/stats", handler.GetStats,middleware.RequireAuth, middleware.CheckPermission(db))
	}

	router.GET("/:shortCode", handler.Redirect , middleware.RateLimitMiddleware(redisClient))
}

// @Summary      create short link
// @Tags         link
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        request   body   dto.LinkCreateDto  true  "request"
// @Success      200  {object}  nil
// @Failed       400  {object}  error
// @Failed       401  {object}  error
// @Failed       500  {object}  error
// @Router       /link/shorten [POST]
func (h *linkHandler) Create(ctx echo.Context) error {

	err := h.service.Create(ctx)
	{
		if err != nil {
			return err
		}
	}

	return nil
}

// @Summary      delete link
// @Tags         link
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path   int64  true  "id"
// @Success      204  {object}  nil
// @Failed       401  {object}  error
// @Failed       400  {object}  error
// @Failed       404  {object}  error
// @Failed       500  {object}  error
// @Router       /link/delete/{id} [DELETE]
func (h *linkHandler) Delete(ctx echo.Context) error {

	err := h.service.Delete(ctx)
	{
		if err != nil {
			return err
		}
	}

	return nil
}

// @Summary      redirect
// @Tags         link
// @Accept       json
// @Produce      json
// @Param        shortCode   path   string  true  "shortCode"
// @Success      302  {object}  nil
// @Failed       404  {object}  error
// @Failed       410  {object}  error
// @Failed       500  {object}  error
// @Router       /{shortCode} [GET]
func (h *linkHandler) Redirect(ctx echo.Context) error {

	err := h.service.Redirect(ctx)
	{
		if err != nil {
			return err
		}
	}

	return nil
}
// @Summary      get link by user id
// @Tags         link
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200  {object}  nil
// @Failed       404  {object}  error
// @Failed       401  {object}  error
// @Failed       500  {object}  error
// @Router       /link/my [GET]
func (h *linkHandler) GetMyLinks(ctx echo.Context) error {

	err := h.service.GetMyLinks(ctx)
	{
		if err != nil {
			return err
		}
	}

	return nil
}

// @Summary      get stats about link
// @Tags         stats
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path   int64  true  "id"
// @Success      200  {object}  nil
// @Failed       400  {object}  error
// @Failed       401  {object}  error
// @Failed       500  {object}  error
// @Router       /link/{id}/stats [GET]
func (h *linkHandler) GetStats(ctx echo.Context)error{

	err := h.stats.Analistics(ctx)
	{
		if err != nil {
			return err
		}
	}

	return nil
}