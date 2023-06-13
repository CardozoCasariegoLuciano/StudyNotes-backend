package storage

import "gorm.io/gorm"

func (st *Storage) Save(aModel interface{}) *gorm.DB {
	return st.db.Save(aModel)
}

func (st *Storage) GetAll(model interface{}) *gorm.DB {
	return st.db.Find(model)
}

func (st *Storage) GetById(id int, model interface{}) *gorm.DB {
	return st.db.Where("id = ?", id).First(model)
}

func (st *Storage) DeleteByID(id int, model interface{}) *gorm.DB {
	return st.db.Unscoped().Delete(model, id)
}
