package controllers

import (
	"net/http"

	"github.com/fpay/gopress"
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

	app.GET("/post", nil)          // 获取一篇文章
	app.GET("/posts", nil)         // 获取文章列表
	app.POST("/post", nil)         // 发表/更新文章
	app.POST("/post/comment", nil) // 发表/修改评论
}

// PostArticle Action
func (c *PostsController) PostArticle(ctx gopress.Context) error {
	return ctx.Render(http.StatusOK, "posts/sample", nil)
}

// GetArticle Action
func (c *PostsController) GetArticles(ctx gopress.Context) error {
	return ctx.Render(http.StatusOK, "posts/sample", nil)
}

// PostComment Action
func (c *PostsController) PostComment(ctx gopress.Context) error {
	return nil
}

// GetComment Action
func (c *PostsController) GetComments(ctx gopress.Context) error {
	return nil
}
