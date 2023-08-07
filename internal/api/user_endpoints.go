package api

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/WendelHime/hydroponics-metrics-collector/internal/logic"
	"github.com/WendelHime/hydroponics-metrics-collector/internal/shared/models"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
)

type UserEndpoints struct {
	logic logic.UserLogic
}

func NewUserEndpoint(l logic.UserLogic) UserEndpoints {
	return UserEndpoints{logic: l}
}

func (e UserEndpoints) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var account models.User

	err := render.Bind(r, &account)
	if err != nil {
		log.Warn().Err(err).Msg("failed to decode account")
		renderErr(w, r, err)
		return
	}

	err = e.logic.CreateAccount(r.Context(), account)
	if err != nil {
		renderErr(w, r, err)
		return
	}

	render.Status(r, http.StatusCreated)
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

func (l LoginResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (e UserEndpoints) SignIn(w http.ResponseWriter, r *http.Request) {
	authorization := r.Header.Get("X-Apigateway-Api-Userinfo")
	if len(authorization) == 0 {
		render.Status(r, http.StatusBadRequest)
		return
	}

	payload := strings.Replace(authorization, "Basic ", "", 1)
	payloadDecoded, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		return
	}
	userPass := strings.Split(string(payloadDecoded), ":")

	credentials := models.Credentials{Email: userPass[0], Password: userPass[1]}
	token, err := e.logic.Login(r.Context(), credentials)
	if err != nil {
		renderErr(w, r, err)
		return
	}

	response := LoginResponse{AccessToken: token.AccessToken}
	render.Render(w, r, response)
	render.Status(r, http.StatusOK)
}
