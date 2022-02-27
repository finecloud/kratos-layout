package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"time"
)

type Greeter struct {
	ID        int64     `gorm:"primarykey; autoIncrement"`
	Site      string    `gorm:"type:varchar(32)"`
	EventType string    `gorm:"type:varchar(64)"`
	Hello     string    `gorm:"hello:varchar(32)"`
	Body      string    `gorm:"type:text"`
	FbBody    string    `gorm:"type:text"`
	Reply     string    `gorm:"type:text"`
	Time      time.Time `gorm:"type:datetime"`
}

type GreeterRepo interface {
	CreateGreeter(context.Context, *Greeter) error
	UpdateGreeter(context.Context, *Greeter) error
}

type GreeterUsecase struct {
	repo GreeterRepo
	log  *log.Helper
}

func NewGreeterUsecase(repo GreeterRepo, logger log.Logger) *GreeterUsecase {
	return &GreeterUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *GreeterUsecase) Create(ctx context.Context, g *Greeter) error {
	return uc.repo.CreateGreeter(ctx, g)
}

func (uc *GreeterUsecase) Update(ctx context.Context, g *Greeter) error {
	return uc.repo.UpdateGreeter(ctx, g)
}
