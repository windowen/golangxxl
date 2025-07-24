package model

type (
	RtcKickOutUserReq struct {
		UserId   int `json:"userId"`
		RoomId   int `json:"roomId"`
		Duration int `json:"duration"`
	}

	RtcKickOutUserResp struct {
		StatusCode int    `json:"statusCode"`
		Message    string `json:"message"`
		Status     string `json:"status"`
		Id         int    `json:"id"`
	}
)
