package features

import (
	g "gorm.io/gorm"
)

type UserInfo struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"not null"`
}

func GetUserInfo(db *g.DB) (UserInfo, error) {
	var user UserInfo
	result := db.Last(&user)
	if result.Error != nil {
		return UserInfo{}, result.Error
	}
	return user, nil
}

func SaveUserInfo(db *g.DB, name string) error {
	user := UserInfo{Name: name}
	result := db.Create(&user)
	return result.Error
}

type Todo struct {
	ID     uint `gorm:"primaryKey"`
	Task   string
	Status string
}

func SaveTodo(db *g.DB, task string) error {
	todo := Todo{Task: task, Status: "pending"}
	result := db.Create(&todo)
	return result.Error
}

func ViewTodos(db *g.DB) ([]Todo, error) {
	var todos []Todo
	result := db.Find(&todos)
	if result.Error != nil {
		return nil, result.Error
	}
	return todos, nil
}

func UpdateTodo(db *g.DB, id uint) error {
	result := db.Model(&Todo{}).Where("id = ?", id).Update("status", "completed")
	return result.Error
}
