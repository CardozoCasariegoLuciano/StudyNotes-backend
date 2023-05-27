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
	return st.db.Where("author = ?", userID).Find(model)
}

func (st *Storage) GetNoteByID(noteID int, model *dbmodels.Note) *gorm.DB {
	return st.db.Where("id = ?", noteID).First(model)
}

func (st *Storage) DeleteNoteByID(noteID int, model *dbmodels.Note) *gorm.DB {
	return st.db.Unscoped().Delete(model, noteID)
}
