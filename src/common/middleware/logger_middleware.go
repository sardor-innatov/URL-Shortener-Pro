package middleware

import (
	"fmt"
	"log"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func InitLogger(router *echo.Echo) {
	// Обязательно первым, чтобы ID был доступен логгеру
	router.Use(middleware.RequestID())

	router.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:    true,
		LogURI:       true,
		LogMethod:    true,
		LogLatency:   true,
		LogError:     true,
		LogRequestID: true,
		HandleError:  true, // позволяет логгеру видеть ошибки из обработчиков
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			// Достаем user_id из контекста (ставится в JWT middleware)
			userID := "guest"
			if val := c.Get("user_id"); val != nil {
				userID = fmt.Sprintf("%v", val)
			}

			errMsg := ""
			if v.Error != nil {
				errMsg = v.Error.Error()
			}

			// Печать в терминал одной строкой
			log.Printf(`{"time":"%s", "request_id":"%s", "user_id":"%s", "method":"%s", "uri":"%s", "status":%d, "latency":"%s", "error":"%s"}`+"\n",
				v.StartTime.Format(time.RFC3339),
				v.RequestID,
				userID,
				v.Method,
				v.URI,
				v.Status,
				v.Latency.String(),
				errMsg,
			)
			return nil
		},
	}))

	router.Use(middleware.Recover())
}
