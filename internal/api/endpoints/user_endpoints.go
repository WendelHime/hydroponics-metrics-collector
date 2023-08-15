package endpoints

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/WendelHime/hydroponics-metrics-collector/internal/logic"
	"github.com/WendelHime/hydroponics-metrics-collector/internal/shared/errors"
	"github.com/WendelHime/hydroponics-metrics-collector/internal/shared/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
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
		errors.RenderErr(w, r, err)
		return
	}

	err = e.logic.CreateAccount(r.Context(), account)
	if err != nil {
		errors.RenderErr(w, r, err)
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
		errors.RenderErr(w, r, err)
		return
	}

	response := LoginResponse{AccessToken: token.AccessToken}
	render.Render(w, r, response)
	render.Status(r, http.StatusOK)
}

type AddDeviceRequest struct {
	Device string `json:"device" validate:"required"`
	UserID string `validate:"required"`
}

func (a *AddDeviceRequest) Bind(r *http.Request) error {
	validate := validator.New()
	err := validate.Struct(a)
	if err != nil {
		return errors.BadRequestErr.WithErr(err).WithMsg("failed to validate request")
	}
	return nil
}

func (e UserEndpoints) AddDevice(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	request := AddDeviceRequest{UserID: userID}
	err := render.Bind(r, &request)
	if err != nil {
		log.Warn().Err(err).Msg("failed to decode add device request")
		errors.RenderErr(w, r, err)
		return
	}

	err = e.logic.AddDevice(r.Context(), userID, request.Device)
	if err != nil {
		log.Error().Err(err).Msg("failed to add new device")
		errors.RenderErr(w, r, err)
		return
	}

	render.Status(r, http.StatusAccepted)
}

type GetDevicesResponse struct {
	UserID  string   `json:"user_id"`
	Devices []string `json:"devices"`
}

func (g GetDevicesResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (e UserEndpoints) GetDevices(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	if len(userID) == 0 {
		errors.RenderErr(w, r, errors.BadRequestErr.WithMsg("missing user ID"))
		return
	}

	devices, err := e.logic.GetDevices(r.Context(), userID)
	if err != nil {
		log.Error().Err(err).Msg("failed to retrieve devices")
		errors.RenderErr(w, r, err)
		return
	}

	response := GetDevicesResponse{
		UserID:  userID,
		Devices: devices,
	}

	render.Render(w, r, response)
	render.Status(r, http.StatusOK)
}
