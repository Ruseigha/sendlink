package handlers

import (
	"encoding/json"
	"math/rand"
	"net/http"

	"github.com/Ruseigha/sendlink/internal/cache"
	"github.com/Ruseigha/sendlink/internal/config"
	"github.com/Ruseigha/sendlink/internal/database"
	"github.com/Ruseigha/sendlink/internal/models"
	"github.com/gorilla/mux"
)

type Handler struct {
    db     *database.PostgresDB
    cache  *cache.RedisCache
    config *config.Config
}

func NewHandler(db *database.PostgresDB, cache *cache.RedisCache, cfg *config.Config) *Handler {
    return &Handler{db: db, cache: cache, config: cfg}
}

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
    response := map[string]string{
        "status":      "healthy",
        "environment": h.config.Environment,
    }
    json.NewEncoder(w).Encode(response)
}

func (h *Handler) CreateShortURL(w http.ResponseWriter, r *http.Request) {
    var req models.CreateURLRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    shortCode := req.CustomCode
    if shortCode == "" {
        shortCode = h.generateShortCode()
    }

    url := &models.URL{
        OriginalURL: req.URL,
        ShortCode:   shortCode,
    }

    if err := h.db.CreateURL(url); err != nil {
        http.Error(w, "Failed to create URL", http.StatusInternalServerError)
        return
    }

    response := models.CreateURLResponse{
        ShortURL:  h.config.BaseURL + "/" + shortCode,
        ShortCode: shortCode,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func (h *Handler) RedirectURL(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    shortCode := vars["code"]

    // Try cache first
    if originalURL, err := h.cache.Get(r.Context(), shortCode); err == nil {
        h.db.IncrementClicks(shortCode)
        http.Redirect(w, r, originalURL, http.StatusFound)
        return
    }

    // Fallback to database
    url, err := h.db.GetURLByShortCode(shortCode)
    if err != nil {
        http.Error(w, "URL not found", http.StatusNotFound)
        return
    }

    // Cache for next time
    h.cache.Set(r.Context(), shortCode, url.OriginalURL, 24*60*60*1000000000)

    h.db.IncrementClicks(shortCode)
    http.Redirect(w, r, url.OriginalURL, http.StatusFound)
}

func (h *Handler) GetStats(w http.ResponseWriter, r *http.Request) {
    stats, err := h.db.GetStats()
    if err != nil {
        http.Error(w, "Failed to get stats", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(stats)
}

func (h *Handler) generateShortCode() string {
    const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    code := make([]byte, h.config.ShortURLLength)
    for i := range code {
        code[i] = charset[rand.Intn(len(charset))]
    }
    return string(code)
}