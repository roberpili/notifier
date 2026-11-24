package http

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rober0xf/notifier/internal/delivery/http/dto"
	"github.com/rober0xf/notifier/internal/domain/entity"
	domainErr "github.com/rober0xf/notifier/internal/domain/errors"
)

// GetByID godoc
// @Summary      Get user by ID
// @Description  Returns a single user by their ID.
// @Tags         users
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  dto.UserPayload
// @Failure      400  {object}  dto.ErrorResponse  "Invalid id"
// @Failure      404  {object}  dto.ErrorResponse  "User not found"
// @Failure      500  {object}  dto.ErrorResponse  "Internal server error"
// @Security     BearerAuth
// @Router       /v1/admin/users/{id} [get]
func (h *UserHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}

	if id <= 0 {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "id must be positive"})
		return
	}

	user, err := h.getUserByIDService.Execute(c.Request.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, domainErr.ErrInvalidUserData):
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "invalid user data"})
		case errors.Is(err, domainErr.ErrUserNotFound):
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "user not found"})
		default:
			slog.ErrorContext(c.Request.Context(), "failed to get user by id", "error", err)
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal server error"})
		}

		return
	}

	c.JSON(http.StatusOK, dto.ToUserResponse(user))
}

// GetByEmail godoc
// @Summary      Get user by email
// @Description  Returns a single user by their email.
// @Tags         users
// @Produce      json
// @Param        email   path      string  true  "User email"
// @Success      200  {object}  dto.UserPayload
// @Failure      400  {object}  dto.ErrorResponse  "Invalid email format"
// @Failure      404  {object}  dto.ErrorResponse  "User not found"
// @Failure      500  {object}  dto.ErrorResponse  "Internal server error"
// @Security     BearerAuth
// @Router       /v1/admin/users/email/{email} [get]
func (h *UserHandler) GetByEmail(c *gin.Context) {
	email := c.Param("email")

	if !strings.Contains(email, "@") {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid email format"})
		return
	}

	user, err := h.getUserByEmailService.Execute(c.Request.Context(), email)
	if err != nil {
		switch {
		case errors.Is(err, domainErr.ErrInvalidEmailFormat):
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid email format"})
		case errors.Is(err, domainErr.ErrUserNotFound):
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "user not found"})
		default:
			slog.ErrorContext(c.Request.Context(), "failed to get user by email", "error", err)
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal server error"})
		}

		return
	}

	c.JSON(http.StatusOK, dto.ToUserResponse(user))
}

// GetAll godoc
// @Summary      Get all users
// @Description  Returns all users. Admin only.
// @Tags         users
// @Produce      json
// @Success      200  {array}   dto.UserPayload
// @Failure      500  {object}  dto.ErrorResponse  "Internal server error"
// @Security     BearerAuth
// @Router       /v1/admin/users [get]
func (h *UserHandler) GetAll(c *gin.Context) {
	users, err := h.getAllUsersService.Execute(c.Request.Context())
	if err != nil {
		slog.ErrorContext(c.Request.Context(), "failed to get all users", "error", err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal server error"})
		return
	}

	response := make([]entity.User, 0, len(users))
	response = append(response, users...)

	c.JSON(http.StatusOK, response)
}
