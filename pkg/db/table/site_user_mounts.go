package table

import (
	"time"
)

type UserMounts struct {
	Id          int       `gorm:"column:id" json:"id"`                     // ID
	UserId      int       `gorm:"column:user_id" json:"user_id"`           // 用户ID
	MountsId    int       `gorm:"column:mounts_id" json:"mounts_id"`       // 坐骑id
	IsSelected  int       `gorm:"column:is_selected" json:"is_selected"`   // 是否使用
	Status      int       `gorm:"column:status" json:"status"`             // 状态 1-正常 2-禁用 3-删除
	ExpiredTime time.Time `gorm:"column:expired_time" json:"expired_time"` // 过期时间
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`     // 创建时间
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"`     // 更新时间
}

func (*UserMounts) TableName() string {
	return "site_user_mounts"
}

func (a *UserMounts) IsEmpty() bool {
	if a == nil {
		return true
	}
	return a.Id == 0
}
