package api

import (
	"net/http"

	"github.com/IPampurin/UrlShortener/pkg/service"
	"github.com/gin-gonic/gin"
	"github.com/wb-go/wbf/logger"
)

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
			switch err {
			case service.ErrShortURLExists:
				c.JSON(http.StatusConflict, ErrorResponse{Error: "короткая ссылка уже занята"})
			default:
				c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "внутренняя ошибка сервера"})
			}
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
