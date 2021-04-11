package configs

import (
	"github.com/BurntSushi/toml"
	"labsystem/util"
)

type MySQLConfig struct {
	Host     string `toml:"host"`
	Port     uint   `toml:"port"`
	Name     string `toml:"user_name"`
	Password string `toml:"password"`
	DBName   string `toml:"db_name"`
}

type RedisConfig struct {
	Host     string `toml:"host"`
	Port     uint   `toml:"port"`
	Password string `toml:"password"`
	DB       uint   `toml:"db"`
}

type DBConfig struct {
	MySQL *MySQLConfig `toml:"mysql"`
	Redis *RedisConfig `toml:"redis"`
}

type LogConfig struct {
	Env    Environment `toml:"environment"`
	Level  string      `toml:"level"`
	Output string      `toml:"output"`
}

type HttpConfig struct {
	Host string `toml:"host"`
	Port int `toml:"port"`
}

var (
	dbConfig  DBConfig
	logConfig LogConfig
	httpConfig HttpConfig
)

func init() {
	// db config
	setDBConfig()
	// log config
	setLogConfig()
	// http config
	setHttpConfig()
}

func setDBConfig() {
	dbConfigText, err := util.ReadAll(CurProjectPath() + "/configs/db.toml")
	if err != nil {
		panic("don't read db.toml")
	}
	if _, err = toml.Decode(string(dbConfigText), &dbConfig); err != nil {
		panic("don't decode db.toml:" + err.Error())
	}
}

func setLogConfig() {
	logConfigText, err := util.ReadAll(CurProjectPath() + "/configs/log.toml")
	if err != nil {
		panic("don't read log.toml")
	}
	if _, err = toml.Decode(string(logConfigText), &logConfig); err != nil {
		panic("don't decode log.toml:" + err.Error())
	}
}

func setHttpConfig() {
	httpConfigText, err := util.ReadAll(CurProjectPath() + "/configs/http.toml")
	if err != nil {
		panic("don't read http.toml")
	}
	if _, err = toml.Decode(string(httpConfigText), &httpConfig); err != nil {
		panic("don't decode http.toml:" + err.Error())
	}
}

func NewMySQLConfig() *MySQLConfig {
	return dbConfig.MySQL
}

func NewRedisConfig() *RedisConfig {
	return dbConfig.Redis
}

func NewLogConfig() *LogConfig {
	return &logConfig
}

func NewHttpConfig() *HttpConfig {
	return &httpConfig
}