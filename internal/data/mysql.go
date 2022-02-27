package data

import (
	"github.com/LikeRainDay/kratos-layout/internal/conf"
	"github.com/LikeRainDay/kratos-layout/internal/repo"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// TABLES 初始化数据库
var TABLES = []repo.Tabled{
	repo.Event{},
}

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

	// 数据库初始化
	for _, table := range TABLES {
		if !db.Migrator().HasTable(table) {
			err = db.Migrator().CreateTable(table)
			if err != nil {
				newLog.Fatalf("failed to create table[%s], err: %v", table.TableName(), err)
			}
		}
		err = db.Migrator().AutoMigrate(table)
		if err != nil {
			newLog.Fatalf("failed to migrate table[%s], err: %v", table.TableName(), err)
		}
	}

	cleanup := func() {
		newLog.Info("closing the data resources")
		_ = sqlDB.Close()
	}
	return &Data{
		db:  db,
		log: newLog,
	}, cleanup, nil
}
