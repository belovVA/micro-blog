package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"micro-blog/internal/converter"
	"micro-blog/internal/handler/dto"
	"micro-blog/internal/handler/pkg/response"
	"micro-blog/internal/model"
	"micro-blog/pkg/logger"
)

type UserService interface {
	Authenticate(ctx context.Context, user *model.User) (*model.User, error)
}

type UserHandler struct {
	Service UserService
	logger  *slog.Logger
}

func NewUserHandler(service UserService, logger *slog.Logger) *UserHandler {
	return &UserHandler{
		Service: service,
		logger:  logger,
	}
}

func (h *UserHandler) Authenticate(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, ErrBodyRequest, http.StatusBadRequest)
		h.logger.Info(ErrBodyRequest, slog.String(logger.ErrorKey, err.Error()))
		return
	}

	v := getValidator(r)
	if err := v.Struct(req); err != nil {
		response.WriteError(w, ErrRequestFields, http.StatusBadRequest)
		h.logger.Info(ErrRequestFields, slog.String(logger.ErrorKey, err.Error()))
		return
	}

	userModel := converter.ToUserModelFromReq(&req)

	user, err := h.Service.Authenticate(r.Context(), userModel)
	if err != nil {
		response.WriteError(w, err.Error(), http.StatusBadRequest)
		h.logger.Info("error to register user", slog.String(logger.ErrorKey, err.Error()))
		return
	}

	resp := converter.ToUserRespFromModel(user)
	h.logger.InfoContext(r.Context(), "successful register")

	response.SuccessJSON(w, resp, http.StatusCreated)
}
