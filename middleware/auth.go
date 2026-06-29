package middleware

import (
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/shihabul3000/Go-Project/apperrors"
	"github.com/shihabul3000/Go-Project/auth"
)

const (
	ContextUserID = "user_id"
	ContextRole   = "role"
)

func JWTAuth(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			header := c.Request().Header.Get(echo.HeaderAuthorization)
			parts := strings.Fields(header)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				return apperrors.ErrUnauthorized
			}

			claims, err := auth.ParseToken(parts[1], secret)
			if err != nil {
				return apperrors.ErrUnauthorized
			}

			c.Set(ContextUserID, claims.UserID)
			c.Set(ContextRole, claims.Role)
			return next(c)
		}
	}
}

func RequireRoles(roles ...string) echo.MiddlewareFunc {
	allowed := make(map[string]struct{}, len(roles))
	for _, role := range roles {
		allowed[role] = struct{}{}
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			role, ok := c.Get(ContextRole).(string)
			if !ok {
				return apperrors.ErrUnauthorized
			}

			if _, ok := allowed[role]; !ok {
				return apperrors.ErrForbidden
			}

			return next(c)
		}
	}
}

func UserIDFromContext(c echo.Context) (uint, bool) {
	userID, ok := c.Get(ContextUserID).(uint)
	return userID, ok
}

func RoleFromContext(c echo.Context) (string, bool) {
	role, ok := c.Get(ContextRole).(string)
	return role, ok
}
