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
	Save(anyModel interface{}) *gorm.DB
	GetAll(model interface{}) *gorm.DB
	GetById(id int, model interface{}) *gorm.DB
	DeleteByID(id int, model interface{}) *gorm.DB
	FindUserByEmail(email string, model *dbmodels.User) *gorm.DB
	GetUserByID(id int, model *dbmodels.User) *gorm.DB
	ComparePasswords(hashedPass string, bodyPass string) error
	GetAllGames(userID int, model interface{}) *gorm.DB
	GetGameById(userID int, id int, model interface{}) *gorm.DB
}
