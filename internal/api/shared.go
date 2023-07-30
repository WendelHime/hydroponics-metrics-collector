package api

import (
	"errors"
	"net/http"

	localErrs "github.com/WendelHime/hydroponics-metrics-collector/internal/shared/errors"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
)

func renderErr(w http.ResponseWriter, r *http.Request, err error) {
	var apiErr *localErrs.Error
	if errors.As(err, &apiErr) {
		log.Warn().Err(err).Msg("request failed")
		render.Render(w, r, apiErr)
		return
	}

	log.Warn().Err(err).Msg("internal server error")
	render.Status(r, http.StatusInternalServerError)
}
