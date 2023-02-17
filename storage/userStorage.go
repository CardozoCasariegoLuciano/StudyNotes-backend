package storage

import (
	"github.com/CardozoCasariegoLuciano/StudyNotes-backend/models"
	"gorm.io/gorm"
)

func (st *Storage) CreateUser(user *models.User) *gorm.DB {
	return st.db.Save(user)
}

func (st *Storage) FindUserByEmail(email string, model *models.User) *gorm.DB {
	return st.db.Where("email = ?", email).First(model)
}
