package http

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rober0xf/notifier/internal/delivery/http/dto"
	domainErr "github.com/rober0xf/notifier/internal/domain/errors"
	"github.com/rober0xf/notifier/internal/middleware"
	"github.com/rober0xf/notifier/internal/usecase/user"
)

// Update godoc
// @Summary      Update a user
// @Description  Updates a user's data by ID. At least one field is required.
// @Description  Possible 400 errors: invalid id, invalid email format, invalid password, missing fields.
// @Description  Possible 409 errors: email already in use, username already in use.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id       path      int                   true  "User ID"
// @Param        payload  body      dto.UpdateUserRequest true  "Fields to update"
// @Success      200      {object}  dto.UserPayload
// @Failure      400      {object}  dto.ErrorResponse  "Invalid request"
// @Failure      401      {object}  dto.ErrorResponse  "Unauthorized"
// @Failure      403      {object}  dto.ErrorResponse  "Forbidden"
// @Failure      404      {object}  dto.ErrorResponse  "User not found"
// @Failure      409      {object}  dto.ErrorResponse  "Conflict"
// @Failure      500      {object}  dto.ErrorResponse  "Internal server error"
// @Security     BearerAuth
// @Router       /v1/auth/users/{id} [put]
func (h *UserHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}

	if id <= 0 {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "id must be positive"})
		return
	}

	// parse json to dto
	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
		return
	}

	if req.Username == nil && req.Email == nil && req.Password == nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "at least one field is required"})
		return
	}

	// create the user with the new data
	input := user.UpdateUserInput{
		ID:       id,
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}

	role, err := middleware.GetRoleFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}

	updatedUser, err := h.updateUserService.Execute(c.Request.Context(), input, userID, role)
	if err != nil {
		switch {
		case errors.Is(err, domainErr.ErrUserNotFound):
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "user not found"})
		case errors.Is(err, domainErr.ErrEmailAlreadyExists):
			c.JSON(http.StatusConflict, dto.ErrorResponse{Error: "email already in use"})
		case errors.Is(err, domainErr.ErrUsernameAlreadyExists):
			c.JSON(http.StatusConflict, dto.ErrorResponse{Error: "username already in use"})
		case errors.Is(err, middleware.ErrForbidden):
			c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "cannot change other user data"})
		case errors.Is(err, domainErr.ErrInvalidEmailFormat):
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid email format"})
		case errors.Is(err, domainErr.ErrInvalidPassword):
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid password"})
		default:
			slog.ErrorContext(c.Request.Context(), "failed to update user", "error", err)
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal server error"})
		}

		return
	}

	c.JSON(http.StatusOK, dto.ToUserResponse(updatedUser))
}
