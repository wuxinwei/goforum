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
	Male   GenderType = "male"
	Female GenderType = "female"
)

type User struct {
	ID        uint       `gorm:"primary_key"`
	CreatedAt time.Time  `json:"created_at,omitempty"`
	UpdatedAt time.Time  `json:"updated_at,omitempty"`
	DeletedAt *time.Time `sql:"index"`
	// Login info
	Username string `json:"username" gorm:"column:username;unique_index" sql:"not null"`
	Password string `json:"password" gorm:"column:password;type:CHAR(60)" sql:"not null"`
	// Profile
	Name     string     `json:"name,omitempty" gorm:"column:name"`
	Gender   GenderType `json:"gender,omitempty" gorm:"column:gender;type:ENUM('male', 'female')"`
	Birthday time.Time  `json:"birthday,omitempty" gorm:"column:birthday;type:DATE"`
	Country  string     `json:"country,omitempty" gorm:"column:country"`
	Province string     `json:"province,omitempty" gorm:"column:province"`
	City     string     `json:"city,omitempty" gorm:"column:city"`
	Coin     uint       `json:"coin,omitempty" gorm:"column:coin"`
	// Post && Comment
	Posts    []Post    `json:"posts,omitempty" gorm:"ForeignKey:UserID;AssociationForeignKey:ID"`
	Comments []Comment `json:"comments,omitempty" gorm:"ForeignKey:UserID;AssociationForeignKey:ID"`
}

func (u *User) TableName() string {
	return "users"
}
