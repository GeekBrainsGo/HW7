package server

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
)

// Only in test files are visible
var (
	testStrings = []string{
		"Hello Test!",
		"Any text",
		"Русский текст",
		"{\"data\":\"value\"}",
	}
)

func TestHandleEcho(t *testing.T) {
	lg := logrus.New()
	lg.SetLevel(0)
	serv := NewServer(context.Background(), lg).Start(":8080")

	for _, tcase := range testStrings {

		reader := bytes.NewReader([]byte(tcase))
		req, err := http.NewRequest("POST", "/echo", reader)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(serv.HandleEcho)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("status code isn't 200, got: %v", status)
		}

		if rr.Body.String() != tcase {
			t.Errorf("unexpected body, got: %v, want %v", rr.Body.String(), tcase)
		}

	}

	serv.Stop()
}

func BenchmarkHandleEcho(b *testing.B) {
	lg := logrus.New()
	lg.SetLevel(0)
	serv := NewServer(context.Background(), lg).Start(":8080")
	data := "DATA TEXT"

	for i := 0; i < b.N; i++ {

		reader := bytes.NewReader([]byte(data))
		req, err := http.NewRequest("POST", "/echo", reader)
		if err != nil {
			b.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(serv.HandleEcho)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			b.Errorf("status code isn't 200, got: %v", status)
		}

		if rr.Body.String() != data {
			b.Errorf("unexpected body, got: %v, want %v", rr.Body.String(), data)
		}

	}

	serv.Stop()
}
