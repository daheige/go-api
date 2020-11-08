package config

import (
	"github.com/daheige/thinkgo/mysql"
	"github.com/daheige/thinkgo/yamlconf"
)

var (
	// AppEnv app_env
	AppEnv string

	// AppDebug app debug
	AppDebug bool

	conf   *yamlconf.ConfigEngine
	dbConf = &mysql.DbConf{}
)
