package tests

import (
	"net/url"
	"path"
	"testing"
	"url-shortener/internal/api"
	"url-shortener/internal/http-server/handlers/url/save"
	"url-shortener/internal/lib/random"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/require"
)

const (
	host = "localhost:8081"
)

func TestURLShortener_SaveURL(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}
	e := httpexpect.Default(t, u.String())

	e.POST("/url").
		WithJSON(save.Request{
			URL:   gofakeit.URL(),
			Alias: random.NewRandomString(10),
		}).
		WithBasicAuth("admin", "admin").
		Expect().
		Status(200).
		JSON().
		Object().
		ContainsKey("alias")
}

func TestURLShortener_SaveRedirectRemove(t *testing.T) {
	testCases := []struct {
		name  string
		url   string
		alias string
		error string
	}{
		{
			name:  "valid url",
			url:   gofakeit.URL(),
			alias: gofakeit.Word() + gofakeit.Word(),
		},
		{
			name:  "invalid url",
			url:   "invalid url",
			alias: gofakeit.Word() + gofakeit.Word(),
			error: "field URL is not a valid URL",
		},
		{
			name:  "empty alias",
			url:   gofakeit.URL(),
			alias: "",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			u := url.URL{
				Scheme: "http",
				Host:   host,
			}

			e := httpexpect.Default(t, u.String())

			resp := e.POST("/url").
				WithJSON(save.Request{
					URL:   tt.url,
					Alias: tt.alias,
				}).
				WithBasicAuth("admin", "admin").
				Expect().
				Status(200).
				JSON().
				Object()

			if tt.error != "" {
				resp.NotContainsKey("alias")

				resp.Value("error").String().IsEqual(tt.error)

				return
			}

			alias := tt.alias
			if tt.alias != "" {
				resp.Value("alias").String().IsEqual(tt.alias)
			} else {
				resp.Value("alias").String().NotEmpty()

				alias = resp.Value("alias").String().Raw()
			}

			testRedirect(t, alias, tt.url)

			reqDel := e.DELETE("/"+path.Join("url", alias)).
				WithBasicAuth("admin", "admin").
				Expect().
				Status(200).
				JSON().
				Object()

			reqDel.Value("status").String().IsEqual("OK")

			testRedirectNotFound(t, alias)
		})
	}
}

func testRedirectNotFound(t *testing.T, alias string) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
		Path:   alias,
	}
	_, err := api.GetRedirect(u.String())
	require.ErrorIs(t, err, api.ErrInvalidStatusCode)
}

func testRedirect(t *testing.T, alias string, urlToRedirect string) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
		Path:   alias,
	}

	redirectedToURL, err := api.GetRedirect(u.String())
	require.NoError(t, err)

	require.Equal(t, urlToRedirect, redirectedToURL)
}
