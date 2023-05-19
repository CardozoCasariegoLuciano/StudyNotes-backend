package storage

import "gorm.io/gorm"

func (st *Storage) Create(aModel interface{}) *gorm.DB {
	return st.db.Save(aModel)
}
