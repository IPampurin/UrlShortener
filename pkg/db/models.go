package db

import (
	"time"
)

// Link представляет запись в таблице links
type Link struct {
	ID          int64     `db:"id"`           // ID — внутренний идентификатор ссылки (автоинкремент)
	ShortURL    string    `db:"short_url"`    // ShortURL — короткий идентификатор (например, "abc123"), уникален в пределах таблицы
	OriginalURL string    `db:"original_url"` // OriginalURL — исходный длинный URL
	CreatedAt   time.Time `db:"created_at"`   // CreatedAt — дата и время создания записи
	IsCustom    bool      `db:"is_custom"`    // IsCustom — флаг, указывающий, что short_url задан пользователем, а не сгенерирован системой
	ClicksCount int32     `db:"clicks_count"` // ClicksCount — количество переходов по ссылке (денормализованное поле)
}

// Analytics представляет запись о переходе по короткой ссылке
type Analytics struct {
	ID         int64     `db:"id"`          // ID — уникальный идентификатор записи о переходе
	LinkID     int64     `db:"link_id"`     // LinkID — идентификатор ссылки, по которой совершён переход
	AccessedAt time.Time `db:"accessed_at"` // AccessedAt — момент времени, когда произошёл переход
	UserAgent  string    `db:"user_agent"`  // UserAgent — строка User-Agent браузера или клиента
	IPAddress  *string   `db:"ip_address"`  // IPAddress — IP-адрес посетителя (может быть nil, если не сохраняется)
	Referer    *string   `db:"referer"`     // Referer — URL источника перехода (может быть nil)
}
