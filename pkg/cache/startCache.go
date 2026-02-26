package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/IPampurin/UrlShortener/pkg/configuration"
	"github.com/wb-go/wbf/logger"
	"github.com/wb-go/wbf/redis"
)

/*
в кэше должны хранится пары ShortURL - OriginalURL
также должны быть сочетания OriginalURL - ShortURL // при чём ShortURL в этом случае только сгенерированные! (кастомные ShortURL надо проверять первыми)
*/

// Cache хранит подключение к БД Redis
type Cache struct {
	*redis.Client
}

// InitCache запускает работу с Redis
func InitCache(ctx context.Context, cfgCache *configuration.ConfCache, log logger.Logger) (*Cache, error) {

	// определяем конфигурацию подключения к Redis
	options := redis.Options{
		Address:   fmt.Sprintf("%s:%d", cfgCache.HostName, cfgCache.Port),
		Password:  cfgCache.Password,
		MaxMemory: "100mb",
		Policy:    "allkeys-lru",
	}

	// пробуем подключиться
	clientRedis, err := redis.Connect(options)
	if err != nil {
		return nil, fmt.Errorf("ошибка установки соединения с Redis: %v\n", err)
	}

	// проверяем подключение
	if err := clientRedis.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("ошибка подключения к Redis: %v\n", err)
	}

	// получаем экземпляр
	cache := &Cache{clientRedis}

	// загружаем начальные данные
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	/*
		err = loadDataToCache(ctx, cfgCache)
		if err != nil {
			log.Error("ошибка загрузки первичных данных в кэш: %v", err)
			return nil, err
		}
	*/
	log.Info("Кэш прогрет и работает.")

	return cache, nil
}

/*
// loadDataToCache загружает данные за последнее время в кэш при старте
func loadDataToCache(ctx context.Context, cfg *configuration.ConfCache) error {

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
*/
