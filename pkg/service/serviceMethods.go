package service

import (
	"context"
	"time"

	"github.com/IPampurin/UrlShortener/pkg/db"
)

// CreateLink создаёт новую ссылку или возвращает существующую.
func (s *Service) CreateLink(ctx context.Context, originalURL, customShort string) (*db.Link, error) {

	// 1. Если задан кастомный short, проверяем уникальность
	if customShort != "" {
		existing, _ := s.link.GetLinkByShortURL(ctx, customShort)
		if existing != nil {
			return nil, ErrShortURLExists
		}
	}

	// 2. Проверяем, есть ли уже такая оригинальная ссылка (опционально)
	links, err := s.link.GetLinkByOriginalURL(ctx, originalURL)
	if err != nil {
		return nil, err
	}
	if len(links) > 0 {
		// возвращаем первую найденную
		return links[0], nil
	}

	// 3. Генерируем short, если не задан
	shortURL := customShort
	if shortURL == "" {
		// генерируем, пока не найдём уникальный
		for {
			shortURL = NewRandomString(0)
			existing, _ := s.link.GetLinkByShortURL(ctx, shortURL)
			if existing == nil {
				break
			}
		}
	}

	// 4. Создаём объект и сохраняем
	link := &db.Link{
		ShortURL:    shortURL,
		OriginalURL: originalURL,
		IsCustom:    customShort != "",
		ClicksCount: 0,
	}
	if err := s.link.SaveLink(ctx, link); err != nil {
		return nil, err
	}
	return link, nil
}

// GetLink возвращает ссылку по short URL.
func (s *Service) GetLink(ctx context.Context, shortURL string) (*db.Link, error) {
	link, err := s.link.GetLinkByShortURL(ctx, shortURL)
	if err != nil {
		return nil, err
	}
	if link == nil {
		return nil, ErrLinkNotFound
	}
	return link, nil
}

// RecordClick записывает аналитику перехода и увеличивает счётчик.
func (s *Service) RecordClick(ctx context.Context, shortURL, userAgent, ip, referer string) error {
	link, err := s.link.GetLinkByShortURL(ctx, shortURL)
	if err != nil {
		return err
	}
	if link == nil {
		return ErrLinkNotFound
	}
	analytics := &db.Analytics{
		LinkID:     link.ID,
		AccessedAt: time.Now(),
		UserAgent:  userAgent,
		IPAddress:  ip,
		Referer:    referer,
	}
	if err := s.analytics.SaveAnalytics(ctx, analytics); err != nil {
		return err
	}
	return s.link.IncrementClicks(ctx, int64(link.ID))
}

// GetAnalytics возвращает ссылку и все переходы.
func (s *Service) GetAnalytics(ctx context.Context, shortURL string) (*db.Link, []*db.Analytics, error) {
	link, err := s.link.GetLinkByShortURL(ctx, shortURL)
	if err != nil {
		return nil, nil, err
	}
	if link == nil {
		return nil, nil, ErrLinkNotFound
	}
	analytics, err := s.analytics.GetAnalyticsByLinkID(ctx, link.ID)
	if err != nil {
		return nil, nil, err
	}
	return link, analytics, nil
}

// LastLinks возвращает последние 20 ссылок.
func (s *Service) LastLinks(ctx context.Context) ([]*db.Link, error) {
	return s.link.LastLinks(ctx)
}
