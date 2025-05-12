package mysql

import (
	"fmt"
	"github.com/Pacific73/gorm-cache/cache"
	gormCacheConfig "github.com/Pacific73/gorm-cache/config"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"hzer/configs"
	"hzer/internal/mysql/dbaccess"
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

	mysqlDB, err = gorm.Open(mysql.Open(url), &gorm.Config{})
	//mysqlDB.LogMode(false)
	if err != nil {
		panic(err)
	}
	sqlDB, err := mysqlDB.DB()
	if err != nil {
		panic(err)
	}
	//始终保持的tcp连接数，即使连接都关闭了
	sqlDB.SetMaxIdleConns(30)
	//最大tcp连接数
	sqlDB.SetMaxOpenConns(300)

	if config.Mysql.RedisCache && config.Redis.Enable {
		//添加缓存中间件
		redisClient := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port),
			Password: config.Redis.Password,
		})
		gormCache, _ := cache.NewGorm2Cache(&gormCacheConfig.CacheConfig{
			CacheLevel:           gormCacheConfig.CacheLevelAll,
			CacheStorage:         gormCacheConfig.CacheStorageRedis,
			RedisConfig:          cache.NewRedisConfigWithClient(redisClient),
			InvalidateWhenUpdate: true, // when you create/update/delete objects, invalidate cache
			CacheTTL:             5000, // 5000 ms
			CacheMaxItemCnt:      5,    // if length of objects retrieved one single time
			// exceeds this number, then don't cache
		})

		//缓存中间件附加到gorm
		err = mysqlDB.Use(gormCache)
		if err != nil {
			panic(err)
		}
	}
	dbaccess.SetGlobalDB(mysqlDB)
	InitTable()
}

func InitTable() {
	//TODO: 添加数据表后在此注册
	mysqlDB.AutoMigrate(&dbaccess.User{})
}
