package api

import (
	"encoding/base64"
	"encoding/json"
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

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&account)
	if err != nil {
		log.Warn().Err(err).Msg("failed to decode account")
		render.Status(r, http.StatusBadRequest)
		return
	}

	err = e.logic.CreateAccount(r.Context(), account)
	if err != nil {
		// TODO add error package and validate returned errors
		log.Error().Err(err).Msg("failed to create account")
		render.Status(r, http.StatusInternalServerError)
		return
	}

	render.Status(r, 200)
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

func (e UserEndpoints) Login(w http.ResponseWriter, r *http.Request) {
	authorization := r.Header.Get("X-Apigateway-Api-Userinfo")
	if len(authorization) == 0 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	payload := strings.Replace(authorization, "Basic ", "", 1)
	payloadDecoded, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	userPass := strings.Split(string(payloadDecoded), ":")

	credentials := models.Credentials{Email: userPass[0], Password: userPass[1]}
	accessToken, err := e.logic.Login(r.Context(), credentials)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	response := LoginResponse{AccessToken: accessToken}
	encoder := json.NewEncoder(w)
	err = encoder.Encode(response)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	render.Status(r, 200)
}
