package db

import (
	"context"
	"fmt"
)

// SaveLink добавляет запись в таблицу links БД
func (d *DataBase) SaveLink(ctx context.Context, link *Link) error {

	query := `   INSERT INTO links (short_url, original_url, created_at, is_custom, clicks_count)
                 VALUES ($1, $2, NOW(), $3, $4)
			  RETURNING id, created_at`

	err := d.Postgres.Pool.QueryRow(ctx, query, link.ShortURL, link.OriginalURL, link.IsCustom, link.ClicksCount).
		Scan(&link.ID, &link.CreatedAt)
	if err != nil {
		return fmt.Errorf("ошибка добавления записи о ссылке в SaveLink: %w", err)
	}

	return nil
}

// GetLinkByShortURL получает из таблицы links БД запись по короткой ссылке
func (d *DataBase) GetLinkByShortURL(ctx context.Context, shortURL string) (*Link, error) {

	query := `SELECT * 
	            FROM links 
			   WHERE short_url = $1`

	linkDB := &Link{}

	err := d.Pool.QueryRow(ctx, query, shortURL).
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

// GetLinkByOriginalURL получает из таблицы links БД запись по длинной ссылке
func (d *DataBase) GetLinkByOriginalURL(ctx context.Context, originalURL string) ([]*Link, error) {

	query := `SELECT * 
	            FROM links 
			   WHERE original_url = $1`

	rows, err := d.Pool.Query(ctx, query, originalURL)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении списка ссылок в GetLinkByOriginalURL: %w", err)
	}
	defer rows.Close()

	links := make([]*Link, 0)
	for rows.Next() {
		var l Link
		err := rows.Scan(
			&l.ID,
			&l.ShortURL,
			&l.OriginalURL,
			&l.CreatedAt,
			&l.IsCustom,
			&l.ClicksCount,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка при сканировании строки списка ссылок в GetLinkByOriginalURL: %w", err)
		}

		links = append(links, &l)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при итерации по списку ссылок в GetLinkByOriginalURL: %w", err)
	}

	return links, nil
}

// IncrementClicks увеличивает счётчик переходов по ссылке
func (d *DataBase) IncrementClicks(ctx context.Context, linkID int64) error {

	query := `UPDATE links
	             SET clicks_count = clicks_count + 1
			   WHERE id = $1`

	_, err := d.Pool.Exec(ctx, query, linkID)
	if err != nil {
		return fmt.Errorf("ошибка увеличения счётчика переходов в IncrementClicks: %w", err)
	}

	return nil
}

// GetLinks получает крайние по времени 20 записей по сокращению ссылок
func (d *DataBase) GetLinks(ctx context.Context) ([]*Link, error) {

	const limit = 20

	query := `SELECT *
	            FROM links
			   LIMIT $1`

	rows, err := d.Pool.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении списка ссылок в GetLinks: %w", err)
	}
	defer rows.Close()

	links := make([]*Link, 0)
	for rows.Next() {
		var l Link
		err := rows.Scan(
			&l.ID,
			&l.ShortURL,
			&l.OriginalURL,
			&l.CreatedAt,
			&l.IsCustom,
			&l.ClicksCount,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка при сканировании строки списка ссылок в GetLinks: %w", err)
		}

		links = append(links, &l)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при итерации по списку ссылок в GetLinks: %w", err)
	}

	return links, nil
}
