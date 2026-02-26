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
	cache     interfaces.CacheMetods
}

func InitService(ctx context.Context, storage *db.DataBase, cache *cache.Cache, appLogger logger.Logger) (*Service, error) {

	svc := &Service{}

	return svc, nil
}
