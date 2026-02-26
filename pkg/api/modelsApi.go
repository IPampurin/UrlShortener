package api

import "time"

// CreateRequest - запрос на создание короткой ссылки
type CreateRequest struct {
	OriginalURL string `json:"original_url" binding:"required,url"`
	CustomShort string `json:"custom_short" binding:"omitempty,alphanum,max=50"`
}

// FullResponseLink - ответ на успешное создание или запрос данных
type FullResponseLink struct {
	ShortURL    string    `json:"short_url"`
	OriginalURL string    `json:"original_url"`
	CreatedAt   time.Time `json:"created_at"`
	ClicksCount int       `json:"clicks_count"`
}

// ClickEntry - информация об одном переходе (для аналитики)
type ClickEntry struct {
	AccessedAt time.Time `json:"accessed_at"`
	UserAgent  string    `json:"user_agent"`
	IPAddress  string    `json:"ip_address,omitempty"`
	Referer    string    `json:"referer,omitempty"`
}

// AnalyticsResponse - полный ответ для GET /analytics/:short_url
type AnalyticsResponse struct {
	ShortURL    string       `json:"short_url"`
	OriginalURL string       `json:"original_url"`
	CreatedAt   time.Time    `json:"created_at"`
	ClicksCount int          `json:"clicks_count"`
	Clicks      []ClickEntry `json:"clicks"`
}

// ErrorResponse - стандартный ответ с ошибкой
type ErrorResponse struct {
	Error string `json:"error"`
}
