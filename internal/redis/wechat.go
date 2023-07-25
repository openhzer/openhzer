package redis

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/silenceper/wechat/v2/miniprogram/auth"
	"hzer/pkg/util"
	"time"
)

// CreateWechatTempUID 用于创建不重复的微信临时uuid
func CreateWechatTempUID() (uuid string, err error) {
	redisconn := pool.Get()
	defer redisconn.Close()
	var ex bool
	for {
		uuid, _ = util.GetUUID()
		keys := fmt.Sprintf("Hzer:JWT:Wechat:Temp:%s", uuid)
		ex, err = redis.Bool(redisconn.Do(
			"SETNX", keys, time.Now().Unix(),
		))
		if err != nil {
			return "", err
		}
		if ex {
			redisconn.Do(
				"EXPIRE", keys, 300,
			)
			break
		}
	}
	return
}

// SaveWechatTempInfo 缓存临时微信信息
func SaveWechatTempInfo(uuid string, auth auth.ResCode2Session) error {
	redisconn := pool.Get()
	defer redisconn.Close()
	keys := fmt.Sprintf("Hzer:JWT:Wechat:Info:%s", uuid)
	bys, _ := json.Marshal(&auth)
	_, err := redisconn.Do("set", keys, string(bys))
	return err
}

func GetWechatTempInfo(uuid string) {

}
