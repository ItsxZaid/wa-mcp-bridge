package whatsapp

import (
	"context"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
)

type Bot struct {
	// TODO: Add fields for WhatsApp connection and state
}

func New() (*Bot, error) {
	return &Bot{}, nil
}

func createWAInstance() {
	dbLog := waLog.Stdout("Database", "", false)
	ctx := context.Background()
	container, err := sqlstore.New(ctx, "sqlite3", "file:whatsapp.db?_foreign_keys=on", dbLog)
	if err != nil {
		panic(err)
	}

	deviceStore, err := container.GetFirstDevice(ctx)
	if err != nil {
		panic(err)
	}

	clientLog := waLog.Stdout("Client", "", false)
	client := whatsmeow.NewClient(deviceStore, clientLog)

	// TODO: Add QR Code handling and connection logic
}
