package data

import (
	"github.com/LikeRainDay/kratos-layout/internal/conf"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Data .
type Data struct {
	db  *gorm.DB
	log *log.Helper
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	newLog := log.NewHelper(logger)
	db, err := gorm.Open(mysql.Open(c.GetMysql().GetSource()), &gorm.Config{})
	if err != nil {
		newLog.Fatalf("failed opening connection to mysql: %v", err)
	}

	// 初始化sql连接池
	sqlDB, err := db.DB()
	if err != nil {
		newLog.Fatalf("failed get sql database: %v", err)
	}
	sqlDB.SetMaxIdleConns(int(c.GetMysql().GetMaxIdl()))
	sqlDB.SetMaxOpenConns(int(c.GetMysql().GetMaxOpen()))
	sqlDB.SetConnMaxLifetime(c.GetMysql().GetConnMaxLift().AsDuration())

	cleanup := func() {
		newLog.Info("closing the data resources")
		_ = sqlDB.Close()
	}
	return &Data{
		db:  db,
		log: newLog,
	}, cleanup, nil
}
