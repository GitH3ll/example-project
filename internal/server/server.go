package server

import (
	"encoding/json"
	"github.com/GitH3ll/example-project/internal/model"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strconv"
)

type user struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type controller interface {
	AddUser(user model.User) error
	GetUser(id int) (model.User, error)
}

type Server struct {
	listenURI  string
	logger     *logrus.Logger
	controller controller
	r          chi.Router
}

func NewServer(logger *logrus.Logger, controller controller) *Server {
	router := chi.NewRouter()

	return &Server{
		listenURI:  ":8000",
		logger:     logger,
		controller: controller,
		r:          router,
	}
}

func (s *Server) StartServer() {
	srv := http.Server{
		Addr:    s.listenURI,
		Handler: s.r,
	}

	err := srv.ListenAndServe()
	if err != nil {
		s.logger.Error(err)
	}
}

func (s *Server) RegisterRoutes() {
	s.r.Post("/user/add", s.HandleAddUser)
	s.r.Get("/user/{userID}", s.HandleGetUser)
}

func (s *Server) HandleAddUser(w http.ResponseWriter, r *http.Request) {
	var u user

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			s.logger.Error(err)
		}
	}(r.Body)

	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		s.handleError(err, http.StatusInternalServerError, w)
		return
	}

	err = s.controller.AddUser(model.User(u))
	if err != nil {
		s.handleError(err, http.StatusInternalServerError, w)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "userID")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		s.handleError(err, http.StatusBadRequest, w)
		return
	}

	u, err := s.controller.GetUser(id)
	if err != nil {
		s.handleError(err, http.StatusInternalServerError, w)
		return
	}

	err = json.NewEncoder(w).Encode(&u)
	if err != nil {
		s.handleError(err, http.StatusInternalServerError, w)
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleError(err error, status int, w http.ResponseWriter) {
	w.WriteHeader(status)
	_, err = w.Write([]byte(err.Error()))
	if err != nil {
		s.logger.Error(err)
	}
}
