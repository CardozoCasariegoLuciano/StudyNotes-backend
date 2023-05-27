package storage

import "gorm.io/gorm"

func (st *Storage) Save(aModel interface{}) *gorm.DB {
	return st.db.Save(aModel)
}
