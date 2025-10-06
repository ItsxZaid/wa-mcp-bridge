package whatsapp

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"

	"github.com/skip2/go-qrcode"

	waLog "go.mau.fi/whatsmeow/util/log"
)

func (b *Bot) Login() error {
	var initError error

	b.startOnce.Do(func () {
		if b.client != nil && b.client.IsConnected() {
		return
	}

	dbLog := waLog.Stdout("Database", "DEBUG", true)

	ctx := context.Background()
	container, err := sqlstore.New(ctx, "sqlite3", "file:wa-mcp-bridge .db?_foreign_keys=on", dbLog)
		if err != nil {
			initError = fmt.Errorf("[WHATSAPP] Failed to initialize SQLite store %w", err)
			return 
		}

	deviceStore, err :=  container.GetFirstDevice(ctx)
	if err != nil {
		initError =  fmt.Errorf("[WHATSAPP] Failed to grab first device data %w", err)
		return 
	}

	waClientLog := waLog.Stdout("WAClient", "DEBUG", true)
	b.client = whatsmeow.NewClient(deviceStore, waClientLog)

	b.client.AddEventHandler(b.eventHandler)

	err = b.client.Connect()
	if err != nil {
		initError =  fmt.Errorf("failed to login with whatsapp")
		return 
	}
	})

	return initError
}

	func (b *Bot) eventHandler(event interface{}) {
			switch v := event.(type) {
			case *events.Connected:
				close(b.qrChan)
			case *events.QR:
				for _, code := range v.Codes {
					qrBytes, err := qrcode.Encode(code, qrcode.Highest, 256)
					if err != nil {
						log.Printf("Error converting code to QRBytes %v", err)
						continue
					}

					b.qrChan <- base64.RawStdEncoding.EncodeToString(qrBytes)

					time.Sleep(time.Second * 20)
				}
			}
	}