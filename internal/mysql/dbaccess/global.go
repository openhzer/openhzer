package dbaccess

import "gorm.io/gorm"

var globalDB *gorm.DB

// BaseModel 包含基础字段和DB实例
type BaseModel struct {
	DB *gorm.DB `gorm:"-" json:"-"`
}

// TableNamer 定义获取表名的接口
type TableNamer interface {
	TableName() string
}

// Model 是一个组合接口，包含Gorm和自定义接口
type Model interface {
	TableNamer
	SetDB(*gorm.DB)
	GetDB() *gorm.DB
}

func (t *BaseModel) GetDB() *gorm.DB {
	return t.DB
}

func (t *BaseModel) SetDB(db *gorm.DB) {
	t.DB = db
}

func SetGlobalDB(db *gorm.DB) {
	globalDB = db
}

func GetGlobalDB() *gorm.DB {
	return globalDB
}

// GetModelDB 获取模型对应的DB实例
// dbaccess.GetModelDB[dbaccess.User]()
func GetModelDB[T any, PT interface {
	*T
	Model
}]() PT {
	var model T
	ptrModel := PT(&model)
	ptrModel.SetDB(globalDB.Model(ptrModel).Table(ptrModel.TableName()))
	return ptrModel
}
