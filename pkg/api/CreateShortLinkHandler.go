package api

import (
	"net/http"

	"github.com/IPampurin/UrlShortener/pkg/db"
	"github.com/IPampurin/UrlShortener/pkg/generator"
	"github.com/gin-gonic/gin"
	"github.com/wb-go/wbf/logger"
)

func CreateShortLinkHandler(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {

		var createShortLinkRequest CreateShortLinkRequest
		var linkApi *LinkApi
		var err error

		if err = c.ShouldBindJSON(&createShortLinkRequest); err != nil {
			log.Ctx(c.Request.Context()).Error("неверный формат запроса", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат запроса"})
			return
		}

		// если пользователь задал своё имя короткой ссылки
		if createShortLinkRequest.CustomShort != "" {
			// смотрим есть ли такая короткая ссылка в БД
			linkApi, err = withCustomShortLinkProcess(c, &createShortLinkRequest)
			if err != nil {
				log.Ctx(c.Request.Context()).Error("ошибка базы данных", "error", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка базы данных"})
				return
			}
		} else {
			// если пользователь не задал своё имя короткой ссылки
			linkApi, err = noCustomShortLinkProcess(c, &createShortLinkRequest)
			if err != nil {
				log.Ctx(c.Request.Context()).Error("ошибка базы данных", "error", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка базы данных"})
				return
			}
		}

		c.JSON(http.StatusCreated, &linkApi)
	}
}

// noCustomShortLinkProcess описывает действия, если пользователь не задал ShortURL
func noCustomShortLinkProcess(c *gin.Context, createShortLinkRequest *CreateShortLinkRequest) (*LinkApi, error) {

	var linkApi *LinkApi
	var linkDB *db.LinkDB

	// проверяем наличие OriginalURL в БД
	linkDB, err := db.GetClientDB().GetLinkByOriginalURL(c.Request.Context(), createShortLinkRequest.OriginalURL)
	if err != nil {
		return nil, err
	}
	// если в БД нет такого OriginalURL, генерируем ShortURL и отправляем на создание записи в БД
	if linkDB == nil {
		shortURL := generator.NewRandomString(0)
		// TODO проверить shortURL на уникальность, возможно генерировать в цикле
		linkDB, err = db.GetClientDB().CreateLink(c.Request.Context(), shortURL, createShortLinkRequest.OriginalURL, false)
		if err != nil {
			return nil, err
		}
	}
	// заполняем соответствующую структуру
	linkApi = &LinkApi{
		ShortURL:    linkDB.ShortURL,
		OriginalURL: linkDB.OriginalURL,
		CreatedAt:   linkDB.CreatedAt,
		ClicksCount: linkDB.ClicksCount,
	}

	return linkApi, nil
}

// withCustomShortLinkProcess описывает действия, если пользователь задал ShortURL
func withCustomShortLinkProcess(c *gin.Context, createShortLinkRequest *CreateShortLinkRequest) (*LinkApi, error) {

	var linkApi *LinkApi
	var linkDB *db.LinkDB

	// проверяем есть ли такая же ShortURL в БД
	linkDB, err := db.GetClientDB().GetLinkByShortURL(c.Request.Context(), createShortLinkRequest.OriginalURL)
	if err != nil {
		return nil, err
	}

	// проверяем наличие OriginalURL в БД
	linkDB, err = db.GetClientDB().GetLinkByOriginalURL(c.Request.Context(), createShortLinkRequest.OriginalURL)
	if err != nil {
		return nil, err
	}
	// если в БД нет такого OriginalURL, генерируем ShortURL и отправляем на создание записи в БД
	if linkDB == nil {
		shortURL := generator.NewRandomString(0)
		// TODO проверить shortURL на уникальность, возможно генерировать в цикле
		linkDB, err = db.GetClientDB().CreateLink(c.Request.Context(), shortURL, createShortLinkRequest.OriginalURL, false)
		if err != nil {
			return nil, err
		}
	}
	// заполняем соответствующую структуру
	linkApi = &LinkApi{
		ShortURL:    linkDB.ShortURL,
		OriginalURL: linkDB.OriginalURL,
		CreatedAt:   linkDB.CreatedAt,
		ClicksCount: linkDB.ClicksCount,
	}

	return linkApi, nil
}
