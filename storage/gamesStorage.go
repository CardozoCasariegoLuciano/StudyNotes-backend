package storage

import "gorm.io/gorm"

func (st *Storage) GetAllGames(id int, model interface{}) *gorm.DB {
	return st.db.Where("user_id = ?", id).Find(model)
}

func (st *Storage) GetGameById(userID int, id int, model interface{}) *gorm.DB {
	return st.db.Where("id = ? AND user_id = ?", id, userID).First(model)
}
