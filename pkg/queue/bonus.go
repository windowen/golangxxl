package queue

import (
	"encoding/json"
	"fmt"
)

type CheckRegisterBonus struct {
	UserId int `json:"user_id"` // 用户id
}

func (mc CheckRegisterBonus) String() string {
	data, err := json.MarshalIndent(mc, "", "  ")
	if err != nil {
		return fmt.Sprintf("CheckRegisterBonus: error marshaling to JSON: %v", err)
	}
	return string(data)
}
