package interfaces

import (
	"context"
	"time"

	"github.com/IPampurin/UrlShortener/pkg/db"
)

// методы по таблице Link
type LinkMethods interface {
	SaveLink(ctx context.Context, link *db.Link) error
	GetLinkByShortURL(ctx context.Context, shortURL string) (*db.Link, error)
	GetLinkByOriginalURL(ctx context.Context, originalURL string) ([]*db.Link, error)
	IncrementClicks(ctx context.Context, linkID int64) error
	GetLinks(ctx context.Context) ([]*db.Link, error)
}

// методы по таблице Analytics
type AnalyticsMethods interface {
	SaveAnalytics(ctx context.Context, analytics *db.Analytics) error
	GetAnalyticsByLinkID(ctx context.Context, linkID int) ([]*db.Analytics, error)
	CountClicksByDay(ctx context.Context, linkID int, from, to time.Time) (map[string]int, error)
	CountClicksByMonth(ctx context.Context, linkID int, from, to time.Time) (map[string]int, error)
	CountClicksByUserAgent(ctx context.Context, linkID int) (map[string]int, error)
}

// методы кэша
type CacheMetods interface {
	Get(ctx context.Context, key string) (*db.Link, error)
	Set(ctx context.Context, key string, link *db.Link, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	LoadDataToCache(ctx context.Context, warming time.Duration, ttl time.Duration) error
}
