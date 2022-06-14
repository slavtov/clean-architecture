package http

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/slavken/clean-architecture/internal/config"
	"github.com/slavken/clean-architecture/internal/domain/models"
	"github.com/slavken/clean-architecture/internal/domain/usecases"
	"github.com/slavken/clean-architecture/internal/middleware"
	"github.com/slavken/clean-architecture/pkg/logger"
	"github.com/slavken/clean-architecture/pkg/utils"
)

type handler struct {
	articleUseCase usecases.ArticleUseCase
	userUseCase    usecases.UserUseCase
	log            logger.Logger
}

func newHandler(
	au usecases.ArticleUseCase,
	uu usecases.UserUseCase,
	log logger.Logger,
) *handler {
	return &handler{
		articleUseCase: au,
		userUseCase:    uu,
		log:            log,
	}
}

func Init(
	cfg *config.Config,
	e *echo.Group,
	au usecases.ArticleUseCase,
	uu usecases.UserUseCase,
	log logger.Logger,
) {
	h := newHandler(au, uu, log)
	auth := middleware.Auth(cfg, uu, log)

	e.GET("/articles", h.GetAll)
	e.GET("/articles/:id", h.GetByID)
	e.POST("/articles", h.Store, auth)
	e.PUT("/articles/:id", h.Update, auth)
	e.DELETE("/articles/:id", h.Delete, auth)
}

// GetAll godoc
// @Tags Articles
// @Summary Get all articles
// @Accept json
// @Produce json
// @Success 200 {object} models.ArticlesList
// @Failure 500 {object} swagger.Error
// @Router /articles [get]
func (h *handler) GetAll(c echo.Context) error {
	res, err := h.articleUseCase.GetAll()
	if err != nil {
		h.log.Errorf("article.UseCase.GetAll: %v", err)
		return err
	}

	return c.JSON(http.StatusOK, &models.ArticlesList{
		TotalCount: len(res),
		Articles:   res,
	})
}

// GetByID godoc
// @Tags Articles
// @Summary Get article by ID
// @Accept json
// @Produce json
// @Param id path string true "Article ID"
// @Success 200 {object} models.Article
// @Failure 400,404,500 {object} swagger.Error
// @Router /articles/{id} [get]
func (h *handler) GetByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.ErrNotFound
	}

	article, err := h.articleUseCase.GetByID(id)
	if err != nil {
		h.log.Errorf("article.UseCase.GetByID: %v", err)
		return err
	}

	return c.JSON(http.StatusOK, article)
}

// GetByID godoc
// @Tags Articles
// @Summary Add article
// @Accept json
// @Produce json
// @Param body body swagger.ArticleRequest true "Body"
// @Security ApiKeyAuth
// @Success 200 {object} models.Article
// @Failure 400,401,500 {object} swagger.Error
// @Router /articles [post]
func (h *handler) Store(c echo.Context) error {
	a := new(models.Article)

	if err := c.Bind(a); err != nil {
		return echo.ErrBadRequest
	}

	a.AuthorID = utils.GetCtxID(c)

	createdArticle, err := h.articleUseCase.Store(a)
	if err != nil {
		h.log.Errorf("article.UseCase.Store: %v", err)
		return err
	}

	return c.JSON(http.StatusCreated, createdArticle)
}

// Update godoc
// @Tags Articles
// @Summary Update article
// @Accept json
// @Produce json
// @Param id path string true "Article ID"
// @Param body body swagger.ArticleRequest true "Body"
// @Security ApiKeyAuth
// @Success 200 {object} models.Article
// @Failure 400,401,404,500 {object} swagger.Error
// @Router /articles/{id} [put]
func (h *handler) Update(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.ErrNotFound
	}

	a := new(models.Article)

	if err := c.Bind(a); err != nil {
		return echo.ErrBadRequest
	}

	a.ID = id
	a.AuthorID = utils.GetCtxID(c)

	updatedArticle, err := h.articleUseCase.Update(a)
	if err != nil {
		h.log.Errorf("article.UseCase.Update: %v", err)
		return err
	}

	return c.JSON(http.StatusOK, updatedArticle)
}

// Delete godoc
// @Tags Articles
// @Summary Delete article
// @Accept json
// @Produce json
// @Param id path string true "Article ID"
// @Security ApiKeyAuth
// @Success 204
// @Failure 400,401,404,500 {object} swagger.Error
// @Router /articles/{id} [delete]
func (h *handler) Delete(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.ErrNotFound
	}

	if err := h.articleUseCase.Delete(models.Article{
		ID:       id,
		AuthorID: utils.GetCtxID(c),
	}); err != nil {
		h.log.Errorf("article.UseCase.Delete: %v", err)
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
