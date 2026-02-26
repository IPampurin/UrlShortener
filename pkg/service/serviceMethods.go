package service

import (
	"context"
	"fmt"
	"time"

	"github.com/IPampurin/UrlShortener/pkg/db"
)

// CreateLink создаёт новую короткую ссылку
// если customShort не пуст, проверяет его уникальность
// если originalURL уже существует, возвращает существующую ссылку
// генерирует случайный shortURL при необходимости
// сохраняет ссылку в БД и, если кэш доступен, записывает её туда
func (s *Service) CreateLink(ctx context.Context, originalURL, customShort string) (*db.Link, error) {

	// 1. Если задан кастомный short, проверяем уникальность
	if customShort != "" {
		shortURL, err := s.link.GetLinkByShortURL(ctx, customShort)
		if err != nil {
			return nil, fmt.Errorf("ошибка проверки кастомной shortURL: %w", err)
		}
		if shortURL != nil {
			return nil, fmt.Errorf("короткая ссылка уже занята")
		}
	}

	// 2. Проверяем, существует ли уже такая оригинальная ссылка (она может быть не одна)
	links, err := s.link.GetLinkByOriginalURL(ctx, originalURL)
	if err != nil {
		return nil, fmt.Errorf("ошибка поиска по оригинальному URL: %w", err)
	}

	if len(links) > 0 {
		// Возвращаем первую найденную (можно выбрать самую свежую)
		// Если есть кэш, можно положить её туда для ускорения будущих запросов
		if s.cache != nil {
			_ = s.cache.SetLink(ctx, links[0].ShortURL, links[0], 0) // TBD: взять из конфига
		}
		return links[0], nil
	}

	// 3. Генерируем shortURL, если не задан
	shortURL := customShort
	if shortURL == "" {
		// Генерируем, пока не найдём уникальный
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
		return nil, fmt.Errorf("ошибка сохранения ссылки: %w", err)
	}

	// 5. Сохраняем в кэш, если он доступен
	if s.cache != nil {
		_ = s.cache.SetLink(ctx, shortURL, link, 0) // TTL можно задать из конфига
	}

	return link, nil
}

// GetLinkByShortURL возвращает ссылку по короткому URL
func (s *Service) GetLinkByShortURL(ctx context.Context, shortURL string) (*db.Link, error) {

	// проверяем кэш
	if s.cache != nil {
		link, err := s.cache.GetLink(ctx, shortURL)
		if err != nil {
			// логируем ошибку кэша, но не прерываем
			// (как добавить логгер в сервис)
		}
		if link != nil {
			return link, nil
		}
	}

	// идём в БД
	link, err := s.link.GetLinkByShortURL(ctx, shortURL)
	if err != nil {
		return nil, err
	}
	if link == nil {
		return nil, nil // не найдено
	}

	// сохраняем в кэш для будущих запросов
	if s.cache != nil {
		_ = s.cache.SetLink(ctx, shortURL, link, 0)
	}

	return link, nil
}

// GetLinkByOriginalURL возвращает все ссылки по оригинальному URL
// (кэширование здесь не применяем пока что)
func (s *Service) GetLinkByOriginalURL(ctx context.Context, originalURL string) ([]*db.Link, error) {

	return s.link.GetLinkByOriginalURL(ctx, originalURL)
}

// IncrementClicks увеличивает счётчик переходов по ссылке (вызывается после успешного редиректа)
func (s *Service) IncrementClicks(ctx context.Context, linkID int64) error {

	return s.link.IncrementClicks(ctx, linkID)
}

// GetLinks возвращает последние 20 ссылок (для UI)
func (s *Service) GetLinks(ctx context.Context) ([]*db.Link, error) {

	return s.link.GetLinks(ctx)
}

// SaveAnalytics сохраняет информацию о переходе
func (s *Service) SaveAnalytics(ctx context.Context, analytics *db.Analytics) error {

	return s.analytics.SaveAnalytics(ctx, analytics)
}

// GetAnalyticsByLinkID возвращает все переходы по ссылке
func (s *Service) GetAnalyticsByLinkID(ctx context.Context, linkID int) ([]*db.Analytics, error) {

	return s.analytics.GetAnalyticsByLinkID(ctx, linkID)
}

// CountClicksByDay возвращает статистику переходов по дням
func (s *Service) CountClicksByDay(ctx context.Context, linkID int, from, to time.Time) (map[string]int, error) {

	return s.analytics.CountClicksByDay(ctx, linkID, from, to)
}

// CountClicksByMonth возвращает статистику переходов по месяцам
func (s *Service) CountClicksByMonth(ctx context.Context, linkID int, from, to time.Time) (map[string]int, error) {

	return s.analytics.CountClicksByMonth(ctx, linkID, from, to)
}

// CountClicksByUserAgent возвращает статистику переходов по User-Agent
func (s *Service) CountClicksByUserAgent(ctx context.Context, linkID int) (map[string]int, error) {

	return s.analytics.CountClicksByUserAgent(ctx, linkID)
}
