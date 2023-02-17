package models

import "gorm.io/gorm"

type Istorage interface {
	CreateUser(user *User) *gorm.DB
	FindUserByEmail(email string, model *User) *gorm.DB
}
