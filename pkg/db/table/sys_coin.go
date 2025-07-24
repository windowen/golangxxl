package table

import "time"

type SysCoin struct {
	Id          int       `gorm:"column:id" json:"id"`
	Code        string    `gorm:"column:code" json:"code"`
	Symbol      string    `gorm:"column:symbol" json:"symbol"`
	Name        string    `gorm:"column:name" json:"name"`
	Status      int       `gorm:"column:status" json:"status"`
	IconUrl     string    `gorm:"column:icon_url" json:"icon_url"`
	CountryCode string    `gorm:"column:country_code" json:"country_code"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
}

func (SysCoin) TableName() string {
	return "sys_coin"
}
