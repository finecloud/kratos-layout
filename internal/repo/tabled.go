package repo

import "time"

type Tabled interface {
	TableName() string
}

type Event struct {
	ID        int64     `gorm:"primarykey; autoIncrement"`
	Site      string    `gorm:"type:varchar(32)"`
	EventType string    `gorm:"type:varchar(64)"`
	Body      string    `gorm:"type:text"`
	FbBody    string    `gorm:"type:text"`
	Reply     string    `gorm:"type:text"`
	Time      time.Time `gorm:"type:datetime"`
}

func (Event) TableName() string {
	return "event"
}
