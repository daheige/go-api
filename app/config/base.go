package config

import (
	"github.com/daheige/thinkgo/mysql"
	"github.com/daheige/thinkgo/yamlConf"
)

var AppEnv string
var AppDebug bool
var conf *yamlConf.ConfigEngine
var dbConf = &mysql.DbConf{}
