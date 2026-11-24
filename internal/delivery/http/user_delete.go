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
)

// Delete godoc
// @Summary      Delete a user
// @Description  Deletes a user account by ID.
// @Description  Possible 400 errors: invalid id, invalid user data.
// @Description  Possible 401 errors: userID not in context, role not in context.
// @Description  Possible 403 errors: cannot delete another user's account.
// @Tags         users
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      204  "No content"
// @Failure      400  {object}  dto.ErrorResponse  "Invalid id or user data"
// @Failure      401  {object}  dto.ErrorResponse  "Unauthorized"
// @Failure      403  {object}  dto.ErrorResponse  "Forbidden"
// @Failure      404  {object}  dto.ErrorResponse  "User not found"
// @Failure      500  {object}  dto.ErrorResponse  "Internal server error"
// @Router       /v1/admin/users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
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

	err = h.deleteUserService.Execute(c.Request.Context(), id, userID, role)
	if err != nil {
		switch {
		case errors.Is(err, domainErr.ErrInvalidUserData):
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid user data"})
		case errors.Is(err, middleware.ErrForbidden):
			c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "cannot delete other user account"})
		case errors.Is(err, domainErr.ErrUserNotFound):
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "user not found"})
		default:
			slog.ErrorContext(c.Request.Context(), "failed to delete user", "error", err)
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal server error"})
		}

		return
	}

	c.Status(http.StatusNoContent)
}
