package httpclient

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	netUrl "net/url"
	"strings"

	"liveJob/pkg/common/config"
)

var (
	defaultHttpClient *http.Client
	Header            = map[string]string{
		"Content-type": "application/json",
	}
)

type BasicAuth struct {
	Username string
	Password string
}

// POST 使用样例
// headerMap := make(map[string]string)
// basicAuth := httpclient.BasicAuth{}
// resp, err := httpclient.POST(url, []byte(queryDSL), headerMap, basicAuth)
func POST(path string, data []byte, header map[string]string, basicAuth ...BasicAuth) ([]byte, error) {
	payload := strings.NewReader(string(data))
	req, _ := http.NewRequest("POST", path, payload)
	req.Header.Add("content-type", "application/json")

	if len(basicAuth) > 0 {
		if basicAuth[0].Username != "" && basicAuth[0].Password != "" {
			req.SetBasicAuth(basicAuth[0].Username, basicAuth[0].Password)
		}
	}

	for key, value := range header {
		req.Header.Add(key, value)
	}

	rsp, err := HttpClient.Do(req)
	if err != nil {
		return []byte(""), err
	}
	defer rsp.Body.Close()

	body, err := io.ReadAll(rsp.Body)
	if err != nil {
		return []byte(""), err
	}

	if rsp.StatusCode == 200 {
		return body, nil
	}

	return []byte(""), nil
}

func GET(path string, header map[string]string, basicAuth ...BasicAuth) ([]byte, error) {

	req, _ := http.NewRequest("GET", path, nil)

	for key, value := range header {
		req.Header.Add(key, value)
	}

	if len(basicAuth) > 0 {
		if basicAuth[0].Username != "" && basicAuth[0].Password != "" {
			req.SetBasicAuth(basicAuth[0].Username, basicAuth[0].Password)
		}
	}

	rsp, err := HttpClient.Do(req)
	if err != nil {
		return []byte(""), err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", rsp.StatusCode)
	}

	body, err := io.ReadAll(rsp.Body)
	if err != nil {
		return []byte(""), err
	}

	return body, nil
}

// POSTJson 使用样例
// res, err := httpclient.POSTJson(postUrl, []byte(JsonStr), header, httpclient.GetClient(false, time.Second*10)
func POSTJson(path string, data []byte, header map[string]string, cli *http.Client) ([]byte, error) {
	if cli == nil {
		cli = HttpClient
	}

	payload := bytes.NewReader(data)
	req, err := http.NewRequest("POST", path, payload)
	if err != nil {
		return nil, err
	}

	req.Header.Add("content-type", "application/json")

	for key, value := range header {
		req.Header.Add(key, value)
	}

	rsp, err := cli.Do(req)
	if err != nil {
		return []byte(""), err
	}
	defer rsp.Body.Close()

	body, err := io.ReadAll(rsp.Body)
	if err != nil {
		return []byte(""), err
	}

	return body, nil
}

// ProxyPostJson  使用样例
// res, err := httpclient.ProxyPostJson(url, jsonStr, map[string]string{"": ""})
func ProxyPostJson(path string, data []byte, header map[string]string) ([]byte, error) {
	proxy := getProxy()
	myClient := &http.Client{Transport: &http.Transport{Proxy: proxy}}

	payload := bytes.NewReader(data)
	req, err := http.NewRequest("POST", path, payload)
	if err != nil {
		return nil, err
	}

	req.Header.Add("content-type", "application/json")

	for key, value := range header {
		req.Header.Add(key, value)
	}

	rsp, err := myClient.Do(req)
	if err != nil {
		return []byte(""), err
	}
	defer rsp.Body.Close()

	body, err := io.ReadAll(rsp.Body)
	if err != nil {
		return []byte(""), err
	}

	return body, nil
}

// newClient for connection re-use
func getProxy() func(*http.Request) (*netUrl.URL, error) {
	proxy := http.ProxyFromEnvironment
	if len(config.Config.App.ProxyURL) != 0 {
		par, err := netUrl.Parse(config.Config.App.ProxyURL)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		proxy = http.ProxyURL(par)
	}
	return proxy
}

func ProxyGet(path string, header map[string]string) ([]byte, error) {
	proxy := getProxy()
	myClient := &http.Client{Transport: &http.Transport{Proxy: proxy}}

	req, _ := http.NewRequest("GET", path, nil)
	for key, value := range header {
		req.Header.Add(key, value)
	}
	rsp, err := myClient.Do(req)
	if err != nil {
		return []byte(""), err
	}
	defer rsp.Body.Close()

	body, err := io.ReadAll(rsp.Body)
	if err != nil {
		return []byte(""), err
	}

	return body, nil
}
