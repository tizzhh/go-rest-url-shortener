package redirect

import (
	"errors"
	"log/slog"
	"net/http"
	resp "url-shortener/internal/api/response"
	"url-shortener/internal/storage"
	sl "url-shortener/pkg/logger/slog"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

//go:generate go run github.com/vektra/mockery/v2 --name=URLGetter
type URLGetter interface {
	GetURL(alias string) (string, error)
}

var log *slog.Logger = sl.GetLogger()

func Redirect(urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const caller = "handlers.url.redirect.Redirect"

		log = log.With(
			slog.String("caller", caller),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("failed to get alias from url", slog.String("url", r.URL.Path))
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}

		destinationURL, err := urlGetter.GetURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url for alias not found", slog.String("alias", alias))
			render.JSON(w, r, resp.Error("not found"))
			return
		}
		if err != nil {
			log.Info("failed to retrieve url", sl.Err(err))
			render.JSON(w, r, resp.Error("internal error"))
			return
		}

		log.Info(
			"url for alias retrieved",
			slog.String("url", destinationURL),
			slog.String("alias", alias),
		)

		http.Redirect(w, r, destinationURL, http.StatusFound)
	}
}
