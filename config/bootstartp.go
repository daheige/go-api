package config

import (
	"errors"

	"github.com/daheige/thinkgo/redisCache"
	"github.com/daheige/thinkgo/yamlConf"

	"github.com/gomodule/redigo/redis"
)

var AppEnv string
var conf *yamlConf.ConfigEngine

func InitConf(path string) {
	conf = yamlConf.NewConf()
	conf.LoadConf(path + "/app.yaml")
}

func InitRedis() {
	//初始化redis
	redisConf := &redisCache.RedisConf{}
	conf.GetStruct("RedisCommon", redisConf)

	// log.Println(redisConf)
	redisConf.SetRedisPool("default")
}

//从连接池中获取redis client
func GetRedisObj(name string) (redis.Conn, error) {
	conn := redisCache.GetRedisClient(name)
	if conn == nil {
		return nil, errors.New("get redis client error")
	}

	return conn, nil
}
