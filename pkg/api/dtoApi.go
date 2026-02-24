package api

import "time"

// CreateShortLinkRequest - тело запроса на создание короткой ссылки
type CreateShortLinkRequest struct {
	OriginalURL string `json:"original_url" binding:"required,url"`
	CustomShort string `json:"custom_short" binding:"omitempty,alphanum,max=50"`
}

// LinkApi - ответ на успешное создание
type LinkApi struct {
	ShortURL    string    `json:"short_url"`
	OriginalURL string    `json:"original_url"`
	CreatedAt   time.Time `json:"created_at"`
	ClicksCount int       `json:"clicks_count"` // всегда 0 при создании
}

// ErrorResponse - стандартный ответ с ошибкой
type ErrorResponse struct {
	Error string `json:"error"`
}
