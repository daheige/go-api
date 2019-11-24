package config

import (
	"github.com/daheige/thinkgo/mysql"
	"github.com/daheige/thinkgo/yamlconf"
)

var (
	AppEnv   string
	AppDebug bool

	conf   *yamlconf.ConfigEngine
	dbConf = &mysql.DbConf{}
)
