package common

import (
	"Hertz-Scaffold/biz/model"
	"github.com/duke-git/lancet/v2/netutil"
	"net/http"
	"time"
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

func ForwardFormUrl(tenant *model.PlatformTenant, url string, payload map[string]string) (int, map[string]interface{}) {
	url = tenant.Host + url
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
