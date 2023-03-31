package tests

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"hzer/internal/redis"
)

func redisJson(c *gin.Context) {
	rc := redis.GetCoon()
	defer rc.Close()
	cs := struct {
		Openid string
		Key    string
		Uid    string
	}{
		Openid: "aassd",
		Key:    "keysss",
		Uid:    "tuid",
	}
	b, _ := json.Marshal(&cs)
	rc.Do("set", "testjson", string(b))
}
