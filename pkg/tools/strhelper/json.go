package strhelper

import "encoding/json"

func Struct2Json(ss interface{}) (string, error) {
	res, err := json.Marshal(ss)
	if nil != err {
		return "", err
	}
	return string(res), nil
}

func Json2Struct(body string, ss interface{}) error {
	err := json.Unmarshal([]byte(body), ss)
	if nil != err {
		return err
	}
	return nil
}
