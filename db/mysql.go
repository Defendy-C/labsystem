package db

import (
	"fmt"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"labsystem/configs"
	"labsystem/util/logger"
)

type MySQLDB struct {
	*configs.MySQLConfig
	DB *gorm.DB
}

func NewMySQL() *MySQLDB {
	config := configs.NewMySQLConfig()
	logger.Log.Info("connecting mysql...", zap.Any("config", config))
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true", config.Name, config.Password, config.Host, config.Port, config.DBName)
	if db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{}); err != nil {
		logger.Log.Panic("mysql connecting failed", zap.String("dsn", dsn), zap.Error(err))
		return nil
	} else {
		logger.Log.Info("mysql connecting finished")
		return &MySQLDB{config, db}
	}
}

func (db *MySQLDB) Ping() error {
	err := db.DB.Error
	return err
}
