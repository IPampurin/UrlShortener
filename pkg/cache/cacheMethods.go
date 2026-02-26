package cache

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/IPampurin/UrlShortener/pkg/db"
	"github.com/wb-go/wbf/redis"
	"github.com/wb-go/wbf/retry"
)

// LoadDataToCache загружает данные за последнее время в кэш при старте
func (c *Cache) LoadDataToCache(ctx context.Context, lastLinks []*db.Link, ttl time.Duration) error {

	// определяем стратегию ретраев
	strategy := retry.Strategy{Attempts: 3, Delay: 100 * time.Millisecond, Backoff: 2}

	// сохраняем данные в redis
	for _, link := range lastLinks {
		key := link.ShortURL
		err := c.SetWithExpirationAndRetry(ctx, strategy, key, link, ttl)
		if err != nil {
			log.Printf("ошибка добавления ссылки %s при прогреве кэша: %v", key, err)
			continue
		}
	}

	return nil
}

// GetLink возвращает ссылку из кэша по короткому URL (или nil, nil)
func (c *Cache) GetLink(ctx context.Context, shortURL string) (*db.Link, error) {

	data, err := c.Get(ctx, shortURL)
	if err != nil {
		if errors.Is(err, redis.NoMatches) {
			return nil, nil // не найдено
		}
		return nil, err
	}

	var link db.Link
	if err := json.Unmarshal([]byte(data), &link); err != nil {
		return nil, err
	}

	return &link, nil
}

// SetLink сохраняет ссылку в кэш с указанным TTL
func (c *Cache) SetLink(ctx context.Context, shortURL string, link *db.Link, ttl time.Duration) error {

	data, err := json.Marshal(link)
	if err != nil {
		return err
	}

	return c.SetWithExpiration(ctx, shortURL, data, ttl)
}

// DeleteLink удаляет ссылку из кэша
func (c *Cache) DeleteLink(ctx context.Context, shortURL string) error {

	return c.Del(ctx, shortURL)
}
