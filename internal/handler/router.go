package handler

import (
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"micro-blog/internal/middleware"
)

const (
	ErrBodyRequest   = "Invalid Request Body"
	ErrRequestFields = "Invalid Request Fields"
	ErrUUIDParsing   = "Invalid UUID"
	ErrNotFound      = "Not Found"
)

type Service interface {
	UserService
	PostService
}

type Router struct {
	service Service
	logger  *slog.Logger
}

func NewRouter(service Service, logger *slog.Logger) http.Handler {
	r := http.NewServeMux()
	router := &Router{
		service: service,
		logger:  logger,
	}

	validate := middleware.NewValidator().Middleware
	recovery := middleware.Recovery(logger)

	// Утилита для оборачивания хендлера
	wrap := func(h http.Handler) http.Handler {
		return recovery(validate(h))
	}

	r.Handle("/register", methodOnly(http.MethodPost, wrap(http.HandlerFunc(router.authHandler))))
	r.Handle("/posts", wrap(http.HandlerFunc(router.postsHandler)))
	r.Handle("/posts/", methodOnly(http.MethodPost, wrap(http.HandlerFunc(router.postLikeHandler))))

	return r
}

// Позволяет ограничить вызов ручки только конкретным методом
func methodOnly(method string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func getValidator(r *http.Request) *validator.Validate {
	if v, ok := r.Context().Value("validator").(*validator.Validate); ok {
		return v
	}
	return validator.New()
}

func (r *Router) authHandler(w http.ResponseWriter, req *http.Request) {
	h := NewUserHandler(r.service, r.logger)
	h.Authenticate(w, req)
}

func (r *Router) postsHandler(w http.ResponseWriter, req *http.Request) {
	h := NewPostHandler(r.service, r.logger)
	switch req.Method {
	case http.MethodPost:
		h.CreatePost(w, req)
	case http.MethodGet:
		h.GetPostList(w, req)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (r *Router) postLikeHandler(w http.ResponseWriter, req *http.Request) {
	h := NewPostHandler(r.service, r.logger)
	h.LikePost(w, req)
}
