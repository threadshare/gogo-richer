package model

import (
	"time"
)

// User struct
type User struct {
	ID        int64      `gorm:"primary_key" json:"id"`
	Username  string     `gorm:"type:varchar(25);unique_index" json:"username"`
	Password  string     `gorm:"not null" json:"-"`
	Level     int64      `json:"level"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `sql:"index" json:"-"`
}

// GetUserByID ...
func GetUserByID(id interface{}) (user User, err error) {
	err = DB.First(&user, id).Error
	return
}

// GetUser ...
func GetUser(username interface{}) (user User, err error) {
	err = DB.Where("username = ?", username).First(&user).Error
	return
}

// ListUser ...
func (user User) ListUser(size, page int64, order string) (total int64, users []User, err error) {
	err = DB.Model(&User{}).Where("level < ? OR id = ?", user.Level, user.ID).Count(&total).Error
	if err != nil {
		return
	}
	if size == -1 {
		size = 1000
	}
	err = DB.Where("level < ? OR id = ?", user.Level, user.ID).Order(toUnderScoreCase(order)).Limit(size).Offset((page - 1) * size).Find(&users).Error
	return
}
