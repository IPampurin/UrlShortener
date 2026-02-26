package api

import (
	"context"
	"net/http"
	"time"

	"github.com/IPampurin/UrlShortener/pkg/db"
	"github.com/IPampurin/UrlShortener/pkg/service"
	"github.com/gin-gonic/gin"
	"github.com/wb-go/wbf/logger"
)

// CreateShortLink обрабатывает POST /shorten
func CreateShortLink(svc *service.Service, log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req CreateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Ctx(c.Request.Context()).Error("неверный формат запроса", "error", err)
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "неверный формат запроса"})
			return
		}

		link, err := svc.CreateLink(c.Request.Context(), req.OriginalURL, req.CustomShort)
		if err != nil {
			log.Ctx(c.Request.Context()).Error("ошибка создания ссылки", "error", err)
			// Проверяем известные ошибки
			if err.Error() == "короткая ссылка уже занята" {
				c.JSON(http.StatusConflict, ErrorResponse{Error: err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "внутренняя ошибка сервера"})
			return
		}

		resp := FullResponseLink{
			ShortURL:    link.ShortURL,
			OriginalURL: link.OriginalURL,
			CreatedAt:   link.CreatedAt,
			ClicksCount: link.ClicksCount,
		}

		c.JSON(http.StatusCreated, resp)
	}
}

// Redirect обрабатывает GET /s/:short_url
func Redirect(svc *service.Service, log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {

		shortURL := c.Param("short_url")
		link, err := svc.GetLinkByShortURL(c.Request.Context(), shortURL)
		if err != nil {
			log.Ctx(c.Request.Context()).Error("ошибка получения ссылки", "error", err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "внутренняя ошибка"})
			return
		}
		if link == nil {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "ссылка не найдена"})
			return
		}

		// Асинхронно записываем аналитику
		go func() {
			// создаём новый контекст, чтобы не зависеть от запроса
			ctx := context.Background()
			analytics := &db.Analytics{
				LinkID:     link.ID,
				AccessedAt: time.Now(),
				UserAgent:  c.GetHeader("User-Agent"),
				IPAddress:  c.ClientIP(),
				Referer:    c.GetHeader("Referer"),
			}
			if err := svc.SaveAnalytics(ctx, analytics); err != nil {
				log.Ctx(ctx).Error("ошибка записи аналитики", "error", err)
			}
			// увеличиваем счётчик переходов
			if err := svc.IncrementClicks(ctx, int64(link.ID)); err != nil {
				log.Ctx(ctx).Error("ошибка увеличения счётчика", "error", err)
			}
		}()

		c.Redirect(http.StatusFound, link.OriginalURL)
	}
}

// GetAnalytics обрабатывает GET /analytics/:short_url
func GetAnalytics(svc *service.Service, log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {

		shortURL := c.Param("short_url")
		link, err := svc.GetLinkByShortURL(c.Request.Context(), shortURL)
		if err != nil {
			log.Ctx(c.Request.Context()).Error("ошибка получения ссылки", "error", err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "внутренняя ошибка"})
			return
		}
		if link == nil {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "ссылка не найдена"})
			return
		}

		// получаем аналитику
		analytics, err := svc.GetAnalyticsByLinkID(c.Request.Context(), link.ID)
		if err != nil {
			log.Ctx(c.Request.Context()).Error("ошибка получения аналитики", "error", err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "внутренняя ошибка"})
			return
		}

		// преобразуем в ClickEntry
		clicks := make([]ClickEntry, len(analytics))
		for i, a := range analytics {
			clicks[i] = ClickEntry{
				AccessedAt: a.AccessedAt,
				UserAgent:  a.UserAgent,
				IPAddress:  a.IPAddress,
				Referer:    a.Referer,
			}
		}

		resp := AnalyticsResponse{
			ShortURL:    link.ShortURL,
			OriginalURL: link.OriginalURL,
			CreatedAt:   link.CreatedAt,
			ClicksCount: link.ClicksCount,
			Clicks:      clicks,
		}

		c.JSON(http.StatusOK, resp)
	}
}

// GetLinks обрабатывает GET /links (список последних ссылок для UI)
func GetLinks(svc *service.Service, log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {

		links, err := svc.GetLinks(c.Request.Context())
		if err != nil {
			log.Ctx(c.Request.Context()).Error("ошибка получения списка ссылок", "error", err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "внутренняя ошибка"})
			return
		}

		resp := make([]FullResponseLink, len(links))
		for i, l := range links {
			resp[i] = FullResponseLink{
				ShortURL:    l.ShortURL,
				OriginalURL: l.OriginalURL,
				CreatedAt:   l.CreatedAt,
				ClicksCount: l.ClicksCount,
			}
		}

		c.JSON(http.StatusOK, resp)
	}
}
