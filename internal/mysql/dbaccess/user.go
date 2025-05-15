package dbaccess

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `gorm:"type:varchar(50);not null"`
	Password string `gorm:"type:varchar(50);not null"`
}

func (u *User) TableName() string {
	return "users"
}
