package cmd

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "url_shortener_pro/docs"
	"url_shortener_pro/src/common/config"
	custom_middleware "url_shortener_pro/src/common/middleware"
	"url_shortener_pro/src/common/seeder"
	link_cmd "url_shortener_pro/src/services/link_service"
	click_service "url_shortener_pro/src/services/link_service/click/service"
	user_cmd "url_shortener_pro/src/services/user_service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func Exec() {

	var router = echo.New()

	custom_middleware.InitLogger(router)

	var (
		projectEnv       = configureEnv(router)
		pgConfigInstance = config.NewConfig(
			projectEnv.PgHost,
			projectEnv.PgUser,
			projectEnv.PgPassword,
			projectEnv.PgDB,
			projectEnv.PgPort,
		)
	)

	db, err := config.NewGorm(&pgConfigInstance)
	{
		if err != nil {
			panic(err)
		}
	}

	migrate(db)

	clickSvc := click_service.NewClickService(db, projectEnv.ClickWorkers)

	config.RegisterValidator(router)

	router.GET("/swagger/*", echoSwagger.WrapHandler)

	{
		user_cmd.Cmd(router, db)
		link_cmd.Cmd(router, db, clickSvc)
	}

	seeder.SeedRoles(db)
	createOrg(db, router)

	router.Use(custom_middleware.RateLimitMiddleware(config.GetRedis()))
	router.Use(middleware.BodyLimit("1M"))

	server := &http.Server{
		Addr:              ":8080", // или берите из конфига
		Handler:           router,  // наш echo router
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
		ReadHeaderTimeout: 2 * time.Second, // Защита от Slowloris
	}

	go func() {
		err := server.ListenAndServe()
		{
			if err != nil && err != http.ErrServerClosed {
				log.Fatalf("Server error: %v", err)
			}
		}

	}()

	// 2. Создаем канал для прослушивания сигналов ОС (Ctrl+C, kill)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Программа "замирает" здесь, пока не придет сигнал
	<-stop
	println("\nShutdown signal received. Cleaning up...")

	// 3. Устанавливаем таймаут для закрытия Echo (например, 10 секунд)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := router.Shutdown(ctx); err != nil {
		println("Echo shutdown error: %v", err)
	}

	// 4. ГАСИМ ВОРКЕРОВ (ждем, пока допишут клики из очереди)
	clickSvc.Shutdown()

	println("Server stopped gracefully")
}
