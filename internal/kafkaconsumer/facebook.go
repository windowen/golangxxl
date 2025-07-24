package kafkaconsumer

import (
	"encoding/json"
	"fmt"

	"go.uber.org/zap"

	"queueJob/pkg/tools/httpclient"
	"queueJob/pkg/tools/utils"
	"queueJob/pkg/zlogger"
)

type UserData struct {
	Emails          []string `json:"em,omitempty"`
	Phones          []string `json:"ph,omitempty"`
	ClientIP        string   `json:"client_ip_address,omitempty"`
	ClientUserAgent string   `json:"client_user_agent,omitempty"`
	Fbc             string   `json:"fbc,omitempty"`
	Fbp             string   `json:"fbp,omitempty"`
}

type Content struct {
	ID               string `json:"id"`
	Quantity         int    `json:"quantity"`
	DeliveryCategory string `json:"delivery_category"`
}

type CustomData struct {
	Currency    string     `json:"currency"`
	Value       float64    `json:"value"`
	Contents    []*Content `json:"contents"`
	ContentType string     `json:"content_type"`
	ContentIds  []string   `json:"content_ids"`
}

type FacebookEvent struct {
	EventName      string      `json:"event_name"`
	EventTime      int64       `json:"event_time"`
	ActionSource   string      `json:"action_source"`
	EventSourceURL string      `json:"event_source_url"`
	UserData       *UserData   `json:"user_data"`
	CustomData     *CustomData `json:"custom_data"`
}

func SendFacebookEvent(pixelID, accessToken string, event *FacebookEvent) error {
	url := fmt.Sprintf("https://graph.facebook.com/v18.0/%s/events?access_token=%s", pixelID, accessToken)

	// Hash sensitive data
	for i, email := range event.UserData.Emails {
		event.UserData.Emails[i] = utils.HashSHA256(email)
	}
	for i, phone := range event.UserData.Phones {
		event.UserData.Phones[i] = utils.HashSHA256(phone)
	}

	payload := map[string]interface{}{
		"data": []*FacebookEvent{event},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %v", err)
	}

	resp, err := httpclient.ProxyPostJson(url, body, nil)
	if err != nil {
		return fmt.Errorf("http post failed: %v", err)
	}

	zlogger.Infow("Facebook API response status",
		zap.String("pixelID", pixelID),
		zap.Any("event", event), zap.String("resp", string(resp)))
	return nil
}
