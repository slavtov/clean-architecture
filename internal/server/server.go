package server

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	_ "github.com/slavken/go-clean-architecture/docs"
	"github.com/slavken/go-clean-architecture/internal/config"
	"github.com/slavken/go-clean-architecture/pkg/logger"
	"github.com/slavken/go-clean-architecture/pkg/store/redis"
)

type Server struct {
	cfg    *config.Config
	router *echo.Echo
	db     *sqlx.DB
	redis  redis.Store
	log    logger.Logger
}

func New(
	cfg *config.Config,
	db *sqlx.DB,
	rdb redis.Store,
	log logger.Logger,
) *Server {
	return &Server{
		cfg:    cfg,
		router: echo.New(),
		db:     db,
		redis:  rdb,
		log:    log,
	}
}

func (s *Server) Run() error {
	s.middleware()
	s.handlers()

	return s.router.Start(s.cfg.Server.Addr)
}
