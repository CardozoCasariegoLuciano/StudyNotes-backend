package apimodels

import (
	dbmodels "github.com/CardozoCasariegoLuciano/StudyNotes-backend/models/dbModels"
	"gorm.io/gorm"
)

/*
Nota
Cuando actualices la interfaz Istorage no te olvides de correr el comando:

mockgen -source=models/apiModels/storage.go -destination=./handlers/mocks/IStorageMocks.go

para mockearla y asi no romper los test viejos y poder testeaar nos nuevos cambios
*/

type Istorage interface {
	Create(anyModel interface{}) *gorm.DB
	FindUserByEmail(email string, model *dbmodels.User) *gorm.DB
	ComparePasswords(hashedPass string, bodyPass string) error
	GetAllNotes(userID int, model *dbmodels.Notes) *gorm.DB
	GetNoteByID(noteID int, model *dbmodels.Note) *gorm.DB
	DeleteNoteByID(noteID int, model *dbmodels.Note) *gorm.DB
}
