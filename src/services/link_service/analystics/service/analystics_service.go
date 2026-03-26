package service

import (
	"log"
	"net/http"
	"strconv"
	"time"
	"url_shortener_pro/src/services/link_service/analystics/dto"
	model "url_shortener_pro/src/services/link_service/analystics/model"
	click_model "url_shortener_pro/src/services/link_service/click/model"
	link_model "url_shortener_pro/src/services/link_service/link/model"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type AnalyticsService interface {
	Analistics(ctx echo.Context) error
}

type analyticsService struct {
	db *gorm.DB
}

func NewAnalysticsService(db *gorm.DB) AnalyticsService {
	service := analyticsService{
		db: db,
	}

	go service.RunStatsAggregator()

	return &service
}

func (s *analyticsService) Analistics(ctx echo.Context) error {
	idStr := ctx.Param("id")
	linkID, err := strconv.ParseInt(idStr, 10, 64)
	{
		if err != nil {
			println(err.Error())
			return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid link ID"})
		}
	}

	role := ctx.Get("role")
	val := ctx.Get("id")
	idFloat, ok := val.(float64)
	{
		if !ok {
			log.Println("invalid user id type in context")
			return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "invalid user id type in context"})
		}
	}
	userId := int64(idFloat)

	if role == "user" {
		var link link_model.Link
		err = s.db.Model(&link).
			Where("id = ? AND user_id = ? AND is_deleted = false ", linkID, userId).
			First(&link).Error
		{
			if err != nil {
				return ctx.JSON(http.StatusNotFound, echo.Map{"error": "link not found"})
			}
		}
	}

	var stats model.LinkStats
	s.db.Where("link_id = ?", linkID).First(&stats)

	var topCountries []dto.CountryStatsDto
	// Считаем ТОП стран (индекс по link_id и country обязателен!)
	s.db.Model(&click_model.Click{}).
		Select("country, count(id) as count").
		Where("link_id = ?", linkID).
		Group("country").
		Order("count DESC").
		Limit(5).
		Scan(&topCountries)

	var uniqueIPs int64
	s.db.Model(&click_model.Click{}).
		Where("link_id = ?", linkID).
		Distinct("ip_address").
		Count(&uniqueIPs)

	response := dto.StatsDto{
		TotalClicks:   stats.TotalClicks,
		UniqueIps:     uniqueIPs,
		TopCountries:  topCountries,
		Last24hClicks: stats.Last24hClicks,
	}

	return ctx.JSON(http.StatusOK, echo.Map{"data": response})
}

func (s *analyticsService) RunStatsAggregator() {
	ticker := time.NewTicker(1 * time.Minute)
	for range ticker.C {
		s.aggregateStats()
		s.update24hStats()
	}
}

func (s *analyticsService) aggregateStats() {

	var results []struct {
		LinkID uint
		Count  int
	}

	s.db.Model(&click_model.Click{}).
		Select("link_id, count(id) as count").
		Group("link_id").
		Scan(&results)

	for _, res := range results {
		s.db.Exec(`
            INSERT INTO link_stats (link_id, total_clicks, updated_at)
            VALUES (?, ?, NOW())
            ON CONFLICT (link_id) 
            DO UPDATE SET 
                total_clicks = EXCLUDED.total_clicks,
                updated_at = NOW()`,
			res.LinkID, res.Count)
	}
}

func (s *analyticsService) update24hStats() {

	err := s.db.Exec(`
    UPDATE link_stats
    SET 
    last_24h_clicks = COALESCE((
        SELECT count(*) 
        FROM clicks 
        WHERE clicks.link_id = link_stats.link_id 
          AND clicks.created_at > NOW() - INTERVAL '24 hours'
    ), 0),
    updated_at = NOW();
    `).Error

	if err != nil {
		log.Printf("Error updating 24h stats: %v", err)
	}
}
