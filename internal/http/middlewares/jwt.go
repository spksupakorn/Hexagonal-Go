package middlewares

import (
	"dungeons-dragon-service/internal/http/custom"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type JWTMiddleware struct {
	secret []byte
}

func NewJWTMiddleware(secret string) *JWTMiddleware {
	return &JWTMiddleware{secret: []byte(secret)}
}

func (m *JWTMiddleware) Parse(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		auth := c.Request().Header.Get("Authorization")
		if auth == "" {
			// No token provided: proceed as unauthenticated
			return next(c)
		}
		tokenStr := auth
		// Support "Bearer <token>"
		if len(auth) > 7 && auth[:7] == "Bearer " {
			tokenStr = auth[7:]
		}

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return m.secret, nil
		})
		if err != nil || !token.Valid {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			// Set user in context
			if sub, ok := claims["sub"]; ok {
				switch v := sub.(type) {
				case string:
					c.Set("userID", v)
				}
			}
			if role, ok := claims["role"].(string); ok {
				c.Set("role", role)
			}
		}
		return next(c)
	}
}

func RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().Method == http.MethodOptions {
			return next(c)
		}
		if c.Get("userID") == nil {
			return c.JSON(http.StatusUnauthorized, custom.BuildResponse(custom.Unauthorized, "authentication required"))
		}
		return next(c)
	}
}

func RequireAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().Method == http.MethodOptions {
			return next(c)
		}
		role, _ := c.Get("role").(string)
		if role != "admin" {
			return c.JSON(http.StatusForbidden, custom.BuildResponse(custom.Forbidden, "admin access required"))
		}
		return next(c)
	}
}

func GetUserID(c echo.Context) (string, bool) {
	id, ok := c.Get("userID").(string)
	return id, ok
}

func IsAuthenticated(c echo.Context) bool {
	_, ok := c.Get("userID").(string)
	return ok
}
