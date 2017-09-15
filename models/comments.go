package models

import "time"

// Comment post comment
type Comment struct {
	ID        uint       `json:"id,omitempty" gorm:"primary_key"`
	CreatedAt time.Time  `json:"created_at,omitempty"`
	UpdatedAt time.Time  `json:"updated_at,omitempty"`
	DeletedAt *time.Time `sql:"index"`
	// comment info
	UserID   uint   `json:"user_id" gorm:"index" sql:"not null"`              // user ID, foreign key
	PostID   uint   `json:"post_id" gorm:"index" sql:"not null"`              // post ID, foreign key
	TargetID uint   `json:"target_id" gorm:"column:target_id" sql:"not null"` // 该评论针对的用户
	Content  string `json:"content" gorm:"not null;column:content;type:TEXT" sql:"not null"`
}

func (c *Comment) TableName() string {
	return "comments"
}
