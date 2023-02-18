package testhelpers

import (
	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/models"
	"gorm.io/gorm"
)

// SimpleFakeDatabase
type SimpleFakeDatabase struct {
}

func (f *SimpleFakeDatabase) CreateUser(user *models.User) *gorm.DB {
	return nil
}

func (f *SimpleFakeDatabase) FindUserByEmail(email string, model *models.User) *gorm.DB {
	return &gorm.DB{}
}

// EmailTakenFakeDataBase
type EmailTakenFakeDatabase struct {
}

func (f *EmailTakenFakeDatabase) CreateUser(user *models.User) *gorm.DB {
	return nil
}

func (f *EmailTakenFakeDatabase) FindUserByEmail(email string, model *models.User) *gorm.DB {
	return &gorm.DB{RowsAffected: 1}
}
