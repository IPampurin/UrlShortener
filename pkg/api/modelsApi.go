package api

import "time"

// CreateRequest - тело запроса на создание короткой ссылки
type CreateRequest struct {
	OriginalURL string `json:"original_url" binding:"required,url"`
	CustomShort string `json:"custom_short" binding:"omitempty,alphanum,max=50"`
}

// FullResponseLink - ответ на успешное создание или запрос данных
type FullResponseLink struct {
	ShortURL    string    `json:"short_url"`
	OriginalURL string    `json:"original_url"`
	CreatedAt   time.Time `json:"created_at"`
	ClicksCount int       `json:"clicks_count"` // всегда 0 при создании
}

type FullResponseAnalitics struct {
}

// ErrorResponse - стандартный ответ с ошибкой
type ErrorResponse struct {
	Error string `json:"error"`
}
