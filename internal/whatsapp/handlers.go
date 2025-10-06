package whatsapp

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
	"wa-mcp-bridge/internal/store"

	"github.com/go-chi/chi/v5"
	"go.mau.fi/whatsmeow"
)

type Bot struct {
	store    store.Store
	client *whatsmeow.Client

	qrChan chan string
	startOnce sync.Once
}

const pingInterval = time.Second * 10

func New(store store.Store) (*Bot, error) { 
	return &Bot{
		store: store,
		qrChan: make(chan string, 1),
	}, nil
}

func (b *Bot) RegisterRoutes(r chi.Router) {
	r.Get("/status", b.handleStatus())
	r.Get("/login", b.handleLogin())
	r.Get("/qr", b.handleQR())
}

// Handlers
func (b *Bot) handleStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if b.client != nil && b.client.IsConnected() {
			w.Write([]byte("Connected"))
		}

		w.Write([]byte("Disconnected"))
	}
}

func (b *Bot) handleLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := b.Login(); err != nil {
			log.Printf("ERROR: failed to connect to whatsapp: %v", err)
			http.Error(w, "Failed to connect to whatsapp", http.StatusInternalServerError)
		}

		w.Write([]byte("Whatsapp Connection initiated."))
	}
}



func (b *Bot) handleQR() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if b.client != nil && b.client.IsConnected() {
			log.Printf("Bot already running.")
			http.Error(w, "Bot already connected to a whatsapp account.", http.StatusConflict)
			return 
		} 
		
		if b.client == nil {
				if err := b.Login(); err != nil {
					http.Error(w, "Failed to start a connection to whatsapp", http.StatusInternalServerError)
					return 
				}
		}


		ctx := r.Context()

		hdr := w.Header()
		hdr.Set("Content-Type", "text/event-stream")
		hdr.Set("Cache-Control", "no-cache")
		hdr.Set("Connection", "keep-alive")
		hdr.Set("X-Accel-Buffering", "no")

		f, ok := w.(http.Flusher)
		if !ok {
			return
		}	
		
		io.WriteString(w, ": ping\n\n")
		f.Flush()
			
		L:
		for {
			select {
			case <- ctx.Done():
				break L
			case <-time.After(pingInterval):
				io.WriteString(w, ": heart-beat\n\n")
				f.Flush()
			case qrCode, ok := <-b.qrChan:
				if !ok {
					io.WriteString(w, "event: error\ndata: The QR code process has timed out or completed.\n\n")
					f.Flush()
					break L
				}

				io.WriteString(w, "event: qr\n")
				fmt.Fprintf(w, "data: %s\n\n", qrCode)
			}
		}
	}
}


























































// For real my mind saying learn rust wtf.
// What am i even doing 
// Bitch as mind or me ? 
// Lol. i need to learn english first