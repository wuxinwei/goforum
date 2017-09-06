package models

import (
	"database/sql/driver"
	"time"
)

// GenderType
type GenderType string

// Scan override method
func (u *GenderType) Scan(value interface{}) error { *u = GenderType(value.([]byte)); return nil }

// Value override method
func (u GenderType) Value() (driver.Value, error) { return []byte(u), nil }

const (
	male   GenderType = "male"
	female GenderType = "female"
)

type User struct {
	ID        uint       `structs:"-" gorm:"primary_key"`
	CreatedAt time.Time  `structs:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt time.Time  `structs:"updated_at,omitempty" json:"updated_at,omitempty"`
	DeletedAt *time.Time `structs:"-" sql:"index"`
	// Login info
	UserName string `structs:"username" json:"username" gorm:"not null;column:username;unique_index"`
	Password string `structs:"password" json:"password" gorm:"not null;column:password;type:CHAR(60)"`
	// Profile
	Name     string     `structs:"name,omitempty" json:"name,omitempty" gorm:"column:name"`
	Gender   GenderType `structs:"gender,omitempty" json:"gender,omitempty" gorm:"column:gender;type:ENUM('male', 'female')"`
	Birthday time.Time  `structs:"birthday,omitempty" json:"birthday,omitempty" gorm:"column:birthday;type:DATE"`
	Country  string     `structs:"country,omitempty" json:"country,omitempty" gorm:"column:country"`
	Province string     `structs:"province,omitempty" json:"province,omitempty" gorm:"column:province"`
	City     string     `structs:"city,omitempty" json:"city,omitempty" gorm:"column:city"`
	Coin     uint       `structs:"coin,omitempty" json:"coin,omitempty" gorm:"column:coin"`
	// Article
	Posts []Post `structs:"-" json:"posts,omitempty"`
}

func (u *User) TableName() string {
	return "users"
}
