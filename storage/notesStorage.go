package storage

import (
	dbmodels "github.com/CardozoCasariegoLuciano/StudyNotes-backend/models/dbModels"
	"gorm.io/gorm"
)

/*
Nota
Cuando actualices agregues un nuevo metodo actualiza
la interfaz Istorage
models/apimodels/Istorage
*/
func (st *Storage) GetAllNotes(userID int, model *dbmodels.Notes) *gorm.DB {
	return st.db.Where("user_id = ?", userID).Find(model)
}
