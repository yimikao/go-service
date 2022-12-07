package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"tutgo/business/comment"
	"tutgo/db"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type server interface {
	Run()
	Stop(ctx context.Context)
}

type svr struct {
	Server *http.Server
}

func (s *svr) Run() {
	log.Printf("server up and run at port %s", s.Server.Addr)

	go func() {
		log.Fatal(s.Server.ListenAndServe())
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	s.Stop(context.Background())

}

func (s *svr) Stop(ctx context.Context) {
	if err := s.Server.Shutdown(ctx); err != nil {

		log.Fatalf("couldn't gracefully shutdown server %v", err)
	}
}

func NewServer(addr string) server {
	return &svr{
		Server: &http.Server{
			Addr:    addr,
			Handler: setupRoutes(addr),
		},
	}
}

func setupRoutes(port string) *chi.Mux {

	r := chi.NewRouter()

	r.Use(middleware.Logger)

	pool, err := db.New()
	if err != nil {
		log.Fatalf("cant create database %v", err)
	}

	cr := comment.NewCommentLayer(pool)
	ch := NewCommentHandler(cr)

	r.Route("/comments", func(r chi.Router) {
		r.Get("/", ch.getComments)
		r.Post("/", AuthUser(ch.createComment))
	})

	return r
}

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
}

func NewErrorResp(m string, c int) *Response {
	return &Response{
		Success: false,
		Message: m,
		Code:    c,
		Data:    nil,
	}
}

func NewResp(m string, c interface{}) *Response {
	return &Response{
		Success: true,
		Message: m,
		Code:    http.StatusOK,
		Data:    c,
	}
}

type CommentHandler struct {
	storer comment.Storer
}

func NewCommentHandler(s comment.Storer) CommentHandler {
	return CommentHandler{
		storer: s,
	}
}

func AuthUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("authed") != "" {
			next.ServeHTTP(w, r)
			return
		}
		w.WriteHeader(http.StatusUnauthorized)
	})
}

func (h *CommentHandler) createComment(w http.ResponseWriter, r *http.Request) {
	var (
		req = new(comment.NewComment)
		res = new(Response)
		err error
	)

	if err = json.NewDecoder(r.Body).Decode(req); err != nil {
		res = NewErrorResp(err.Error(), http.StatusBadRequest)
		if err = json.NewEncoder(w).Encode(res); err != nil {
			log.Printf("couldn't send error response %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	c, err := h.storer.Create(&comment.Comment{
		Comment: req.Comment,
		UserID:  req.UserID,
	})

	if err != nil {
		res = NewErrorResp(err.Error(), http.StatusInternalServerError)
		if err = json.NewEncoder(w).Encode(res); err != nil {
			log.Printf("couldn't send error response %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	res = NewResp("comment created!", c)
	if err = json.NewEncoder(w).Encode(res); err != nil {
		log.Printf("couldn't send error response %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *CommentHandler) getComments(w http.ResponseWriter, r *http.Request) {
	var res = new(Response)

	cms, err := h.storer.All()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			res = NewErrorResp(err.Error(), http.StatusNotFound)
			if err = json.NewEncoder(w).Encode(res); err != nil {
				log.Printf("couldn't send error response %v\n", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			return
		}

		res = NewErrorResp(err.Error(), http.StatusInternalServerError)
		if err = json.NewEncoder(w).Encode(res); err != nil {
			log.Printf("couldn't send error response %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	res = NewResp("comments fetched!", cms)
	if err = json.NewEncoder(w).Encode(res); err != nil {
		log.Printf("couldn't send error response %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
