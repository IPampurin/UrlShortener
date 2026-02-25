package cache

import (
	"context"

	"github.com/IPampurin/UrlShortener/pkg/configuration"
	"github.com/wb-go/wbf/logger"
)

/*
в кэше должны хранится пары ShortURL - OriginalURL
также должны быть сочетания OriginalURL - ShortURL // при чём ShortURL в этом случае только сгенерированные! (кастомные ShortURL надо проверять первыми)
*/

func InitRedis(ctx context.Context, cfgCache *configuration.ConfCache, log logger.Logger) error {

	return nil
}
