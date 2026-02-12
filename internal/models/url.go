package models

import "time"

type URL struct {
    ID          int       `json:"id"`
    OriginalURL string    `json:"original_url"`
    ShortCode   string    `json:"short_code"`
    Clicks      int       `json:"clicks"`
    CreatedAt   time.Time `json:"created_at"`
    ExpiresAt   *time.Time `json:"expires_at,omitempty"`
}

type CreateURLRequest struct {
    URL       string `json:"url"`
    CustomCode string `json:"custom_code,omitempty"`
}

type CreateURLResponse struct {
    ShortURL  string `json:"short_url"`
    ShortCode string `json:"short_code"`
}

type Stats struct {
    TotalURLs  int `json:"total_urls"`
    TotalClicks int `json:"total_clicks"`
}