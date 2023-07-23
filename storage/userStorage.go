package storage

import (
	dbmodels "github.com/CardozoCasariegoLuciano/StudyNotes-backend/models/dbModels"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

/*
Nota
Cuando actualices agregues un nuevo metodo actualiza
la interfaz Istorage
models/apimodels/Istorage
*/
func (st *Storage) FindUserByEmail(email string, model *dbmodels.User) *gorm.DB {
	return st.db.Where("email = ?", email).First(model)
}

func (st *Storage) ComparePasswords(hashedPass string, bodyPass string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(bodyPass))
}

func (st *Storage) GetUserByID(userID int, model *dbmodels.User) *gorm.DB {
	return st.db.
		Select("name, email, image, role, id").
		Where("id = ?", userID).
		First(model)
}
