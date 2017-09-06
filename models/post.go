package models

import (
	"time"
)

// Post 文章
type Post struct {
	ID        uint       `structs:"-" gorm:"primary_key"`
	CreatedAt time.Time  `structs:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt time.Time  `structs:"updated_at,omitempty" json:"updated_at,omitempty"`
	DeletedAt *time.Time `structs:"-" sql:"index"`
	// Post info
	UserID   uint      `structs:"user_id" json:"user_id" gorm:"index"`                    // 文章发布者的ID, User 的外键
	Title    string    `structs:"title" json:"title" gorm:"column:title"`                 // 文章标题
	Tag      string    `structs:"tag,omitempty" json:"tags,omitempty" gorm:"column:tags"` // 文章标签, 格式: tag1,tag2,tag3,tag4,...
	Author   string    `structs:"author" json:"author" gorm:"column:author"`              // 作者名
	Article  string    `structs:"article" json:"article" gorm:"column:article"`           // 文章内容
	Comments []Comment `structs:"-" json:"comments,omitempty"`                            // 文章评论
}

func (p *Post) TableName() string {
	return "posts"
}

// Comment 文章评论
type Comment struct {
	ID        uint       `structs:"-" gorm:"primary_key"`
	CreatedAt time.Time  `structs:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt time.Time  `structs:"updated_at,omitempty" json:"updated_at,omitempty"`
	DeletedAt *time.Time `structs:"-" sql:"index"`
	// Comment info
	PostID  uint   `structs:"post_id" json:"post_id" gorm:"index"`          // 评论的指定文章ID
	UserID  uint   `structs:"user_id" json:"user_id" gorm:"index"`          // 评论做针对的用户
	Comment string `structs:"comment" json:"comment" gorm:"column:article"` // 评论内容
}

func (c *Comment) TableName() string {
	return "comments"
}
