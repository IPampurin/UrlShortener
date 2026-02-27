package service

import (
	"context"
	"errors"
	"time"

	"github.com/IPampurin/UrlShortener/pkg/db"
	"github.com/wb-go/wbf/logger"
)

// CreateShortLink создаёт новую короткую ссылку.

func (s *Service) CreateShortLink(ctx context.Context, log logger.Logger, originalURL, customShort string) (*ResponseLink, error) {

	// 1. Если задан кастомный short, проверяем уникальность
	if customShort != "" {
		existing, err := s.link.GetLinkByShortURL(ctx, customShort)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			return nil, errors.New("короткая ссылка уже занята")
		}
	}

	// 2. Проверяем, есть ли уже такая оригинальная ссылка
	links, err := s.link.GetLinkByOriginalURL(ctx, originalURL)
	if err != nil {
		return nil, err
	}
	if len(links) > 0 {
		latest := links[0]
		for _, l := range links {
			if l.CreatedAt.After(latest.CreatedAt) {
				latest = l
			}
		}
		if s.cache != nil {
			if err := s.cache.SetLink(ctx, latest.ShortURL, latest); err != nil {
				log.Ctx(ctx).Error("ошибка сохранения в кэш", "error", err)
			}
		}
		return toResponseLink(latest), nil
	}

	// 3. Генерируем shortURL, если не задан
	shortURL := customShort
	if shortURL == "" {
		for {
			shortURL = NewRandomString(0)
			existing, _ := s.link.GetLinkByShortURL(ctx, shortURL)
			if existing == nil {
				break
			}
		}
	}

	// 4. Создаём новую ссылку
	link, err := s.link.CreateLink(ctx, originalURL, shortURL, customShort != "")
	if err != nil {
		return nil, err
	}

	// 5. Сохраняем в кэш
	if s.cache != nil {
		if err := s.cache.SetLink(ctx, shortURL, link); err != nil {
			log.Ctx(ctx).Error("ошибка сохранения в кэш", "error", err)
		}
	}

	return toResponseLink(link), nil
}

// ShortLinkInfo возвращает информацию о ссылке по короткому идентификатору
func (s *Service) ShortLinkInfo(ctx context.Context, log logger.Logger, shortURL string) (*ResponseLink, error) {

	if s.cache != nil {
		link, err := s.cache.GetLink(ctx, shortURL)
		if err != nil {
			log.Ctx(ctx).Error("ошибка получения из кэша", "error", err)
		}
		if link != nil {
			return toResponseLink(link), nil
		}
	}

	link, err := s.link.GetLinkByShortURL(ctx, shortURL)
	if err != nil {
		return nil, err
	}
	if link == nil {
		return nil, nil
	}

	if s.cache != nil {
		if err := s.cache.SetLink(ctx, shortURL, link); err != nil {
			log.Ctx(ctx).Error("ошибка сохранения в кэш", "error", err)
		}
	}

	return toResponseLink(link), nil
}

// ShortLinkAnalytics возвращает аналитику по ссылке
func (s *Service) ShortLinkAnalytics(ctx context.Context, log logger.Logger, shortURL string) (*ResponseAnalytics, error) {

	link, err := s.link.GetLinkByShortURL(ctx, shortURL)
	if err != nil {
		return nil, err
	}
	if link == nil {
		return nil, nil
	}

	analytics, err := s.analytics.GetAnalyticsByLinkID(ctx, link.ID)
	if err != nil {
		return nil, err
	}

	// получаем агрегаты
	from := time.Now().AddDate(0, -1, 0) // например, последний месяц для демонстрации
	to := time.Now()

	clicksByDay, err := s.analytics.CountClicksByDay(ctx, link.ID, from, to)
	if err != nil {
		log.Ctx(ctx).Error("ошибка агрегации по дням", "error", err)
		// не фатально, можно оставить пустым
	}

	clicksByMonth, err := s.analytics.CountClicksByMonth(ctx, link.ID, from, to)
	if err != nil {
		log.Ctx(ctx).Error("ошибка агрегации по месяцам", "error", err)
	}

	clicksByUA, err := s.analytics.CountClicksByUserAgent(ctx, link.ID)
	if err != nil {
		log.Ctx(ctx).Error("ошибка агрегации по user-agent", "error", err)
	}

	followLinks := make([]FollowLink, len(analytics))
	for i, a := range analytics {
		followLinks[i] = FollowLink{
			AccessedAt: a.AccessedAt,
			UserAgent:  a.UserAgent,
			IPAddress:  a.IPAddress.String(),
			Referer:    a.Referer,
		}
	}

	return &ResponseAnalytics{
		Link:              *toResponseLink(link),
		Analytics:         followLinks,
		ClicksByDay:       clicksByDay,
		ClicksByMonth:     clicksByMonth,
		ClicksByUserAgent: clicksByUA,
	}, nil
}

// LastLinks возвращает последние ссылки
func (s *Service) LastLinks(ctx context.Context, log logger.Logger) ([]*ResponseLink, error) {

	links, err := s.link.GetLinks(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*ResponseLink, len(links))
	for i, l := range links {
		result[i] = toResponseLink(l)
	}

	return result, nil
}

// RecordClick сохраняет информацию о переходе
func (s *Service) RecordClick(ctx context.Context, log logger.Logger, linkID int, userAgent, ip, referer string) error {

	return s.analytics.SaveAnalytics(ctx, linkID, time.Now(), userAgent, ip, referer)
}

// IncrementClicks увеличивает счётчик переходов
func (s *Service) IncrementClicks(ctx context.Context, log logger.Logger, linkID int64) error {

	return s.link.IncrementClicks(ctx, linkID)
}

// SearchByOriginalURL ищет среди OriginalURL содержащие query
func (s *Service) SearchByOriginalURL(ctx context.Context, log logger.Logger, query string) ([]*ResponseLink, error) {

	links, err := s.link.SearchByOriginalURL(ctx, query)
	if err != nil {
		return nil, err
	}

	result := make([]*ResponseLink, len(links))
	for i, l := range links {
		result[i] = toResponseLink(l)
	}

	return result, nil
}

// SearchByShortURL ищет среди ShortURL содержащие query
func (s *Service) SearchByShortURL(ctx context.Context, log logger.Logger, query string) ([]*ResponseLink, error) {

	links, err := s.link.SearchByShortURL(ctx, query)
	if err != nil {
		return nil, err
	}

	result := make([]*ResponseLink, len(links))
	for i, l := range links {
		result[i] = toResponseLink(l)
	}

	return result, nil
}

// toResponseLink преобразует db.Link в service.ResponseLink
func toResponseLink(l *db.Link) *ResponseLink {

	return &ResponseLink{
		ID:          l.ID,
		ShortURL:    l.ShortURL,
		OriginalURL: l.OriginalURL,
		CreatedAt:   l.CreatedAt,
		ClicksCount: l.ClicksCount,
	}
}
