package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func forbiddenResponse(c echo.Context) error {
	return c.JSON(http.StatusForbidden, map[string]interface{}{"message": http.StatusText(http.StatusForbidden)})
}

func JWTMiddleware(jwtSign string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if strings.Contains(c.Request().URL.Path, "/login") {
				return next(c)
			}

			signature := strings.Split(c.Request().Header.Get("Authorization"), " ")
			if len(signature) < 2 {
				return forbiddenResponse(c)
			}
			if signature[0] != "Bearer" {
				return forbiddenResponse(c)
			}

			claim := jwt.MapClaims{}
			token, err := jwt.ParseWithClaims(signature[1], claim, func(t *jwt.Token) (interface{}, error) {
				_, ok := t.Method.(*jwt.SigningMethodHMAC)
				if !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
				}

				return []byte(jwtSign), nil
			})
			if err != nil {
				return forbiddenResponse(c)
			}

			method, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok || method != jwt.SigningMethodHS256 {
				return forbiddenResponse(c)
			}

			expAt, err := claim.GetExpirationTime()
			if err != nil {
				return forbiddenResponse(c)
			}

			if time.Now().After(expAt.Time) {
				return forbiddenResponse(c)
			}

			userID, _ := claim["id"].(string)
			role, _ := claim["role"].(string)
			c.Set("id", userID)
			c.Set("role", role)

			return next(c)
		}
	}
}

func stringSliceContains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}

	return false
}

// func RBACMiddleware(roles []string) echo.MiddlewareFunc {
// 	return func(next echo.HandlerFunc) echo.HandlerFunc {
// 		return func(c echo.Context) error {
// 			role, _ := c.Get("role").(string)
// 			if role == "superadmin" {
// 				return next(c)
// 			}

// 			if stringSliceContains(roles, role) {
// 				return next(c)
// 			}

// 			return forbiddenResponse(c)
// 		}
// 	}
// }

func ACLMiddleware(rolesMap map[string]bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			role, _ := c.Get("role").(string)
			if role == "superadmin" {
				return next(c)
			}

			if rolesMap[role] {
				return next(c)
			}

			return forbiddenResponse(c)
		}
	}
}
