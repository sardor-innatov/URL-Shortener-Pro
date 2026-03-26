package middleware

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"os"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {

		authHeader := ctx.Request().Header.Get("Authorization")
		fmt.Println(authHeader)
		{
			if authHeader == "" || len(authHeader) < 8 || authHeader[:7] != "Bearer " {
				log.Println("unauthorized")
				return ctx.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
			}
		}
		fmt.Println("good Bearer")
		tokenString := authHeader[7:]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			{
				if !ok {
					log.Println("unexpected signing method")
					return nil, ctx.JSON(http.StatusBadRequest, echo.Map{"error": "unexpected signing method"})
				}
			}

			return []byte(os.Getenv("SECRET")), nil
		})
		{
			if err != nil || !token.Valid {
				log.Println("invalid token", err.Error())
				return ctx.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
			}
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		{
			if !ok {
				log.Println("couldnt get claims")
				return ctx.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
			}
		}

		// check if token is expired
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			log.Println("token expired")
			return ctx.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
		}

		//find the user with token

		id, ok := claims["id"]
		{
			if !ok {
				log.Println("couldnt get id from claim")
				return ctx.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
			}
		}

		email, ok := claims["email"].(string)
		{
			if !ok {
				log.Println("couldnt get email from claim")
				return ctx.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
			}
		}

		role, ok := claims["role"].(string)
		fmt.Println(role)
		{
			if !ok {
				log.Println("couldnt get role from claim")
				return ctx.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
			}
		}

		ctx.Set("id", id)
		ctx.Set("email", email)
		ctx.Set("role", role)

		fmt.Println(ctx.Get("role"))
		return next(ctx)
	}
}

func CheckPermission(db *gorm.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			log.Println("checking permissions")
			// var roleID int64
			// {
			// 	roleVal := ctx.Get("role_id")
			// 	roleFloat := roleVal.(float64)
			// 	roleID = int64(roleFloat)
			// }
			role := ctx.Get("role")
			path := ctx.Path() // путь из Echo (например, /api/v1/order/:id)
			method := ctx.Request().Method
			log.Println(role)
			var count int64

			permission := map[string][]string{
				path: {method},
			}

			jsonValue, _ := json.Marshal(permission)
			log.Println(jsonValue)
			err := db.Table("roles").
				Where("role_name = ? AND permissions @> ?", role, jsonValue).
				Count(&count).Error
			{
				if err != nil {
					log.Println(err.Error())
					return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
				}
				if count == 0 {
					log.Println(0)
					return ctx.JSON(http.StatusForbidden, echo.Map{"error": "you dont have permission to do this"})
				}
			}
			log.Println("next")
			return next(ctx)
		}
	}
}
