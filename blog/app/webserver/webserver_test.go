package webserver

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/assert"
)

func TestApiPost_Create(t *testing.T) {

	config := NewConfig()
	_, err := toml.DecodeFile("../../configs/blog.toml", config)
	if err != nil {
		t.Fatal(err)
	}

	db, err := newSession(config.DatabaseConnectionString)
	if err != nil {
		t.Fatal(err)
	}

	serv := newServer(db, config)

	testCases := []struct {
		name         string
		payload      interface{}
		expectedCode int
	}{
		{
			name: "Create New Post Valid",
			payload: map[string]interface{}{
				"title": "NewTestPost",
				"short": "Short decription",
				"body":  "Full body",
			},
			expectedCode: http.StatusCreated,
		},
		{
			name:         "Invalid params",
			payload:      `{"ds":"ss"}`,
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPut, "/api/v1/post", b)
			serv.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}

}
