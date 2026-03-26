package service

import (
	"log"
	"net/http"
	"strconv"
	role_dto "url_shortener_pro/src/services/user_service/role/dto"
	role_model "url_shortener_pro/src/services/user_service/role/model"
	user_model "url_shortener_pro/src/services/user_service/user/model"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type RoleService interface {
	CreateRole(ctx echo.Context,dto *role_dto.RoleCreateDto) error
	Update(ctx echo.Context,dto *role_dto.RoleCreateDto) error
	Delete(ctx echo.Context) error
	GetById(ctx echo.Context) (*role_model.Role, error)
	GetAllRoles(ctx echo.Context) ([]role_model.Role, error)
	GetAllPermissions(ctx echo.Context) ([]role_model.Permission, error)
	GiveRole(ctx echo.Context, dto *role_dto.UserRoleCreate)(error)
}

type roleService struct {
	db *gorm.DB
}

func NewRoleService(db *gorm.DB) RoleService {
	return &roleService{
		db: db,
	}
}

func (s *roleService) CreateRole(ctx echo.Context,dto *role_dto.RoleCreateDto) error {
	
	err := ctx.Bind(&dto)
	{
		if err != nil {
			log.Println("failed to read from json", err)
			return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
		}
	}
	log.Println("bindong json")

	role := role_model.Role{
		RoleName: dto.RoleName,
		Permissions: dto.Permissions,
	}

	result := s.db.Create(&role)
	err = result.Error
	{
		if err != nil {
			log.Println("failed to create role", err.Error())
			return ctx.JSON(http.StatusInternalServerError, echo.Map{"error" : err.Error()})
		}
	}

	return ctx.JSON(http.StatusBadRequest, echo.Map{"id": role.Id})
}

func (s *roleService) Update(ctx echo.Context,dto *role_dto.RoleCreateDto) error {

	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	{
		if err != nil {
			log.Println(err.Error())
			return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid quiz ID"})
		}
	}

	err = ctx.Bind(&dto)
	{
		if err != nil {
			log.Println("failed to read from json", err)
			return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
		}
	}

	var role role_model.Role
	err = s.db.Table("roles").
		Where("id = ?", id).
		First(&role).Error
	{
		if err != nil{
			log.Println("role not found", err.Error())
			return ctx.JSON(http.StatusNotFound, echo.Map{"error": err.Error()})
		}
	}

	role = role_model.Role{
		Id: id,
		RoleName: dto.RoleName,
		Permissions: dto.Permissions,
	}

	result := s.db.Save(&role)
	err = result.Error
	{
		if err != nil {
			log.Println("failed to update role", err.Error())
			return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}
	}

	return nil
}

func (s *roleService) Delete(ctx echo.Context) error {

	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	{
		if err != nil {
			log.Println(err.Error())
			return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid quiz ID"})
		}
	}

	var role role_model.Role
	err = s.db.Table("roles").
		Where("id = ?", id).
		First(&role).Error
	{
		if err != nil{
			log.Println("role not found", err.Error())
			return ctx.JSON(http.StatusNotFound, echo.Map{"error": err.Error()})
		}
	}

	result := s.db.Delete(&role)
	err = result.Error
	{
		if err != nil {
			log.Println("failed to delete", err.Error())
			return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}
	}

	return nil
}

func (s *roleService) GetById(ctx echo.Context) (*role_model.Role, error) {

	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	{
		if err != nil {
			log.Println(err.Error())
			return nil, ctx.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid role id"})
		}
	}

	var role role_model.Role
	err = s.db.Table("roles").
		Where("id = ?", id).
		First(&role).Error
	{
		if err != nil{
			log.Println("role not found", err.Error())
			return nil, ctx.JSON(http.StatusNotFound, echo.Map{"error": err.Error()})
		}
	}

	return &role, nil
}

func (s *roleService) GetAllRoles(ctx echo.Context) ([]role_model.Role, error) {

	var roles []role_model.Role

	result := s.db.Find(&roles)
	err := result.Error
	{
		if err != nil{
			log.Println("failed to get roles", err.Error())
			return nil, ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}
	}

	return roles, nil
}

func (s *roleService) GetAllPermissions(ctx echo.Context) ([]role_model.Permission, error) {

	var permissions []role_model.Permission

	result := s.db.Find(&permissions)
	err := result.Error
	{
		if err != nil{
			log.Println("failed to get roles", err.Error())
			return nil, ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}
	}

	return permissions, nil
}

func (s *roleService) GiveRole(ctx echo.Context, dto *role_dto.UserRoleCreate) (error) {
	err := ctx.Bind(&dto)
	{
		if err != nil {
			log.Println("failed to read from json", err)
			return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
		}
	}
	log.Println("bindong json")
	// checking role existence
	var role role_model.Role
	err = s.db.Table("roles").
		Where("id = ?", dto.RoleId).
		First(&role).Error
	{
		if err != nil {
			log.Println("role not found", err.Error())
			return ctx.JSON(http.StatusNotFound, echo.Map{"error" : err.Error()})
		}
	}
	log.Println("checked role existance")
	// checking user existence
	var user user_model.User
	err = s.db.Table("users").
		Where("id = ?", dto.UserId).
		First(&user).Error
	{
		if err != nil {
			log.Println("user not found", err.Error())
			return ctx.JSON(http.StatusNotFound, echo.Map{"error" : err.Error()})
		}
	}
	log.Println("checked user existence")

	result := s.db.Table("user_roles").Create(&dto)
	err = result.Error
	{
		if err != nil {
			log.Println("failed to give role", err.Error())
			return ctx.JSON(http.StatusInternalServerError, echo.Map{"error" : err.Error()})
		}
	}
	log.Println("created role")

	return nil
}