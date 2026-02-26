package interfaces

import (
	"context"
	"time"

	"github.com/IPampurin/UrlShortener/pkg/db"
)

// методы по таблице Link
type LinkMethods interface {
	CreateLink(ctx context.Context, originalURL, shortURL string, isCustom bool) (*db.Link, error)
	GetLinkByShortURL(ctx context.Context, shortURL string) (*db.Link, error)
	GetLinkByOriginalURL(ctx context.Context, originalURL string) ([]*db.Link, error)
	IncrementClicks(ctx context.Context, linkID int64) error
	GetLinks(ctx context.Context) ([]*db.Link, error)
	GetLinksOfPeriod(ctx context.Context, period time.Duration) ([]*db.Link, error)
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
type CacheMethods interface {
	GetLink(ctx context.Context, shortURL string) (*db.Link, error)
	SetLink(ctx context.Context, shortURL string, link *db.Link, ttl time.Duration) error
	DeleteLink(ctx context.Context, shortURL string) error
	LoadDataToCache(ctx context.Context, lastLiks []*db.Link, ttl time.Duration) error
}
