package controllers

import (
	"net/http"

	"github.com/fatih/structs"
	"github.com/fpay/gopress"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"github.com/wuxinwei/goforum/middleware"
	"github.com/wuxinwei/goforum/models"
	"github.com/wuxinwei/goforum/utils"
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
	app.GET("/test", c.CreateRender("errors/500"))

	app.GET("/", c.CreateRender("user/login"))
	app.GET("/user", c.CreateRender("user/login"))
	app.GET("/user/login", c.CreateRender("user/login"))
	app.POST("/user/login", c.Login, middleware.CipherPassword())
	app.GET("/user/register", c.CreateRender("user/register"))
	app.POST("/user/register", c.Register, middleware.CipherPassword())
	app.GET("/user/profile", c.Profile, middleware.Auth())
	app.POST("/user/profile", c.SetProfile, middleware.Auth())
}

// CreateRender got a page render by specific path
func (c *UserController) CreateRender(renderPath string) func(context gopress.Context) error {
	return func(ctx gopress.Context) error {
		err := ctx.Render(http.StatusOK, renderPath, "Error Code")
		if err != nil {
			logrus.Errorf("Render Error: %s", err)
		}
		return err
	}
}

// Login Action
func (c *UserController) Login(ctx gopress.Context) error {
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

	// TODO, 登录后需要保存 session 信息, 之后所有大部分读取操作均从 redis 中读取
	if err := db.Where("username = ? AND password = ?", user.UserName, user.Password).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ctx.Render(http.StatusInternalServerError, "user/login", nil)
		}
		return ctx.Render(http.StatusBadRequest, "errors/500", map[string]interface{}{
			"error": err.Error(),
		})
	}
	// TODO: redirect 需要重新考虑
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
	if err := db.Where("username = ? ", user.UserName).First(&models.User{}).Error; err == nil {
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
		UserName: username,
	}
	db, err := utils.GetDB(c.app)
	if err != nil {
		return ctx.Render(http.StatusInternalServerError, "errors/500", map[string]interface{}{
			"error": err.Error(),
		})
	}
	if err := db.Where("username = ?", user.UserName).First(&user).Error; err != nil {
		return ctx.Render(http.StatusInternalServerError, "errors/500", map[string]interface{}{
			"error": err.Error(),
		})
	}
	// clear password
	user.Password = ""
	return ctx.Render(http.StatusOK, "/user/profile", structs.Map(&user))
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
	if err := db.Model(&user).Where("username = ?", user.UserName).Update(userMap).Error; err != nil {
		return ctx.Render(http.StatusInternalServerError, "errors/500", map[string]interface{}{
			"error": err.Error(),
		})
	}
	return ctx.Redirect(http.StatusPermanentRedirect, "https://www.baidu.com")
}
