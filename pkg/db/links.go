package db

import (
	"context"
	"fmt"
)

// CreateLink добавляет запись в таблицу links БД
func (c *ClientPostgres) CreateLink(ctx context.Context, shortURL, originalURL string, isCustom bool) (*LinkDB, error) {

	query := `   INSERT INTO links (short_url, original_url, created_at, is_custom, clicks_count)
                 VALUES ($1, $2, NOW(), $3, $4)
              RETURNING id, created_at`

	linkDB := &LinkDB{}

	err := c.Pool.QueryRow(ctx, query, shortURL, originalURL, isCustom, 0).
		Scan(&linkDB.ID,
			&linkDB.ShortURL,
			&linkDB.OriginalURL,
			&linkDB.CreatedAt,
			&linkDB.IsCustom,
			&linkDB.ClicksCount)
	if err != nil {
		return nil, fmt.Errorf("ошибка добавления записи о ссылке в CreateLink: %w", err)
	}

	return linkDB, nil
}

// GetLinkByShortURL получает из таблицы links БД запись по короткой ссылке
func (c *ClientPostgres) GetLinkByShortURL(ctx context.Context, shortURL string) (*LinkDB, error) {

	query := `SELECT * 
	            FROM links 
			   WHERE short_url = $1`

	linkDB := &LinkDB{}

	err := c.Pool.QueryRow(ctx, query, shortURL).
		Scan(&linkDB.ID,
			&linkDB.ShortURL,
			&linkDB.OriginalURL,
			&linkDB.CreatedAt,
			&linkDB.IsCustom,
			&linkDB.ClicksCount)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения записи о ссылке в GetLinkByShortURL: %w", err)
	}

	return linkDB, nil
}

// IncrementClicks увеличивает счётчик переходов по ссылке
func (c *ClientPostgres) IncrementClicks(ctx context.Context, linkID int64) error {

	query := `UPDATE links
	             SET clicks_count = clicks_count + 1
			   WHERE id = $1`

	_, err := c.Pool.Exec(ctx, query, linkID)
	if err != nil {
		return fmt.Errorf("ошибка увеличения счётчика переходов в IncrementClicks: %w", err)
	}

	return nil
}

// LastLinks получает крайние по времени 20 записей по сокращению ссылок
func (c *ClientPostgres) LastLinks(ctx context.Context) ([]*LinkDB, error) {

	const limit = 20

	query := `SELECT *
	            FROM links
			   LIMIT $1`

	rows, err := c.Pool.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении списка ссылок в LastLinks: %w", err)
	}
	defer rows.Close()

	var links []*LinkDB
	for rows.Next() {
		var l LinkDB
		err := rows.Scan(
			&l.ID,
			&l.ShortURL,
			&l.OriginalURL,
			&l.CreatedAt,
			&l.IsCustom,
			&l.ClicksCount,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка при сканировании строки списка ссылок в LastLinks: %w", err)
		}

		links = append(links, &l)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при итерации по списку ссылок в LastLinks: %w", err)
	}

	return links, nil
}
