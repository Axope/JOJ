package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UUID     string `json:"uuid" gorm:"type:CHAR(36); not null; unique_index:idx_uuid; comment:'uuid'"`
	Username string `json:"username" form:"username" binding:"required" gorm:"type:VARCHAR(50); unique; not null"`
	Password string `json:"password" form:"password" binding:"required" gorm:"type:VARCHAR(50); not null"`
	Admin    int    `json:"admin" gorm:"type:INT; not null"`
}

func NewUser(uuid, username, password string) *User {
	return &User{
		UUID:     uuid,
		Username: username,
		Password: password,
		Admin:    0,
	}
}

func (u *User) TableName() string {
	return "users"
}
