package config

import (
	"errors"
	"log"

	"github.com/daheige/thinkgo/gredigo"
	"github.com/daheige/thinkgo/yamlconf"
	"github.com/gomodule/redigo/redis"
)

// InitConf 初始化配置文件
func InitConf(path string) {
	conf = yamlconf.NewConf()
	err := conf.LoadConf(path + "/app.yaml")
	if err != nil {
		log.Fatalln("init config error: ", err)
	}

	AppEnv = conf.GetString("AppEnv", "production")
	switch AppEnv {
	case "local", "testing", "staging":
		AppDebug = true
	default:
		AppDebug = false
	}

	// 数据库配置
	conf.GetStruct("DbDefault", dbConf)
	dbConf.SetDbPool()              // 建立db连接池
	dbConf.SetEngineName("default") // 为每个db设置一个engine name
}

// InitRedis 初始化redis
func InitRedis() {
	// 初始化redis
	redisConf := &gredigo.RedisConf{}
	conf.GetStruct("RedisCommon", redisConf)

	// log.Println(redisConf)
	redisConf.SetRedisPool("default")
}

// GetRedisObj 从连接池中获取redis client
// 用完就需要调用redisObj.Close()释放连接，防止过多的连接导致redis连接过多
// 导致当前请求而陷入长久等待，从而redis崩溃
func GetRedisObj(name string) (redis.Conn, error) {
	conn := gredigo.GetRedisClient(name)
	if conn == nil || conn.Err() != nil {
		return nil, errors.New("get redis client error")
	}

	return conn, nil
}
