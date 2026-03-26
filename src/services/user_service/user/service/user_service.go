package service

import (
	"context"
	"time"
	role_model "url_shortener_pro/src/services/user_service/role/model"
	user_dto "url_shortener_pro/src/services/user_service/user/dto"
	user_model "url_shortener_pro/src/services/user_service/user/model"

	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService interface {
	Create(ctx echo.Context, dto *user_dto.UserCreateDto) error
	Update(ctx echo.Context, dto *user_dto.UserCreateDto) error
	Delete(ctx echo.Context) error
	GetById(ctx echo.Context) (*user_dto.UserGetDto, error)
	GetAll(ctx echo.Context) ([]user_dto.UserGetDto, error)
	GetByEmail(ctx echo.Context, email string) (*user_dto.UserGetDto, string, int64, error)
	GetFirstRole(ctx echo.Context, id int64) (string, error)
}

type userService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) UserService {
	return &userService{
		db: db,
	}
}

func (s *userService) Create(ctx echo.Context, dto *user_dto.UserCreateDto) error {

	// err := ctx.Bind(&dto)
	// {
	// 	if err != nil {
	// 		log.Println("failed to read from json", err)
	// 		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	// 	}
	// }

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dto.Password), 10)
	{
		if err != nil {
			log.Println("failed to hash password", err.Error())
			return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}
	}

	user := user_model.User{
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		Email:     dto.Email,
		Password:  string(hashedPassword),
	}

	timeOutCtx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	result := s.db.WithContext(timeOutCtx).Create(&user)
	err = result.Error
	{
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *userService) Update(ctx echo.Context, dto *user_dto.UserCreateDto) error {

	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	{
		if err != nil {
			log.Println(err.Error())
			return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid user ID"})
		}
	}

	err = ctx.Bind(&dto)
	{
		if err != nil {
			log.Println("failed to read from json", err)
			return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
		}
	}

	timeOutCtx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	var user user_model.User
	result := s.db.WithContext(timeOutCtx).Table("users").
		Where("id = ?", id).
		First(&user)
	err = result.Error
	{
		if err != nil {
			log.Println("user not found", err)
			return ctx.JSON(http.StatusNotFound, echo.Map{"error": err.Error()})
		}
	}

	{
		user.FirstName = dto.FirstName
		user.LastName = dto.LastName
		user.Email = dto.Email
	}

	result = s.db.WithContext(timeOutCtx).Save(user)
	err = result.Error
	{
		if err != nil {
			log.Println("failed to update", err)
			return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}
	}

	return ctx.JSON(http.StatusOK, echo.Map{"message": "user updated"})
}

func (s *userService) Delete(ctx echo.Context) error {

	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	{
		if err != nil {
			log.Println(err.Error())
			return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid user ID"})
		}
	}

	timeOutCtx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	var user user_model.User
	result := s.db.WithContext(timeOutCtx).Where("id = ?", id).First(&user)
	err = result.Error
	{
		if err != nil {
			log.Println("user not found", err)
			return ctx.JSON(http.StatusNotFound, echo.Map{"error": err.Error()})
		}
	}

	result = s.db.WithContext(timeOutCtx).Delete(user)
	err = result.Error
	{
		if err != nil {
			log.Println("failed to delete", err)
			return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}
	}

	return ctx.JSON(http.StatusOK, echo.Map{"message": "user deleted"})
}

func (s *userService) GetById(ctx echo.Context) (*user_dto.UserGetDto, error) {

	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	{
		if err != nil {
			log.Println(err.Error())
			return nil, ctx.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid user ID"})
		}
	}

	timeOutCtx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	var user user_model.User
	result := s.db.WithContext(timeOutCtx).Where("id = ?", id).First(&user)
	err = result.Error
	{
		if err != nil {
			log.Println("user not found", err)
			return nil, ctx.JSON(http.StatusNotFound, echo.Map{"error": err.Error()})
		}
	}

	dto := user_dto.UserGetDto{
		Id:        user.Id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}

	return &dto, nil
}

func (s *userService) GetAll(ctx echo.Context) ([]user_dto.UserGetDto, error) {

	timeOutCtx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	var users []user_model.User
	result := s.db.WithContext(timeOutCtx).Find(&users)
	err := result.Error
	{
		if err != nil {
			log.Println("failed to get users")
			return nil, ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}
	}

	dtos := make([]user_dto.UserGetDto, len(users))

	for i, user := range users {
		dtos[i] = user_dto.UserGetDto{
			Id:        user.Id,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
		}
	}

	return dtos, nil
}

func (s *userService) GetByEmail(ctx echo.Context, email string) (*user_dto.UserGetDto, string, int64, error) {

	timeOutCtx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	var user user_model.User
	result := s.db.WithContext(timeOutCtx).Table("users").
		Where("email = ?", email).
		First(&user)
	err := result.Error
	{
		if err != nil {
			log.Println("user not found", err.Error())
			return nil, "", 0, ctx.JSON(http.StatusNotFound, echo.Map{"error": err.Error()})
		}
	}

	dto := user_dto.UserGetDto{
		Id:        user.Id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}
	password := user.Password
	id := user.Id

	return &dto, password, id, nil
}

func (s *userService) GetFirstRole(ctx echo.Context, id int64) (string, error) {

	timeOutCtx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	type roleName struct {
		RoleName string `gorm:"unique;not null;column:role_name"`
	}

	var role roleName
	err := s.db.WithContext(timeOutCtx).Table("roles").
		Select("roles.role_name").
		Joins("JOIN user_roles on user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", id).
		First(&role).Error
	{
		if role.RoleName == "" {

			s.db.Create(&role_model.UserRole{
				UserId: id,
				RoleId: 2, // "user"
			})
			return "user", nil
		}
		if err != nil {
			log.Println("err", err.Error())
			return "", ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}
	}

	return role.RoleName, nil
}
