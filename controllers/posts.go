package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/fpay/gopress"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/wuxinwei/goforum/models"
	"github.com/wuxinwei/goforum/utils"
)

// PostsController
type PostsController struct {
	app *gopress.App
}

// NewPostsController returns posts controller instance.
func NewPostsController() *PostsController {
	return new(PostsController)
}

// RegisterRoutes registes routes to app
func (c *PostsController) RegisterRoutes(app *gopress.App) {
	c.app = app
	app.GET("/", c.Posts, func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Request().Header.Set("Index-Page", "true")
			return next(c)
		}
	})
	app.GET("/post/:id", c.Post)         // 获取一篇文章
	app.POST("/post", c.Publish)         // 发表/更新文章
	app.GET("/post/list", c.Posts)       // 获取文章列表
	app.POST("/post/comment", c.Comment) // 发表/修改评论
}

// Post Action is that get a article
func (c *PostsController) Post(ctx gopress.Context) error {
	id := ctx.Param("id")
	db, err := utils.GetDB(c.app)
	if err != nil {
		return ctx.Render(http.StatusInternalServerError, "errors/500", errInternal(err))
	}
	var post models.Post
	if err := db.Where("id = ?", id).First(&post).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return ctx.Render(http.StatusInternalServerError, "errors/500", errInternal(err))
		} else {
			return ctx.Render(http.StatusInternalServerError, "errors/404", errRequest(err))
		}
	}
	return ctx.Render(http.StatusOK, "post/post", &post)
}

// Publish Action is that publish a article
func (c *PostsController) Publish(ctx gopress.Context) error {
	var post models.Post
	if err := ctx.Bind(&post); err != nil {
		return ctx.Render(http.StatusBadRequest, "errors/400", errInternal(err))
	}
	db, err := utils.GetDB(c.app)
	if err != nil {
		return ctx.Render(http.StatusInternalServerError, "errors/500", errInternal(err))
	} else if post.ID != 0 {
		// update post
		if err := db.Model(&post).Where("id = ?", post.ID).Updates(&post).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				return ctx.Render(http.StatusNotFound, "errors/404", errInternal(err))
			} else {
				return ctx.Render(http.StatusInternalServerError, "errors/500", errInternal(err))
			}
		}
	} else if err := db.Create(&post).Error; err != nil {
		// create post
		return ctx.Render(http.StatusInternalServerError, "errors/500", errInternal(err))
	}

	// TODO: 考虑将更新用户金币利用ＭＱ缓冲，改成异步操作
	// update user's coin count by transaction
	var user models.User
	t := db.Begin()
	if err := t.Model(&user).Where("id = ?", post.UserID).First(&user).Error; err != nil {
		t.Rollback()
		return ctx.Render(http.StatusInternalServerError, "errors/500", errInternal(err))
	}
	if err := t.Model(&user).Where("id = ?", user.ID).Update("coin", user.Coin+1).Error; err != nil {
		t.Rollback()
		return ctx.Render(http.StatusInternalServerError, "errors/500", errInternal(err))
	}
	if err := t.Commit().Error; err != nil {
		return ctx.Render(http.StatusInternalServerError, "errors/500", errInternal(err))
	}
	return ctx.Redirect(http.StatusPermanentRedirect, fmt.Sprintf("http://localhost:3000/post/%d", post.ID))
}

// Posts Action is that get a list of article
func (c *PostsController) Posts(ctx gopress.Context) error {
	var page struct {
		Tag       string `json:"tag,omitempty"`
		PageIndex int    `json:"page_index"`
	}
	if "true" != ctx.Request().Header.Get("Index-Page") {
		page.Tag = ctx.QueryParam("tag")
		if pageIndex, err := strconv.Atoi(ctx.QueryParam("page_index")); err != nil {
			return ctx.Render(http.StatusBadRequest, "errors/400", map[string]interface{}{
				"error": "invalid page index: " + ctx.QueryParam("page_index"),
			})
		} else if pageIndex < 1 {
			return ctx.Render(http.StatusBadRequest, "errors/400", map[string]interface{}{
				"error": "invalid page index: " + ctx.QueryParam("page_index"),
			})
		} else {
			page.PageIndex = pageIndex
		}
	} else {
		page.PageIndex = 0
	}

	var posts []models.Post
	if db, err := utils.GetDB(c.app); err != nil {
		return ctx.Render(http.StatusInternalServerError, "errors/500", errInternal(err))
	} else if err := db.Order("updated_at DESC").Offset((page.PageIndex - 1) * models.PageSize).Limit(models.PageSize).Find(&posts).Error; err != nil {
		return ctx.Render(http.StatusInternalServerError, "errors/500", errInternal(err))
	}
	data := map[string]interface{}{
		"posts": &posts,
	}
	return ctx.Render(http.StatusOK, "post/list", data)
}

// Comment Action is that issue or update a comment on specific post
func (c *PostsController) Comment(ctx gopress.Context) error {
	var comment models.Comment
	if err := ctx.Bind(&comment); err != nil {
		return ctx.Render(http.StatusInternalServerError, "errors/500", map[string]interface{}{
			"error": err.Error(),
		})
	}

	db, err := utils.GetDB(c.app)
	if err != nil {
		return ctx.Render(http.StatusInternalServerError, "errors/500", map[string]interface{}{
			"error": err.Error(),
		})
	}
	t := db.Begin()
	if err := t.Model(&models.Post{}).Where("id = ?", comment.PostID).First(&models.Post{}).Error; err != nil {
		t.Rollback()
		if err != gorm.ErrRecordNotFound {
			return ctx.Render(http.StatusNotFound, "errors/404", errRequest(err))
		} else {
			return ctx.Render(http.StatusInternalServerError, "errors/500", errInternal(err))
		}
	}
	if comment.ID != 0 {
		// prepare update transaction
		if err := t.Model(&comment).Where("id = ? AND post_id = ?", comment.ID, comment.PostID).Updates(&comment).Error; err != nil {
			t.Rollback()
			if err != gorm.ErrRecordNotFound {
				return ctx.Render(http.StatusNotFound, "errors/404", errRequest(err))
			} else {
				return ctx.Render(http.StatusInternalServerError, "errors/500", errInternal(err))
			}
		}
	} else {
		// prepare create new comment transaction
		if err := t.Model(&comment).Create(&comment).Error; err != nil {
			t.Rollback()
			return ctx.Render(http.StatusInternalServerError, "errors/500", errInternal(err))
		}
		// TODO: 将更新金币的操作作为消息发送MQ, 更新金币操作将为异步
		// update user's coin count
		var user models.User
		if err := t.Model(&user).Where("id = ?", comment.UserID).First(&user).Error; err != nil {
			t.Rollback()
			if err != gorm.ErrRecordNotFound {
				return ctx.Render(http.StatusNotFound, "errors/404", errRequest(err))
			} else {
				return ctx.Render(http.StatusInternalServerError, "errors/500", errInternal(err))
			}
		}
		if err := t.Model(&user).Where("id = ?", user.UserID).Update("coin", user.Coin+1).Error; err != nil {
			t.Rollback()
			return ctx.Render(http.StatusInternalServerError, "errors/500", errInternal(err))
		}
	}
	if err := t.Commit().Error; err != nil {
		return ctx.Render(http.StatusInternalServerError, "errors/500", errInternal(err))
	}
	return ctx.Redirect(http.StatusPermanentRedirect, "https://localhost:3000/post/"+strconv.Itoa(int(comment.PostID)))
}
