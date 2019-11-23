package config

import (
	"github.com/daheige/thinkgo/mysql"
	"github.com/daheige/thinkgo/yamlconf"
)

var AppEnv string
var AppDebug bool
var conf *yamlconf.ConfigEngine
var dbConf = &mysql.DbConf{}
