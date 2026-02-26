package service

import (
	"context"

	"github.com/IPampurin/UrlShortener/pkg/cache"
	"github.com/IPampurin/UrlShortener/pkg/db"
	"github.com/IPampurin/UrlShortener/pkg/interfaces"
	"github.com/wb-go/wbf/logger"
)

type Service struct {
	link      interfaces.LinkMethods
	analytics interfaces.AnalyticsMethods
	cache     interfaces.CacheMethods
}

func InitService(ctx context.Context, storage *db.DataBase, cache *cache.Cache, log logger.Logger) *Service {

	svc := &Service{
		link:      storage, // *db.DataBase реализует LinkMethods
		analytics: storage, // *db.DataBase реализует AnalyticsMethods
		cache:     cache,   // *cache.Cache реализует CacheMethods
	}

	return svc
}
