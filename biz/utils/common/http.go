package common

import (
	"Hertz-Scaffold/biz/model"
	"github.com/duke-git/lancet/v2/netutil"
	"net/http"
	"time"
)

var (
	FormUrlencoded = "application/x-www-form-urlencoded"
)

func ForwardJson(tenant *model.PlatformTenant, url string, payload string) (int, map[string]interface{}) {
	url = tenant.Host + url
	httpClientCfg := netutil.HttpClientConfig{
		SSLEnabled:       true,
		HandshakeTimeout: 10 * time.Second,
	}
	request := &netutil.HttpRequest{
		RawURL:  url,
		Method:  "POST",
		Headers: make(http.Header),
	}
	if payload != "" {
		request.Body = []byte(payload)
	}
	request.Headers.Add("tenantCode", tenant.TenantCode)
	httpClient := netutil.NewHttpClientWithConfig(&httpClientCfg)

	res, _ := httpClient.SendRequest(request)

	var tempMap = make(map[string]interface{})
	err := httpClient.DecodeResponse(res, &tempMap)
	if err != nil {
		return res.StatusCode, nil
	}
	return res.StatusCode, tempMap
}

func ForwardFormUrl(tenantCode string, url string, payload map[string]string) (int, map[string]interface{}) {
	httpClientCfg := netutil.HttpClientConfig{
		SSLEnabled:       true,
		HandshakeTimeout: 10 * time.Second,
	}
	request := &netutil.HttpRequest{
		RawURL:   url,
		Method:   "POST",
		Headers:  make(http.Header),
		FormData: map[string][]string{},
	}
	if payload != nil {
		for k := range payload {
			request.FormData.Add(k, payload[k])
		}
	}
	if tenantCode != "" {
		request.Headers.Add("tenantCode", tenantCode)
	}
	request.Headers.Add("Content-Type", FormUrlencoded)
	httpClient := netutil.NewHttpClientWithConfig(&httpClientCfg)

	res, _ := httpClient.SendRequest(request)

	var tempMap = make(map[string]interface{})
	err := httpClient.DecodeResponse(res, &tempMap)
	if err != nil {
		return res.StatusCode, nil
	}
	return res.StatusCode, tempMap
}
