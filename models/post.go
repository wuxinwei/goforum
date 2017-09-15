package models

import (
	"time"
)

// Post 文章
type Post struct {
	ID        uint       `json:"id,omitempty" gorm:"primary_key"`
	CreatedAt time.Time  `json:"created_at,omitempty"`
	UpdatedAt time.Time  `json:"updated_at,omitempty"`
	DeletedAt *time.Time `sql:"index"`
	// Post Info
	UserID  uint   `json:"user_id" gorm:"index" sql:"not null"` // author user ID, foreign key
	Title   string `json:"title" gorm:"column:title" sql:"not null"`
	Author  string `json:"author" gorm:"column:author" sql:"not null"`
	Content string `json:"content" gorm:"column:Content;type:TEXT" sql:"not null"`
	Views   uint   `json:"views,omitempty" gorm:"column:views" sql:"not null"`
	Replies uint   `json:"replies,omitempty" gorm:"column:replies" sql:"not null"`
	// Mapping relation foreign key
	Comments []Comment `json:"comments,omitempty" gorm:"ForeignKey:PostID;AssociationForeignKey:ID"`              // post comment
	Tags     []Tag     `json:"tags,omitempty" gorm:"many2many:post_tag;ForeignKey:id;AssociationForeignKey:name"` // post has and belong to many tags
	//Tags []Tag `json:"tags,omitempty" gorm:"many2many:post_tag;"` // post has and belong to many tags
}

func (p *Post) TableName() string {
	return "posts"
}
