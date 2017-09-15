package middleware

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/wuxinwei/goforum/conf"
	"github.com/wuxinwei/goforum/models"
	"github.com/wuxinwei/goforum/utils"
	"golang.org/x/crypto/bcrypt"
)

// CipherPassword is action that encrypt password by bcrypt method
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

// StoreSession store the session data
func StoreSession(logger echo.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookies := c.Request().Cookies()
			for _, cookie := range cookies {
				if cookie != nil {
					logger.Warnf("cookie: %#v", cookie)
				}
			}
			sess, err := session.Get(string(conf.SessionSecret), c)
			if err != nil {
				logger.Errorf("Get session failed: %s", err)
			}
			sess.Options = &sessions.Options{
				Path:   c.Path(),
				MaxAge: 7200,
				Secure: true,
			}
			if err := sess.Save(c.Request(), c.Response().Writer); err != nil {
				logger.Errorf("Save session failed: %s", err)
			}
			return next(c)
		}
	}
}

// Auth is action to verify identity of user
func Auth(logger echo.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sess, err := session.Get("session", c)
			if err != nil {
				return c.Render(http.StatusInternalServerError, "errors/500", utils.ErrInternal(err))
			}
			if _, ok := sess.Values["token"]; ok {
				return next(c)
			} else {
				return c.Render(http.StatusForbidden, "errors/403", nil)
			}
		}
	}
}
