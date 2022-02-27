package data

import (
	"context"
	"github.com/LikeRainDay/kratos-layout/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

// TABLES 初始化数据库
var TABLES = []interface{}{
	biz.Greeter{},
}

type Migrate struct {
	data   *Data
	logger log.Logger
}

func NewMigrate(data *Data, logger log.Logger) *Migrate {
	return &Migrate{
		data:   data,
		logger: logger,
	}
}

func (m *Migrate) Seed(ctx context.Context) error {
	db := GetDb(ctx, m.data.DbProvider)
	return migrateDb(db, m.logger)
}

func migrateDb(db *gorm.DB, logger log.Logger) error {
	newLog := log.NewHelper(logger)
	for _, table := range TABLES {
		if !db.Migrator().HasTable(table) {
			err := db.Migrator().CreateTable(table)
			if err != nil {
				newLog.Fatalf("failed to create table[%s], err: %v", table, err)
			}
		}
		err := db.Migrator().AutoMigrate(table)
		if err != nil {
			newLog.Fatalf("failed to migrate table[%s], err: %v", table, err)
		}
	}
	return nil
}
