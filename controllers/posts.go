package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/fpay/gopress"
	"github.com/wuxinwei/goforum/conf"
	"github.com/wuxinwei/goforum/middleware"
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

	app.GET("/", c.Posts)
	app.GET("/post/:id", c.Post)                                      // 获取一篇文章
	app.POST("/post", c.Publish, middleware.Auth(app.Logger))         // 发表/更新文章
	app.GET("/posts", c.Posts)                                        // 获取文章列表
	app.POST("/post/comment", c.Comment, middleware.Auth(app.Logger)) // 发表/修改评论
	app.GET("/post/search", c.TextSearch)                             // 全文检索
}

// Post Action is that get a article
func (c *PostsController) Post(ctx gopress.Context) error {
	id := ctx.Param("id")
	db, err := utils.GetDB(c.app)
	if err != nil {
		return ctx.Render(http.StatusInternalServerError, "errors/500", utils.ErrInternal(err))
	}
	var post models.Post
	if db.Where("id = ?", id).First(&post).RecordNotFound() {
		return ctx.Render(http.StatusInternalServerError, "errors/404", utils.ErrRequest(err))
	}
	return ctx.Render(http.StatusOK, "post/post", &post)
}

// Publish Action is that publish a article
func (c *PostsController) Publish(ctx gopress.Context) error {
	var post models.Post
	if err := ctx.Bind(&post); err != nil {
		return ctx.Render(http.StatusBadRequest, "errors/400", utils.ErrInternal(err))
	}
	db, err := utils.GetDB(c.app)
	if err != nil {
		return ctx.Render(http.StatusInternalServerError, "errors/500", utils.ErrInternal(err))
	}
	t := db.Begin()
	if post.ID != 0 {
		// update post
		if t.Model(&models.Post{}).Where("id = ?", post.ID).First(&models.Post{}).RecordNotFound() {
			t.Rollback()
			return ctx.Render(http.StatusNotFound, "errors/404", utils.ErrInternal(err))
		}
		if err := t.Model(&post).Where("id = ?", post.ID).Updates(&post).Error; err != nil {
			t.Rollback()
			return ctx.Render(http.StatusInternalServerError, "errors/500", utils.ErrInternal(err))
		}
		// update tag
		//for _, tag := range post.Tags {
		//}

		// update mapping relation between tag and post
	} else if err := t.Create(&post).Error; err != nil {
		t.Rollback()
		// create post
		return ctx.Render(http.StatusInternalServerError, "errors/500", utils.ErrInternal(err))
	}

	// create tag -> post mapping relation
	//for _, tag := range post.Tags {
	//	tagMap := models.TagMap{
	//		Tag:    tag.Name,
	//		PostID: post.ID,
	//	}
	//	if err := t.Create(&tagMap).Error; err != nil {
	//		t.Rollback()
	//		return ctx.Render(http.StatusInternalServerError, "errors/500", utils.ErrInternal(err))
	//	}
	//}

	// TODO: 考虑将更新用户金币利用ＭＱ缓冲，改成异步操作
	// update user's coin count by transaction
	var user models.User
	if err := t.Model(&user).Where("id = ?", post.UserID).First(&user).Error; err != nil {
		t.Rollback()
		return ctx.Render(http.StatusInternalServerError, "errors/500", utils.ErrInternal(err))
	}
	if err := t.Model(&user).Where("id = ?", user.ID).Update("coin", user.Coin+1).Error; err != nil {
		t.Rollback()
		return ctx.Render(http.StatusInternalServerError, "errors/500", utils.ErrInternal(err))
	}
	t.Commit()
	return ctx.Redirect(http.StatusPermanentRedirect, fmt.Sprintf("http://localhost:3000/post/%d", post.ID))
}

// Posts Action is that get a list of article
func (c *PostsController) Posts(ctx gopress.Context) error {
	var page struct {
		Tag       string `json:"tag,omitempty"`
		PageIndex int    `json:"page_index"`
	}
	page.Tag = ctx.QueryParam("tag")
	pageIndexStr := ctx.QueryParam("page_index")
	if pageIndex, err := strconv.Atoi(pageIndexStr); err != nil || pageIndex < 1 {
		return ctx.Render(http.StatusBadRequest, "errors/400", utils.ErrRequest(errors.New("invalid page index: "+pageIndexStr)))
	} else {
		page.PageIndex = pageIndex
	}

	var posts []models.Post
	if db, err := utils.GetDB(c.app); err != nil {
		return ctx.Render(http.StatusInternalServerError, "errors/500", utils.ErrInternal(err))
	} else if err := db.Order("updated_at DESC").Offset((page.PageIndex - 1) * conf.PageSize).Limit(conf.PageSize).Find(&posts).Error; err != nil {
		return ctx.Render(http.StatusInternalServerError, "errors/500", utils.ErrInternal(err))
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
		return ctx.Render(http.StatusInternalServerError, "errors/500", utils.ErrInternal(err))
	}

	db, err := utils.GetDB(c.app)
	if err != nil {
		return ctx.Render(http.StatusInternalServerError, "errors/500", utils.ErrInternal(err))
	}
	t := db.Begin()
	if t.Model(&models.Post{}).Where("id = ?", comment.PostID).First(&models.Post{}).RecordNotFound() {
		t.Rollback()
		return ctx.Render(http.StatusNotFound, "errors/404", utils.ErrRequest(errors.New("this post does not exist forever")))
	}
	if comment.ID != 0 {
		// prepare update transaction
		if err := t.Model(&comment).Where("id = ? AND post_id = ?", comment.ID, comment.PostID).Updates(&comment).Error; err != nil {
			t.Rollback()
			return ctx.Render(http.StatusInternalServerError, "errors/500", utils.ErrInternal(err))
		}
	} else {
		// prepare create new comment transaction
		if err := t.Model(&comment).Create(&comment).Error; err != nil {
			t.Rollback()
			return ctx.Render(http.StatusInternalServerError, "errors/500", utils.ErrInternal(err))
		}
		// TODO: 将更新金币的操作作为消息发送MQ, 更新金币操作将为异步
		// update user's coin count
		var user models.User
		if t.Model(&user).Where("id = ?", comment.UserID).First(&user).RecordNotFound() {
			t.Rollback()
			return ctx.Render(http.StatusNotFound, "errors/404", utils.ErrRequest(err))
		}
		if err := t.Model(&user).Where("id = ?", user.ID).Update("coin", user.Coin+1).Error; err != nil {
			t.Rollback()
			return ctx.Render(http.StatusInternalServerError, "errors/500", utils.ErrInternal(err))
		}
	}
	t.Commit()
	return ctx.Redirect(http.StatusPermanentRedirect, fmt.Sprint("http://localhost:3000/post/%#d", comment.PostID))
}

// Tag Action is that filter post by specific tag
// TODO: 支持３个Tag检索
func (c *PostsController) Tag(ctx gopress.Context) error {
	tag := ctx.QueryParam("tag")
	if tag == "" {
		return ctx.Render(http.StatusNotFound, "errors/400", utils.ErrRequest(errors.New("empty tag")))
	}
	return nil
	db, err := utils.GetDB(c.app)
	if err != nil {
		return ctx.Render(http.StatusInternalServerError, "errors/500", utils.ErrInternal(err))
	}
	var posts []models.Post
	if err := db.Joins("INNER JOIN post_tag ON post_tag.post_id = posts.id AND post_tag.tag_name = ?", tag).Find(&posts).Error; err != nil {
		return ctx.Render(http.StatusInternalServerError, "errors/500", utils.ErrInternal(err))
	}
	data := map[string]interface{}{
		"posts": &posts,
	}
	return ctx.Render(http.StatusOK, "post/list", data)
}

// TextSearch action is text search service
// TODO: 搜索结果多页面
func (c *PostsController) TextSearch(ctx gopress.Context) error {
	var req struct {
		Scope string `json:"scope" valid:"in(title|content);"`
		Value string `json:"value" valid:"-"`
	}
	if err := ctx.Bind(&req); err != nil {
		return ctx.Render(http.StatusBadRequest, "errors/404", utils.ErrRequest(err))
	}
	es, err := utils.GetElastic(c.app)
	if err != nil {
		return ctx.Render(http.StatusInternalServerError, "errors/500", utils.ErrInternal(err))
	}
	posts, err := es.SearchPost(req.Scope, req.Value)
	if err != nil {
		return ctx.Render(http.StatusInternalServerError, "errors/500", utils.ErrInternal(err))
	}
	data := map[string]interface{}{
		"posts": &posts,
	}
	return ctx.Render(http.StatusOK, "post/list", data)
}
