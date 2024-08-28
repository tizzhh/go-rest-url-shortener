package redirect

import (
	"net/http/httptest"
	"testing"
	"url-shortener/internal/api"
	"url-shortener/internal/http-server/handlers/url/redirect/mocks"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func TestRedirectHandler(t *testing.T) {
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
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			urlGetterMock := mocks.NewURLGetter(t)

			if tt.respError == "" || tt.mockError != nil {
				urlGetterMock.On("GetURL", tt.alias).
					Return(tt.url, tt.mockError).
					Once()
			}

			r := chi.NewRouter()
			r.Get("/{alias}", Redirect(urlGetterMock))

			tc := httptest.NewServer(r)
			defer tc.Close()

			redirectedToUrl, err := api.GetRedirect(tc.URL + "/" + tt.alias)
			require.NoError(t, err)

			require.Equal(t, tt.url, redirectedToUrl)
		})
	}
}
