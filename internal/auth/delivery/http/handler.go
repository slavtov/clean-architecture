package http

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/slavtov/clean-architecture/internal/config"
	"github.com/slavtov/clean-architecture/internal/domain/models"
	"github.com/slavtov/clean-architecture/internal/domain/usecases"
	"github.com/slavtov/clean-architecture/internal/middleware"
	"github.com/slavtov/clean-architecture/pkg/logger"
	"github.com/slavtov/clean-architecture/pkg/utils"
)

type handler struct {
	cfg         *config.Config
	userUseCase usecases.UserUseCase
	log         logger.Logger
}

func newHandler(
	cfg *config.Config,
	uu usecases.UserUseCase,
	log logger.Logger,
) *handler {
	return &handler{cfg, uu, log}
}

func Init(
	cfg *config.Config,
	e *echo.Group,
	uu usecases.UserUseCase,
	log logger.Logger,
) {
	h := newHandler(cfg, uu, log)
	auth := middleware.Auth(cfg, uu, log)
	clearCookies := middleware.ClearCookies(cfg, log)

	authGroup := e.Group("/auth")
	authGroup.POST("/me", h.Me, auth)
	authGroup.POST("/login", h.Login)
	authGroup.POST("/register", h.Register)
	authGroup.POST("/refresh", h.Refresh)
	authGroup.POST("/logout", h.Logout, auth, clearCookies)
	authGroup.POST("/logout/all", h.LogoutAll, auth, clearCookies)

	e.GET("/users", h.GetAll)
	e.GET("/users/:id", h.GetByID)
	e.PUT("/users/:id", h.Update, auth)
	e.DELETE("/users/:id", h.Delete, auth)
}

// Me godoc
// @Tags Auth
// @Summary Get auth user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} models.User
// @Failure 400,401,404,500 {object} swagger.Error
// @Router /auth/me [post]
func (h *handler) Me(c echo.Context) error {
	userID := utils.GetCtxID(c)

	res, err := h.userUseCase.GetByID(userID)
	if err != nil {
		h.log.Errorf("auth.UseCase.GetByID: %v", err)
		return err
	}

	return c.JSON(http.StatusOK, res)
}

// GetAll godoc
// @Tags Users
// @Summary Get all users
// @Accept json
// @Produce json
// @Success 200 {object} models.UsersList
// @Failure 500 {object} swagger.Error
// @Router /users [get]
func (h *handler) GetAll(c echo.Context) error {
	res, err := h.userUseCase.GetAll()
	if err != nil {
		h.log.Errorf("auth.UseCase.GetAll: %v", err)
		return err
	}

	return c.JSON(http.StatusOK, &models.UsersList{
		TotalCount: len(res),
		Users:      res,
	})
}

// GetByID godoc
// @Tags Users
// @Summary Get user by ID
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} models.User
// @Failure 400,404,500 {object} swagger.Error
// @Router /users/{id} [get]
func (h *handler) GetByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.ErrNotFound
	}

	user, err := h.userUseCase.GetByID(id)
	if err != nil {
		h.log.Errorf("auth.UseCase.GetByID: %v", err)
		return err
	}

	return c.JSON(http.StatusOK, user)
}

// Login godoc
// @Tags Auth
// @Summary Login user
// @Accept json
// @Produce json
// @Param body body swagger.UserRequest true "Body"
// @Success 200 {object} models.AuthUser
// @Failure 400,401 {object} swagger.Error
// @Router /auth/login [post]
func (h *handler) Login(c echo.Context) error {
	u := new(models.User)

	if err := c.Bind(u); err != nil {
		return echo.ErrBadRequest
	}

	user, err := h.userUseCase.Login(u)
	if err != nil {
		h.log.Errorf("auth.UseCase.Login: %v", err)
		return err
	}

	h.setCookies(c, user)

	return c.JSON(http.StatusOK, user)
}

// Register godoc
// @Tags Auth
// @Summary New user
// @Accept json
// @Produce json
// @Param body body swagger.UserRequest true "Body"
// @Success 201 {object} models.AuthUser
// @Failure 400,500 {object} swagger.Error
// @Router /auth/register [post]
func (h *handler) Register(c echo.Context) error {
	u := new(models.User)

	if err := c.Bind(u); err != nil {
		return echo.ErrBadRequest
	}

	createdUser, err := h.userUseCase.Store(u)
	if err != nil {
		h.log.Errorf("auth.UseCase.Store: %v", err)
		return err
	}

	h.setCookies(c, createdUser)

	return c.JSON(http.StatusCreated, createdUser)
}

// Refresh godoc
// @Tags Auth
// @Summary Using the refresh token
// @Accept json
// @Produce json
// @Success 200 {object} models.AuthUser
// @Failure 400,401,404,500 {object} swagger.Error
// @Router /auth/refresh [post]
func (h *handler) Refresh(c echo.Context) error {
	if err := utils.VerifyRefreshToken(
		c,
		h.cfg.Server.JwtRefreshSecret,
		h.log,
	); err != nil {
		h.log.Errorf("verifyRefreshToken: %v", err)
		return echo.ErrUnauthorized
	}

	userID := utils.GetCtxID(c)
	refreshID := utils.GetCtxRefreshID(c)

	if _, err := h.userUseCase.GetToken(
		userID,
		refreshID,
	); err != nil {
		h.log.Errorf("auth.UseCase.GetToken: %v", err)
		return echo.ErrUnauthorized
	}

	u, err := h.userUseCase.GetByID(userID)
	if err != nil {
		h.log.Errorf("auth.UseCase.GetByID: %v", err)
		return err
	}

	user, err := h.userUseCase.Auth(&u)
	if err != nil {
		h.log.Errorf("auth.UseCase.Auth: %v", err)
		return err
	}

	if err := h.userUseCase.DeleteToken(
		userID,
		refreshID,
	); err != nil {
		h.log.Errorf("auth.UseCase.DeleteToken: %v", err)
		return err
	}

	h.setCookies(c, user)

	return c.JSON(http.StatusOK, user)
}

// Update godoc
// @Summary Update user
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param body body swagger.UpdateUser true "Body"
// @Security ApiKeyAuth
// @Success 200 {object} models.User
// @Failure 400,401,403,404,500 {object} swagger.Error
// @Router /users/{id} [put]
func (h *handler) Update(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.ErrNotFound
	}

	if id != utils.GetCtxID(c) {
		return echo.ErrForbidden
	}

	u := new(models.User)

	if err := c.Bind(u); err != nil {
		return echo.ErrBadRequest
	}

	u.ID = id

	updatedUser, err := h.userUseCase.Update(u)
	if err != nil {
		h.log.Errorf("auth.UseCase.Update: %v", err)
		return err
	}

	return c.JSON(http.StatusOK, updatedUser)
}

// Delete godoc
// @Tags Users
// @Summary Delete user
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Security ApiKeyAuth
// @Success 204
// @Failure 400,401,403,404,500 {object} swagger.Error
// @Router /users/{id} [delete]
func (h *handler) Delete(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.ErrNotFound
	}

	if id != utils.GetCtxID(c) {
		return echo.ErrForbidden
	}

	if err := h.userUseCase.Delete(id); err != nil {
		h.log.Errorf("auth.UseCase.Delete: %v", err)
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

// Logout godoc
// @Tags Auth
// @Summary Log out
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 204
// @Failure 401,404 {object} swagger.Error
// @Router /auth/logout [post]
func (h *handler) Logout(c echo.Context) error {
	userID := utils.GetCtxID(c)
	accessID := utils.GetCtxAccessID(c)
	refreshID := utils.GetCtxRefreshID(c)

	if err := h.userUseCase.Logout(userID, &utils.TokenDetails{
		AtID: accessID,
		RtID: refreshID,
	}); err != nil {
		h.log.Errorf("auth.UseCase.Logout: %v", err)
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

// LogoutAll godoc
// @Tags Auth
// @Summary Log out of all devices
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 204
// @Failure 401,404 {object} swagger.Error
// @Router /auth/logout/all [post]
func (h *handler) LogoutAll(c echo.Context) error {
	userID := utils.GetCtxID(c)

	if err := h.userUseCase.LogoutAll(userID); err != nil {
		h.log.Errorf("auth.UseCase.LogoutAll: %v", err)
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *handler) setCookies(c echo.Context, user *models.AuthUser) {
	c.SetCookie(&http.Cookie{
		Name:     "access_token",
		Value:    user.AccessToken,
		Path:     "/",
		MaxAge:   h.cfg.Cookie.AccessToken.MaxAge,
		Secure:   h.cfg.Cookie.AccessToken.Secure,
		HttpOnly: h.cfg.Cookie.AccessToken.HttpOnly,
		SameSite: http.SameSiteStrictMode,
	})

	c.SetCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    user.RefreshToken,
		Path:     "/",
		MaxAge:   h.cfg.Cookie.RefreshToken.MaxAge,
		Secure:   h.cfg.Cookie.RefreshToken.Secure,
		HttpOnly: h.cfg.Cookie.RefreshToken.HttpOnly,
		SameSite: http.SameSiteStrictMode,
	})
}
