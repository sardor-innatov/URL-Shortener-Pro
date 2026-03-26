package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"url_shortener_pro/src/common/config"
	click_service "url_shortener_pro/src/services/link_service/click/service"
	"url_shortener_pro/src/services/link_service/link/dto"
	"url_shortener_pro/src/services/link_service/link/model"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type LinkService interface {
	Create(ctx echo.Context) error
	Redirect(ctx echo.Context) error
	Delete(ctx echo.Context) error
	GetMyLinks(ctx echo.Context) error
}

type linkService struct {
	db       *gorm.DB
	rdb      *redis.Client
	clickSvc click_service.ClickService
}

func NewLinkService(db *gorm.DB, clickSvc click_service.ClickService) LinkService {
	return &linkService{
		db:       db,
		rdb:      config.GetRedis(),
		clickSvc: clickSvc,
	}
}

func (s *linkService) Create(ctx echo.Context) error {

	var dto dto.LinkCreateDto
	err := ctx.Bind(&dto)
	{
		if err != nil {
			print(err.Error())
			return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "invalid json"})
		}
	}

	val := ctx.Get("id")
	idFloat, ok := val.(float64)
	{
		if !ok {
			log.Println("invalid user id type in context")
			return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "invalid user id type in context"})
		}
	}

	id := int64(idFloat)

	dto.URL = prepareURL(dto.URL)

	err = validate(ctx, dto)
	{
		if err != nil {
			return err
		}
	}

	var expiresAt *time.Time
	if dto.ExpiresIn != 0 {
		now := time.Now()
		expiresTime := now.Add(time.Duration(dto.ExpiresIn) * time.Hour)
		expiresAt = &expiresTime
	} else {
		expiresAt = nil
	}

	var shortcode string
	if dto.CustomAlias != "" {

		{
			shortcode = dto.CustomAlias
		}

	} else {

		for {

			{
				shortcode = generateShortCode(10)
			}

			link := model.Link{
				URL:       dto.URL,
				ShortCode: shortcode,
				ExpiresAt: expiresAt,
				UserId:    id,
			}

			timeOutCtx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
			defer cancel()

			err = s.db.WithContext(timeOutCtx).Create(&link).Error
			{
				if err != nil {
					continue
				}
			}

			envproj := config.ProjectEnv()
			shortURL := fmt.Sprintf("%s/%s", envproj.BaseURL, shortcode)

			return ctx.JSON(http.StatusOK, echo.Map{"shortLink": shortURL})
		}
	}

	link := model.Link{
		URL:       dto.URL,
		ShortCode: shortcode,
		ExpiresAt: expiresAt,
		UserId:    id,
	}

	timeOutCtx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	err = s.db.WithContext(timeOutCtx).Create(&link).Error
	{
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "this shortcode already used"})
		}
	}

	envproj := config.ProjectEnv()
	shortURL := fmt.Sprintf("%s/%s", envproj.BaseURL, shortcode)

	return ctx.JSON(http.StatusOK, echo.Map{"shortLink": shortURL})
}

func (s *linkService) Redirect(ctx echo.Context) error {

	shortCode := ctx.Param("shortCode")

	s.recordClick(ctx, shortCode)

	dataStr, err := s.rdb.Get(ctx.Request().Context(), shortCode).Result()
	{
		if err == nil {

			var dto dto.LinkRedisDto
			err = json.Unmarshal([]byte(dataStr), &dto)
			{
				if err != nil {
					return fmt.Errorf("could not marshal dto: %v", err)
				}
			}

			if dto.ExpiresAt != nil && dto.ExpiresAt.Before(time.Now()) {
				return ctx.JSON(http.StatusGone, echo.Map{"error": "link expired"})
			}

			println("from cache")
			return ctx.Redirect(http.StatusFound, dto.URL)
		}
	}

	timeOutCtx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	var link model.Link
	err = s.db.WithContext(timeOutCtx).Table("links").
		Where("short_code = ? AND deleted_at IS NULL", shortCode).
		First(&link).Error
	{
		if err != nil {
			return ctx.JSON(http.StatusNotFound, echo.Map{"error": "short link not found"})
		} else if link.ExpiresAt != nil && link.ExpiresAt.Before(time.Now()) {
			return ctx.JSON(http.StatusGone, echo.Map{"error": "link expired"})
		}
	}

	dto := dto.LinkRedisDto{
		URL:       link.URL,
		ExpiresAt: link.ExpiresAt,
	}

	dataJson, err := json.Marshal(dto)
	{
		if err != nil {
			return fmt.Errorf("could not marshal dto: %v", err)
		}
	}

	err = s.rdb.Set(ctx.Request().Context(), shortCode, dataJson, 10*time.Minute).Err()
	{
		if err != nil {
			panic("could not save to redis err: " + err.Error())
		}
	}

	return ctx.Redirect(http.StatusFound, link.URL)
}

func (s *linkService) Delete(ctx echo.Context) error {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	{
		if err != nil {
			println(err.Error())
			return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid user ID"})
		}
	}

	timeOutCtx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	var link model.Link
	err = s.db.WithContext(timeOutCtx).Table("links").
		Where("id = ? AND is_deleted = false", id).
		First(&link).Error
	{
		if err != nil {
			return ctx.JSON(http.StatusNotFound, echo.Map{"error": "link not found"})
		}
	}

	{
		link.IsDeleted = true

		now := time.Now()
		link.DeletedAt = &now
	}

	err = s.db.WithContext(timeOutCtx).Save(&link).Error
	{
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to delete"})
		}
	}
	err = s.rdb.Del(ctx.Request().Context(), link.ShortCode).Err()
	{
		if err != nil {
			return err
		}
	}

	return ctx.JSON(http.StatusNoContent, echo.Map{})

}
func (s *linkService) GetMyLinks(ctx echo.Context) error {
	val := ctx.Get("id")
	idFloat, ok := val.(float64)
	{
		if !ok {
			log.Println("invalid user id type in context")
			return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "invalid user id type in context"})
		}
	}

	id := int64(idFloat)

	timeOutCtx, cancel := context.WithTimeout(ctx.Request().Context(), 200*time.Millisecond)
	defer cancel()

	var links []dto.LinkGetDto
	err := s.db.WithContext(timeOutCtx).Table("links").
		Select("url", "short_code", "expires_at").
		Where("user_id = ?", id).
		Scan(&links).Error
	{
		if err != nil {
			return err
		}
	}

	return ctx.JSON(http.StatusOK, echo.Map{"data": links})

}

func generateShortCode(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func validate(ctx echo.Context, dto interface{}) error {

	err := ctx.Validate(dto)
	if err != nil {

		castedObject, ok := err.(validator.ValidationErrors)
		{
			if !ok {
				ctx.JSON(http.StatusInternalServerError, echo.Map{"error": "Internal error"})
				return err
			}
		}

		customErrors := map[string]string{
			"LinkCreateDto.URL.required":                 "url required",
			"LinkCreateDto.URL.http_url":                 "invalid url",
			"LinkCreateDto.CustomAlias.max":              "custom alias must be less then 10 characters",
			"LinkCreateDto.CustomAlias.min":              "custom alias must be at least 4 characters",
			"LinkCreateDto.CustomAlias.hostname_rfc1123": "custom alias can only contain a-z, A-Z, 0-9, '-'",
		}

		for _, err := range castedObject {
			key := fmt.Sprintf("%s.%s", err.Namespace(), err.Tag())
			if msg, ok := customErrors[key]; ok {
				ctx.JSON(http.StatusBadRequest, echo.Map{"error": msg})
				return err
			}
		}
	}

	return nil
}

func (s *linkService) recordClick(ctx echo.Context, shortCode string) {

	ip := ctx.RealIP()
	userAgent := ctx.Request().UserAgent()

	timeOutCtx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	var link model.Link
	err := s.db.WithContext(timeOutCtx).Table("links").
		Where("short_code = ?", shortCode).
		First(&link).Error
	{
		if err != nil {
			return
		}
	}

	s.clickSvc.Record(link.Id, ip, userAgent)
}

func prepareURL(url string) string {

	cleanURL := strings.ReplaceAll(url, " ", "")
	cleanURL = strings.TrimSpace(cleanURL)

	cleanURL = strings.ToLower(cleanURL)

	return cleanURL
}
