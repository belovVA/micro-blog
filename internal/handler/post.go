package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"micro-blog/internal/converter"
	"micro-blog/internal/handler/dto"
	"micro-blog/internal/handler/pkg/response"
	"micro-blog/internal/model"
	"micro-blog/pkg/logger"
)

type PostService interface {
	CreatePost(ctx context.Context, post *model.Post) (*model.Post, error)
	GetListPost(ctx context.Context) ([]*model.Post, error)
	LikePost(ctx context.Context, like *model.Like) error
}

type PostHandler struct {
	Service PostService
	logger  *slog.Logger
}

func NewPostHandler(service PostService, logger *slog.Logger) *PostHandler {
	return &PostHandler{
		Service: service,
		logger:  logger,
	}
}

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	var req dto.CreatePostReq

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

	postModel, err := converter.ToPostModelFromReq(&req)
	if err != nil {
		response.WriteError(w, ErrUUIDParsing, http.StatusBadRequest)
		h.logger.Info(ErrRequestFields, slog.String(logger.ErrorKey, err.Error()))
		return
	}

	post, err := h.Service.CreatePost(r.Context(), postModel)
	if err != nil {
		response.WriteError(w, err.Error(), http.StatusBadRequest)
		h.logger.Info("error to create post", slog.String(logger.ErrorKey, err.Error()))
		return
	}

	resp := converter.ToPostRespFromModel(post)
	h.logger.InfoContext(r.Context(), "successful created")

	response.SuccessJSON(w, resp, http.StatusCreated)
}

func (h *PostHandler) GetPostList(w http.ResponseWriter, r *http.Request) {
	posts, err := h.Service.GetListPost(r.Context())
	if err != nil {
		response.WriteError(w, err.Error(), http.StatusBadRequest)
		h.logger.Info("error to get posts info", slog.String(logger.ErrorKey, err.Error()))
		return
	}

	postsResp := make([]*dto.PostResp, len(posts))
	for i, post := range posts {
		postsResp[i] = converter.ToPostRespFromModel(post)
	}

	h.logger.InfoContext(r.Context(), "successful get posts list")
	response.SuccessJSON(w, postsResp, http.StatusOK)
}

func (h *PostHandler) LikePost(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	const prefix = "/posts/"
	const suffix = "/like"

	if !strings.HasPrefix(path, prefix) || !strings.HasSuffix(path, suffix) {
		response.WriteError(w, ErrNotFound, http.StatusNotFound)
		h.logger.Info(ErrNotFound, slog.String(logger.ErrorKey, path))
		return
	}

	// Извлекаем id из пути
	idPartStr := strings.TrimSuffix(strings.TrimPrefix(path, prefix), suffix)
	idPartStr = strings.Trim(idPartStr, "/")

	var req dto.LikeRequest

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

	likeModel, err := converter.ToLikeModelFromReq(&req, idPartStr)
	if err != nil {
		response.WriteError(w, ErrUUIDParsing, http.StatusBadRequest)
		h.logger.Info(ErrRequestFields, slog.String(logger.ErrorKey, err.Error()))
		return
	}

	err = h.Service.LikePost(r.Context(), likeModel)
	if err != nil {
		response.WriteError(w, err.Error(), http.StatusBadRequest)
		h.logger.Info("error to like post", slog.String(logger.ErrorKey, err.Error()))
		return
	}

	h.logger.InfoContext(r.Context(), "successful liked post")
	response.SuccessCode(w, http.StatusOK)
}
