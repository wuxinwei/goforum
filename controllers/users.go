package controllers

import (
	"net/http"

	"github.com/fatih/structs"
	"github.com/fpay/gopress"
	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo-contrib/session"
	"github.com/wuxinwei/goforum/middleware"
	"github.com/wuxinwei/goforum/models"
	"github.com/wuxinwei/goforum/utils"
	"golang.org/x/crypto/bcrypt"
)

// UserController the controller part of MVC patter
type UserController struct {
	app *gopress.App
}

// NewUserController returns user controller instance.
func NewUserController() *UserController {
	return new(UserController)
}

// RegisterRoutes register routes to app
func (c *UserController) RegisterRoutes(app *gopress.App) {
	c.app = app

	// test
	app.GET("/test", func(ctx gopress.Context) error {
		return ctx.Render(http.StatusOK, "errors/500", nil)
	})

	app.POST("/user/login", c.Login)
	app.POST("/user/logout", func(ctx gopress.Context) error {
		return ctx.Redirect(http.StatusPermanentRedirect, "http://localhost:3000/")
	})
	app.POST("/user/register", c.Register, middleware.CipherPassword())
	app.GET("/user/register", func(ctx gopress.Context) error {
		return ctx.Render(http.StatusOK, "user/register", nil)
	})
	app.POST("/user/profile", c.SetProfile, middleware.Auth(app.Logger))
	app.GET("/user/profile", c.Profile, middleware.Auth(app.Logger))
}

// Login Action
func (c *UserController) Login(ctx gopress.Context) error {
	sess, err := session.Get("session_id", ctx)
	if err != nil {
		c.app.Logger.Errorf("Get session failed: %s", err)
	}
	sess.Options = &sessions.Options{
		Path:     ctx.Path(),
		MaxAge:   7200,
		HttpOnly: true,
	}
	if err := sess.Save(ctx.Request(), ctx.Response().Writer); err != nil {
		c.app.Logger.Errorf("Save session failed: %s", err)
	}
	var user models.User
	if err := ctx.Bind(&user); err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}
	c.app.Logger.Debugf("Request Content Type: %#v", ctx.Request().Header.Get("Content-Type"))
	c.app.Logger.Debugf("Request Content: %#v", user)
	rawPw := user.Password
	db, err := utils.GetDB(c.app)
	if err != nil {
		return ctx.Render(http.StatusInternalServerError, "errors/500", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// TODO, 登录后需要保存 session 信息, 之后所有大部分读取操作均从 redis 中读取
	if err := db.Where("username = ?", user.Username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ctx.String(http.StatusBadRequest, err.Error())
		}
		return ctx.Render(http.StatusInternalServerError, "errors/500", map[string]interface{}{
			"error": err.Error(),
		})
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(rawPw)); err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}
	return ctx.Redirect(http.StatusPermanentRedirect, "https://127.0.0.1:8080/posts")
}

// Register Action
func (c *UserController) Register(ctx gopress.Context) error {
	var user models.User
	if err := ctx.Bind(&user); err != nil {
		return ctx.Render(http.StatusBadRequest, "errors/500", map[string]interface{}{
			"error": err.Error(),
		})
	}
	db, err := utils.GetDB(c.app)
	if err != nil {
		return ctx.Render(http.StatusBadRequest, "errors/500", map[string]interface{}{
			"error": err.Error(),
		})
	}
	// check cache and db
	if err := db.Where("username = ? ", user.Username).First(&models.User{}).Error; err == nil {
		return ctx.Render(http.StatusOK, "user/login", nil)
	} else if err != nil && err != gorm.ErrRecordNotFound {
		return ctx.Render(http.StatusInternalServerError, "errors/500", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// create new user
	// TODO: 更新 redis/ 更新 session
	if err := db.Create(&user).Error; err != nil {
		return ctx.Render(http.StatusInternalServerError, "errors/500", map[string]interface{}{
			"error": err.Error(),
		})
	}
	c.app.Logger.Errorf("User: %#v", user)
	return ctx.Redirect(http.StatusPermanentRedirect, "https://www.baidu.com")
}

// Profile Action
func (c *UserController) Profile(ctx gopress.Context) error {
	// TODO: 应该在 Redis 中获取, 如果过期则需要重新更新 Redis
	username := ctx.QueryParam("username")
	if username == "" {
		return ctx.Render(http.StatusBadRequest, "errors/400", map[string]interface{}{
			"error": "invalid username",
		})
	}
	user := models.User{
		Username: username,
	}
	db, err := utils.GetDB(c.app)
	if err != nil {
		return ctx.Render(http.StatusInternalServerError, "errors/500", map[string]interface{}{
			"error": err.Error(),
		})
	}
	if err := db.Where("username = ?", user.Username).First(&user).Error; err != nil {
		return ctx.Render(http.StatusInternalServerError, "errors/500", map[string]interface{}{
			"error": err.Error(),
		})
	}
	// clear password
	user.Password = ""
	return ctx.Render(http.StatusOK, "/user/profile", &user)
}

// SetProfile Action
func (c *UserController) SetProfile(ctx gopress.Context) error {
	var user models.User
	if err := ctx.Bind(&user); err != nil {
		return ctx.Render(http.StatusBadRequest, "errors/500", map[string]interface{}{
			"error": err.Error(),
		})
	}
	db, err := utils.GetDB(c.app)
	if err != nil {
		return ctx.Render(http.StatusInternalServerError, "errors/500", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// set profile transaction
	user.Password = ""
	userMap := structs.Map(&user)
	if err := db.Model(&user).Where("username = ?", user.Username).Update(userMap).Error; err != nil {
		return ctx.Render(http.StatusInternalServerError, "errors/500", map[string]interface{}{
			"error": err.Error(),
		})
	}
	return ctx.Redirect(http.StatusPermanentRedirect, "https://www.baidu.com")
}
