package dbaccess

import "github.com/jinzhu/gorm"

var globalDB *gorm.DB

func SetGlobalDB(db *gorm.DB) {
	globalDB = db
}

func GetGlobalDB() *gorm.DB {
	return globalDB
}
