package models

import "gorm.io/gorm"

/*
Nota
Cuando actualices la interfaz Istorage no te olvides de correr el comando:

mockgen -source=models/storage.go -destination=./handlers/mocks/IStorageMocks.go

para mockearla y asi no romper los test viejos y poder testeaar nos nuevos cambios
*/

type Istorage interface {
	CreateUser(user *User) *gorm.DB
	FindUserByEmail(email string, model *User) *gorm.DB
}
