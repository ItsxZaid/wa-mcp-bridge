package whatsapp

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
)

type MockResponseWriter struct {
	bytes.Buffer
	Flushed bool
}

func TestHandleQR_Intergration(t *testing.T) {
	bot, _ := New(nil)
 	router := chi.NewRouter()
 	bot.RegisterRoutes(router)

	time.Sleep(time.Millisecond * 40)

 	ctx, cancel := context.WithCancel(context.Background())

 	req := httptest.NewRequest("GET", "/qr", nil).WithContext(ctx)
 	rr := httptest.NewRecorder()

	handlerFinnished := make(chan bool)

	go func () {
 	router.ServeHTTP(rr, req)

 	handlerFinnished <- true
	} ()

	// whatsmeow lib need some time to handshake with WA Server's
	time.Sleep(time.Second * 5)

	cancel()

	select {
	case <-handlerFinnished:
	case <-time.After(1 * time.Second):
		t.Fatal("handler did not canceled before given time")
	}

 	if rr.Code != http.StatusOK {
		t.Errorf("handler returned a wrong status code: got %v want %v", rr.Code, http.StatusOK)
 	}

 	expectedHeader := "text/event-stream"
 	if contentType := rr.Header().Get("Content-Type"); contentType != expectedHeader {
		t.Errorf("wrong content type: got %v want %v", contentType, expectedHeader)
 	}

 	if !strings.Contains(rr.Body.String(), ": ping") {
		t.Errorf("body does not contain initial ping: %s", rr.Body.String())
	}
}

func TestHandleServeQR(t *testing.T) {
	bot, _ := New(nil)
	mockWriter := &MockResponseWriter{}
	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan bool)

	go func ()  {
		bot.serveQREvents(ctx, mockWriter, mockWriter)
		done <- true
	} ()

	fakeQRCode := "fake-qr-code"
	bot.qrChan <- fakeQRCode

	time.Sleep(time.Millisecond * 50)

	cancel()

	<- done		

	output := mockWriter.String()
	t.Logf("Got ouput: %v", output)

	if !strings.Contains(output, "event: qr") {
		t.Errorf("output does not contain 'event: qr'")
	}

	if !strings.Contains(output, "data: " + fakeQRCode) {
		t.Errorf("output does not contain correct QR Code")
	}

	if !mockWriter.Flushed {
		t.Errorf("Flush() was never called")
	}
}

func (m *MockResponseWriter) Flush() {
	m.Flushed = true
}