package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

func RateLimitMiddleware(redisClient *redis.Client) echo.MiddlewareFunc {

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// 1. Получаем IP пользователя
			ip := c.RealIP()
			key := fmt.Sprintf("rate_limit:%s", ip)

			// 2. Инкрементируем значение в Redis
			count, err := redisClient.Incr(c.Request().Context(), key).Result()
			if err != nil {
				return err
			}

			// 3. Если это новый ключ, устанавливаем время жизни окна (1 минута)
			if count == 1 {
				redisClient.Expire(c.Request().Context(), key, time.Minute)
			}

			// 4. Проверяем лимит
			if count > 100 {
				return c.JSON(http.StatusTooManyRequests, echo.Map{
					"error": "Too Many Requests. Please wait a minute.",
				})
			}

			return next(c)
		}
	}
}
