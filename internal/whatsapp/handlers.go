package whatsapp

import (
	"net/http"
	"wa-mcp-bridge/internal/store"

	"github.com/go-chi/chi/v5"
)

type Bot struct {
	store store.Store
}

func New(store store.Store) (*Bot, error) {
	return &Bot{
		store: store,
	}, nil
}

func (b *Bot) RegisterRoutes(r chi.Router) {
	r.Get("/start", b.handleStartWAInstance())
}

// Handlers
func (b *Bot) handleStartWAInstance() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Starting wa instance"))
	}
}
