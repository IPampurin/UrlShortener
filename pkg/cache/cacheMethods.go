package cache

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IPampurin/UrlShortener/pkg/configuration"
	"github.com/wb-go/wbf/retry"
)

// LoadDataToCache загружает данные за последнее время в кэш при старте
func LoadDataToCache(ctx context.Context, cfg *configuration.ConfCache) error {

	// получаем заказы до установленного порога
	notifications, err := db.GetClientDB().GetNotificationsLastPeriod(ctx, cfg.Warming)
	if err != nil {
		return fmt.Errorf("ошибка получения уведомлений из БД при прогреве кэша: %v", err)
	}
	// проверяем были ли уведомления в БД
	if len(notifications) == 0 {
		return nil
	}

	// определяем стратегию ретраев
	strategy := retry.Strategy{Attempts: 3, Delay: 100 * time.Millisecond, Backoff: 2}

	// сохраняем данные в redis
	for i := range notifications {
		key := notifications[i].UID.String()
		err := GetClientRedis().SetWithExpirationAndRetry(ctx, strategy, key, notifications[i], cfg.TTL)
		if err != nil {
			log.Printf("ошибка добавления уведомления %s при прогреве кэша", key)
			continue
		}
	}

	return nil
}
