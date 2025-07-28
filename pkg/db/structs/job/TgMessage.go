package job

type TgMessageUser struct {
	ID        int64  `json:"id"`
	IsBot     bool   `json:"is_bot,omitempty"` // 注意：omitempty 表示如果 false 可能不出现
	FirstName string `json:"first_name"`
	Username  string `json:"username,omitempty"`
	LangCode  string `json:"language_code,omitempty"`
}

type TgMessage struct {
	MessageID   int            `json:"message_id"`
	From        TgMessageUser  `json:"from"`
	ForwardFrom *TgMessageUser `json:"forward_from,omitempty"` // 可能是 nil
	Text        string         `json:"text"`
	ForwardDate int64          `json:"forward_date"` //时间

	// 其他字段...
}
