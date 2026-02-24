package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/IPampurin/UrlShortener/pkg/db"
	"github.com/gin-gonic/gin"
	"github.com/wb-go/wbf/logger"
)

func CreateShortLinkHandler(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {

		var createShortLinkRequest CreateShortLinkRequest

		if err := c.ShouldBindJSON(&createShortLinkRequest); err != nil {
			log.Ctx(c.Request.Context()).Error("неверный формат запроса", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат запроса"})
			return
		}

		/*
			валидируем чуть-чуть

			проверяем нет ли в кэше такой OriginalURL string `json:"original_url" binding:"required,url"`

			// проверяем, нет ли уже в кэше (на всякий случай)
			if cachedJSON, err := cache.GetClientRedis().Get(c.Request.Context(), key); err == nil && cachedJSON != "" {
				c.JSON(http.StatusConflict, gin.H{"error": "уведомление с таким UID уже существует"})
				return
			}

		*/

		// сохраняем в БД
		if err := db.GetClientDB().CreateLink(c.Request.Context(), n); err != nil {
			log.Ctx(c.Request.Context()).Error("ошибка базы данных", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка базы данных"})
			return
		}

		// сериализуем в JSON и сохраняем в Redis
		if jsonData, err := json.Marshal(n); err == nil {
			_ = cache.GetClientRedis().SetWithExpiration(c.Request.Context(), key, string(jsonData), 10*time.Minute)
		} else {
			log.Ctx(c.Request.Context()).Error("ошибка маршаллинга уведомления", "error", err)
		}

		c.JSON(http.StatusCreated, n)
	}
}
