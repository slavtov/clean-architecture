package server

import (
	"github.com/labstack/echo/v4/middleware"
	articleDelivery "github.com/slavken/go-clean-architecture/internal/article/delivery/http"
	articleRepository "github.com/slavken/go-clean-architecture/internal/article/repository"
	articleUseCase "github.com/slavken/go-clean-architecture/internal/article/usecase"
	authDelivery "github.com/slavken/go-clean-architecture/internal/auth/delivery/http"
	authRepository "github.com/slavken/go-clean-architecture/internal/auth/repository"
	authUseCase "github.com/slavken/go-clean-architecture/internal/auth/usecase"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func (s *Server) middleware() {
	s.router.Pre(middleware.RemoveTrailingSlash())
	s.router.Use(middleware.CORS())
}

func (s *Server) handlers() {
	authRepo := authRepository.NewPGRepository(s.db)
	authRedisRepo := authRepository.NewRedisRepository(s.redis)
	articleRepo := articleRepository.NewPGRepository(s.db)
	articleRedisRepo := articleRepository.NewRedisRepository(s.redis)

	authUC := authUseCase.New(
		s.cfg,
		authRepo,
		authRedisRepo,
		s.log,
	)
	articleUC := articleUseCase.New(
		articleRepo,
		articleRedisRepo,
		s.log,
	)

	if s.cfg.Server.Debug {
		s.router.GET("/swagger/*", echoSwagger.WrapHandler)
	}

	api := s.router.Group("/api")

	authDelivery.Init(
		s.cfg,
		api,
		authUC,
		s.log,
	)
	articleDelivery.Init(
		s.cfg,
		api,
		articleUC,
		authUC,
		s.log,
	)
}
