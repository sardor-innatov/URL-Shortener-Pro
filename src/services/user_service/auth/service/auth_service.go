package auth_service

import (
	"context"
	auth_dto "url_shortener_pro/src/services/user_service/auth/dto"
	user_dto "url_shortener_pro/src/services/user_service/user/dto"
	user_model "url_shortener_pro/src/services/user_service/user/model"
	user_service "url_shortener_pro/src/services/user_service/user/service"

	//user_model "edu_system/src/module/user_service/model"

	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService interface {
	SignUp(ctx echo.Context, dto *auth_dto.SignUpDto) error
	LogIn(ctx echo.Context, dto *auth_dto.LogInDto) error
	Me(ctx echo.Context) (*user_dto.UserGetDto, error)
	createUser(ctx echo.Context, dto *auth_dto.SignUpDto) error
	RefreshToken(ctx echo.Context, dto *auth_dto.RefreshTokenDto) error
}

type authService struct {
	user user_service.UserService
	db   gorm.DB
}

func NewAuthService(db *gorm.DB) AuthService {
	return &authService{
		user: user_service.NewUserService(db),
		db:   *db,
	}
}

func (a *authService) SignUp(ctx echo.Context, dto *auth_dto.SignUpDto) error {

	err := ctx.Bind(&dto)
	{
		if err != nil {
			log.Println("Failed to read request body:", err)
			return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "Failed to read request body"})
		}
	}

	err = a.createUser(ctx, dto)
	{
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "bad request: email must be unique"})
		}
	}
	timeOutCtx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	var user user_model.User
	err = a.db.WithContext(timeOutCtx).Table("users").
		Where("email = ?", dto.Email).
		First(&user).Error
	{
		if err != nil {
			log.Printf("failed to sign up")
			return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
		}
	}

	return ctx.JSON(http.StatusOK, echo.Map{"id": user.Id})
}

func (a *authService) LogIn(ctx echo.Context, dto *auth_dto.LogInDto) error {

	err := ctx.Bind(&dto)
	{
		if err != nil {
			log.Println("Failed to read request body:", err)
			return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "Failed to read request body"})
		}
	}

	// checking if user exists
	_, password, id, err := a.user.GetByEmail(ctx, dto.Email)
	if err != nil {
		return err
	}

	// checking password
	err = bcrypt.CompareHashAndPassword([]byte(password), []byte(dto.Password))
	{
		if err != nil {
			log.Println("Failed to compare password:", err)
			return ctx.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid password"})
		}
	}

	// getting role
	role, err := a.user.GetFirstRole(ctx, id)
	{
		if err != nil {
			return err
		}
	}

	// generate token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    id,
		"email": dto.Email,
		"role":  role,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	{
		if err != nil {
			log.Println("Failed to generate token:", err)
			return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to generate token"})
		}
	}

	return ctx.JSON(http.StatusOK, echo.Map{
		"token": tokenString,
	})
}

func (a *authService) RefreshToken(ctx echo.Context, dto *auth_dto.RefreshTokenDto) error {

	email := ctx.Get("email")
	val := ctx.Get("id")
	idFloat, ok := val.(float64)
	{
		if !ok {
			log.Println("invalid user id type in context")
			return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "invalid user id type in context"})
		}
	}

	id := int64(idFloat)
	log.Println(val)
	log.Println(idFloat)
	log.Println(id)

	// checking if user exists
	_, _, _, err := a.user.GetByEmail(ctx, email.(string))
	if err != nil {
		return err
	}

	err = ctx.Bind(&dto)
	{
		if err != nil {
			log.Println("Failed to read request body:", err)
			return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "Failed to read request body"})
		}
	}
	log.Println(dto.Role)

	timeOutCtx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	var count int64
	err = a.db.WithContext(timeOutCtx).Table("user_roles").
		Joins("JOIN roles on roles.id = user_roles.role_id").
		Where("user_roles.user_id = ? AND roles.role_name = ?", val, dto.Role).
		Count(&count).Error
	{
		if err != nil {
			log.Println("failed to check roles", err.Error())
			return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}
		if count == 0 {
			log.Println("role not found, count ", count)
			return ctx.JSON(http.StatusForbidden, echo.Map{"error": "you dont have such role"})
		}
	}
	// generate token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    id,
		"email": email,
		"role":  dto.Role,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	{
		if err != nil {
			log.Println("Failed to generate token:", err)
			return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to generate token"})
		}
	}

	return ctx.JSON(http.StatusOK, echo.Map{
		"token": tokenString,
	})

}

func (a *authService) Me(ctx echo.Context) (*user_dto.UserGetDto, error) {

	val := ctx.Get("id")
	log.Println(val)

	var id int64
	switch v := val.(type) {
	case float64:
		id = int64(v)
	case int64:
		id = v
	case int:
		id = int64(v)
	default:
		log.Printf("invalid user id type in context: %T", val)
		return nil, ctx.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized: invalid id type"})
	}

	timeOutCtx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	var dto user_dto.UserGetDto
	result := a.db.WithContext(timeOutCtx).Table("users").
		Where("id = ?", id).
		First(&dto)
	err := result.Error
	{
		if err != nil {
			log.Println("user not found", err)
			return nil, ctx.JSON(http.StatusNotFound, echo.Map{"error": err.Error()})
		}
	}

	return &dto, nil
}

func (a *authService) createUser(ctx echo.Context, dto *auth_dto.SignUpDto) error {

	var userDto user_dto.UserCreateDto
	{
		userDto.FirstName = dto.FirstName
		userDto.LastName = dto.LastName
		userDto.Email = dto.Email
		userDto.Password = dto.Password
	}

	err := a.user.Create(ctx, &userDto)
	{
		if err != nil {
			return err
		}
	}

	return nil
}
