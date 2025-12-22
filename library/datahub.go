package library

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

type Metric struct {
	gorm.Model
	GID   string `json:"group_id"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

type SelectOptions func(*gorm.DB)

func init() {
	db, _ = gorm.Open(sqlite.Open("history.db"), &gorm.Config{})
	db.AutoMigrate(&Metric{})
}

func WithWhere(condition string, args ...any) SelectOptions {
	return func(db *gorm.DB) {
		db.Where(condition, args...)
	}
}

func WithOrder(order string) SelectOptions {
	return func(db *gorm.DB) {
		db.Order(order)
	}
}

func Select(opts ...SelectOptions) []Metric {
	var metrics []Metric
	for _, opt := range opts {
		opt(db)
	}
	db.Find(&metrics)
	return metrics
}

func Insert(metric Metric) error {
	return db.Create(&metric).Error
}
