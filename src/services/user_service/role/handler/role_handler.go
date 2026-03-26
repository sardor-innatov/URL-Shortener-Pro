package handler

import (
	"log"
	"net/http"
	"url_shortener_pro/src/common/middleware"
	"url_shortener_pro/src/services/user_service/role/dto"
	"url_shortener_pro/src/services/user_service/role/service"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type RoleHandler interface {
}

type roleHandler struct {
	//	db             *gorm.DB
	service service.RoleService
}

func NewRoleHandler(group *echo.Group, db *gorm.DB) {

	handler := &roleHandler{
		service: service.NewRoleService(db),
	}

	scoreGroup := group.Group("/role", middleware.RequireAuth, middleware.CheckPermission(db))
	{
		scoreGroup.POST("", handler.Create)
		scoreGroup.PUT("/:id", handler.Update)
		scoreGroup.DELETE("/:id", handler.Delete)
		scoreGroup.GET("/:id", handler.GetById)
		scoreGroup.GET("", handler.GetAllRoles)
		scoreGroup.GET("/permissions", handler.GetAllPermissions)
		// scoreGroup.POST("/give", handler.GiveRole)
	}
}

// @Summary      create
// @Tags         role
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        role   body   dto.RoleCreateDto  true  "role"
// @Success      200  {object}  nil
// @Failed       400  {object}  nil
// @Failed       500  {object}  nil
// @Router       /role [POST]
func (h *roleHandler) Create(ctx echo.Context) error {

	var role dto.RoleCreateDto
	log.Println("came to handler")
	err := h.service.CreateRole(ctx, &role)
	{
		if err != nil {
			return err
		}
	}

	return nil
}

// @Summary      update
// @Tags         role
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path   int64  true  "id"
// @Param        role   body   dto.RoleCreateDto  true  "role"
// @Success      200  {object}  nil
// @Failed       400  {object}  nil
// @Failed       404  {object}  nil
// @Failed       500  {object}  nil
// @Router       /role/{id} [PUT]
func (h *roleHandler) Update(ctx echo.Context) error {
	var role dto.RoleCreateDto

	err := h.service.Update(ctx, &role)
	{
		if err != nil {
			return err
		}
	}

	return nil
}

// @Summary      delete
// @Tags         role
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path   int64  true  "id"
// @Success      200  {object}  nil
// @Failed       400  {object}  nil
// @Failed       404  {object}  nil
// @Failed       500  {object}  nil
// @Router       /role/{id} [DELETE]
func (h *roleHandler) Delete(ctx echo.Context) error {

	err := h.service.Delete(ctx)
	{
		if err != nil {
			return err
		}
	}

	return nil
}

// @Summary      get by id
// @Tags         role
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path   int64  true  "id"
// @Success      200  {object}  nil
// @Failed       400  {object}  nil
// @Failed       404  {object}  nil
// @Failed       500  {object}  nil
// @Router       /role/{id} [GET]
func (h *roleHandler) GetById(ctx echo.Context) error {

	role, err := h.service.GetById(ctx)
	{
		if err != nil {
			return err
		}
	}

	return ctx.JSON(http.StatusOK, echo.Map{"role": role})
}

// @Summary      get all
// @Tags         role
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200  {object}  nil
// @Failed       400  {object}  nil
// @Failed       404  {object}  nil
// @Failed       500  {object}  nil
// @Router       /role [GET]
func (h *roleHandler) GetAllRoles(ctx echo.Context) error {

	roles, err := h.service.GetAllRoles(ctx)
	{
		if err != nil {
			return err
		}
	}

	return ctx.JSON(http.StatusOK, echo.Map{"roles": roles})
}

// @Summary      get all
// @Tags         role
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200  {object}  nil
// @Failed       400  {object}  nil
// @Failed       404  {object}  nil
// @Failed       500  {object}  nil
// @Router       /role/permissions [GET]
func (h *roleHandler) GetAllPermissions(ctx echo.Context) error {

	permissions, err := h.service.GetAllPermissions(ctx)
	{
		if err != nil {
			return err
		}
	}

	return ctx.JSON(http.StatusOK, echo.Map{"permissions": permissions})
}

// // @Summary      give role
// // @Tags         role
// // @Accept       json
// // @Produce      json
// // @Security     ApiKeyAuth
// // @Param        dto   body   dto.UserRoleCreate  true  "dto"
// // @Success      200  {object}  nil
// // @Failed       400  {object}  nil
// // @Failed       404  {object}  nil
// // @Failed       500  {object}  nil
// // @Router       /role/give [POST]
// func (h *roleHandler) GiveRole(ctx echo.Context) error {

// 	var dto dto.UserRoleCreate

// 	err := h.service.GiveRole(ctx, &dto)
// 	{
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return ctx.JSON(http.StatusOK, echo.Map{})
// }
