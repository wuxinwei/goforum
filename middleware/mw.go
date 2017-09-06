package middleware

import (
	"net/http"

	"encoding/json"

	"io/ioutil"

	"bytes"

	"github.com/labstack/echo"
	"github.com/wuxinwei/goforum/models"
	"golang.org/x/crypto/bcrypt"
)

// CipherPassword is action that encrypt password by bcrypt cipher method
func CipherPassword() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var user models.User
			req := c.Request()
			if err := c.Bind(&user); err != nil {
				return c.Render(http.StatusBadRequest, "user/400", map[string]interface{}{
					"msg": err.Error(),
				})
			}

			passwordEncrypted, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
			if err != nil {
				return c.Render(http.StatusInternalServerError, "user/500", map[string]interface{}{
					"error": err.Error(),
				})
			}

			user.Password = string(passwordEncrypted)
			newBody, err := json.Marshal(&user)
			if err != nil {
				return c.Render(http.StatusInternalServerError, "user/500", map[string]interface{}{
					"error": err.Error(),
				})
			}
			req.Body = ioutil.NopCloser(bytes.NewReader(newBody))
			req.ContentLength = int64(len(newBody))
			c.SetRequest(req)

			return next(c)
		}
	}
}

// Auth is action to verify identity of user
func Auth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// TODO: 验证session
			return next(c)
		}
	}
}
