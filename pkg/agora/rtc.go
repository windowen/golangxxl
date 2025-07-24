package agora

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"liveJob/pkg/agora/model"
	"liveJob/pkg/common/config"
	"liveJob/pkg/tools/cast"
	"liveJob/pkg/tools/strhelper"
	"liveJob/pkg/zlogger"
)

var RtcClientInstance *rtcClient

type rtcClient struct {
	client *http.Client
	token  string
}

func NewRtcClient() {
	RtcClientInstance = &rtcClient{
		client: &http.Client{
			Timeout: time.Second * 5,
		},
		token: genToken(),
	}
}

// 生成令牌
func genToken() string {
	customerKey := config.Config.Agora.Rtc.CustomerKey
	customerSecret := config.Config.Agora.Rtc.CustomerSecret
	if customerKey == "" || customerSecret == "" {
		zlogger.Errorf("genToken | err: customer key or secret is missing")
		return ""
	}

	plainCredentials := fmt.Sprintf("%s:%s", customerKey, customerSecret)

	token := fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(plainCredentials)))

	return token
}

// post请求
func (rc *rtcClient) postDo(uri string, data any) (string, error) {
	body, err := json.Marshal(data)
	if err != nil {
		zlogger.Errorf("postDo marshal |data:%v| err: %v", cast.ToString(data), err)
		return "", err
	}

	// 创建请求
	req, err := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(body))
	if err != nil {
		zlogger.Errorf("postDo NewRequest err: %v", err)
		return "", err
	}

	// 设置请求头
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", rc.token)

	resp, err := rc.client.Do(req)
	if err != nil {
		zlogger.Errorf("postDo request post |data:%v| err: %v", body, err)
		return "", err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			zlogger.Errorf("postDo Body.Close | err:%v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		zlogger.Errorf("postDo resp.Status |status:%v| err: http response status error", resp.StatusCode)
		return "", fmt.Errorf("http response status error")
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		zlogger.Errorf("postDo ReadAll | err:%v", err)
		return "", err
	}

	return cast.ToString(bodyBytes), nil
}

// RtcKickOutUser 踢出直播间用户
func (rc *rtcClient) RtcKickOutUser(req model.RtcKickOutUserReq) error {
	var (
		agoraCfg = config.Config.Agora
		uri      = fmt.Sprintf("%s%s", agoraCfg.Rtc.Uri, RtcKickOutUser)
	)

	data := map[string]any{
		"appid":           agoraCfg.AppID,
		"cname":           cast.ToString(req.RoomId),
		"uid":             req.UserId,
		"time_in_seconds": req.Duration,
		"privileges":      []string{"join_channel"},
	}

	body, err := rc.postDo(uri, data)
	if err != nil {
		return err
	}

	var rtcKickResp model.RtcKickOutUserResp
	if err = strhelper.Json2Struct(body, &rtcKickResp); err != nil {
		zlogger.Errorln("RtcKickOutUser Json2Struct |body:%v| err:", cast.ToString(body), err)
		return err
	}

	if rtcKickResp.StatusCode != 0 && rtcKickResp.Message != "" {
		zlogger.Errorf("RtcKickOutUser |statusCode:%v| err: %v", rtcKickResp.StatusCode, rtcKickResp.Message)
		return fmt.Errorf(rtcKickResp.Message)
	}

	return nil
}
