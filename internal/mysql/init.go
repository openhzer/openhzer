package mysql

import (
	"fmt"
	"github.com/8treenet/gcache"
	"github.com/8treenet/gcache/option"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"hzer/configs"
)

var (
	mysqlDB   *gorm.DB
	mysqlConf configs.Mysql
)

func InitGorm(config configs.Database) {
	mysqlConf = config.Mysql
	var err error
	//连接gorm
	url := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		mysqlConf.UserName,
		mysqlConf.Password,
		mysqlConf.Host,
		mysqlConf.Port,
		mysqlConf.DataName,
		mysqlConf.Charset)
	mysqlDB, err = gorm.Open("mysql", url)
	mysqlDB.LogMode(false)
	if err != nil {
		panic(err)
	}
	sqlDB := mysqlDB.DB()
	if err != nil {
		panic(err)
	}
	//始终保持的tcp连接数，即使连接都关闭了
	sqlDB.SetMaxIdleConns(30)
	//最大tcp连接数
	sqlDB.SetMaxOpenConns(300)

	//添加缓存中间件
	opt := option.DefaultOption{}
	opt.Expires = 300              //缓存时间, 默认120秒。范围30-43200
	opt.Level = option.LevelSearch //缓存级别，默认LevelSearch。LevelDisable:关闭缓存，LevelModel:模型缓存， LevelSearch:查询缓存
	opt.AsyncWrite = false         //异步缓存更新, 默认false。 insert update delete 成功后是否异步更新缓存。 ps: affected如果未0，不触发更新。
	opt.PenetrationSafe = false    //开启防穿透, 默认false。 ps:防击穿强制全局开启。

	//缓存中间件附加到gorm.DB
	sRedis := fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port)
	gcache.AttachDB(mysqlDB, &opt, &option.RedisOption{
		Addr:     sRedis,
		Password: config.Redis.Password,
	})

	InitTable()
}

func InitTable() {
	//TODO: 添加数据表后在此注册
	//mysqlDB.AutoMigrate(&HzerAdmin{})

}
