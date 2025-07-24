package site

import (
	"gorm.io/gorm"

	"queueJob/pkg/tools/tx"
)

func NewJob(db *gorm.DB) *Job {
	return &Job{
		DB: db,
	}
}

type Job struct {
	DB *gorm.DB
	tx tx.Tx
}
