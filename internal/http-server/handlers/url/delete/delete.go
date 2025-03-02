package delete

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"url-shortener/constants"
	"url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/lib/random"
	"url-shortener/internal/storage"
)

func New(log *slog.Logger, urlDeleter storage.URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log = log.With(
			slog.String("op", constants.UrlDeleteNew),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, response.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			var validateErr validator.ValidationErrors
			if errors.As(err, &validateErr) {
				log.Error("invalid request", sl.Err(err))
				render.JSON(w, r, response.ValidationError(validateErr))
				return
			}

			log.Error("unexpected validation error", sl.Err(err))
			render.JSON(w, r, response.Error("invalid request"))
			return
		}

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(constants.AliasLength)
		}

		id, err := urlDeleter.DeleteUrl(r.Context(), alias)
		if errors.Is(err, storage.ErrAliasExists) {
			log.Info("alias already exists", slog.String("alias", req.Alias))

			render.JSON(w, r, response.Error("alias already exists"))

			return
		}
		if err != nil {
			log.Error("failed to delete url", sl.Err(err))

			render.JSON(w, r, response.Error("failed to delete url"))

			return
		}

		log.Info("url deleted", slog.Int64("id", id))

		render.JSON(w, r, Response{
			Response: response.OK(),
			Alias:    alias,
		})
	}
}
