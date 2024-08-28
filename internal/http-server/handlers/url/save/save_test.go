package save

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"url-shortener/internal/http-server/handlers/url/save/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSaveHandler(t *testing.T) {
	testCases := []struct {
		name      string
		alias     string
		url       string
		respError string
		mockError error
	}{
		{
			name:  "success",
			alias: "test",
			url:   "https://ya.ru",
		},
		{
			name:  "empty alias",
			alias: "",
			url:   "https://ya.ru",
		},
		{
			name:      "empty url",
			alias:     "test",
			url:       "",
			respError: "field URL is a required field",
		},
		{
			name:      "invalid url",
			alias:     "test",
			url:       "invalid invalid",
			respError: "field URL is not a valid URL",
		},
		{
			name:      "SaveURL error",
			alias:     "test",
			url:       "https://ya.ru",
			respError: "failed to add url",
			mockError: errors.New("unexpected error"),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			urlSaverMock := mocks.NewURLSaver(t)

			if tt.respError == "" || tt.mockError != nil {
				urlSaverMock.On("SaveURL", tt.url, mock.AnythingOfType("string")).
					Return(tt.mockError).
					Once()
			}

			handler := NewURL(urlSaverMock)

			input := fmt.Sprintf(`{"url": "%s", "alias": "%s"}`, tt.url, tt.alias)

			req, err := http.NewRequest(http.MethodPost, "/save", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, http.StatusOK)

			body := rr.Body.String()

			var resp Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tt.respError, tt.respError)
		})
	}
}
