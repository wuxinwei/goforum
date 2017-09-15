package models

import "time"

// Tag is a sql table for mapping relation between tag and post

type Tag struct {
	ID        uint       `json:"id,omitempty" gorm:"primary_key"`
	CreatedAt time.Time  `json:"created_at,omitempty"`
	UpdatedAt time.Time  `json:"updated_at,omitempty"`
	DeletedAt *time.Time `sql:"index"`

	Name  string `json:"name" gorm:"column:name" sql:"not null"`
	Posts []Post `gorm:"many2many:post_tag;ForeignKey:name;AssociationForeignKey:id"` // post has and belong to many tags
}

func (t *Tag) TableName() string {
	return "tags"
}
