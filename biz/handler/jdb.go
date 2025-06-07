package handler

import (
	"Hertz-Scaffold/biz/bo"
	"Hertz-Scaffold/biz/constant"
	"Hertz-Scaffold/biz/model"
	"Hertz-Scaffold/biz/service"
	"Hertz-Scaffold/biz/utils/common"
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"io"
	"net/http"
	"net/http/httputil"
	"strings"
	"sync"
	"time"
)

var (
	JdbProxyMap   = make(map[string]*httputil.ReverseProxy)
	JdbMutex      = sync.Mutex{}
	JdbTenantMap  = make(map[string]*model.PlatformTenant)
	JdbRequestMap = make(map[string]*model.PlatformKey)
	USD           = "usd"
)

func init() {
	common.Register(constant.DefaultAPIModule, constant.MethodPost, "jdb/:currency/action", JdbAction)
	proxy()
}

func JdbAction(ctx context.Context, c *app.RequestContext) {
	logger := common.GetCtxLogger(c)
	//
	request := bo.JdbCallbackRequest{}
	err := c.Bind(&request)
	if err != nil {
		return
	}
	x := c.Query("x")

	platformKey := getPlatformKey(c, request.Currency)
	if platformKey == nil {
		return
	}
	uid, done := getUsername(platformKey, x)
	if done {
		return
	}
	tenant := getPlatformTenant(uid, c)
	if tenant == nil {
		return
	}

	forward(tenant, x, c, logger)

	//body
	//if strings.HasPrefix(uid, "9wyl") {
	//	url := "https://testapi.abcvip.website/callback/game/call/jdb/action?x=GxhqU_KBGxPLYEWVYjk4RPrjxj0JI3Y7LFmm0Pik2Nh6Gs-F7lB5eF66CfOw8b06JOWe5f_k_3o-qxpXFAt5SPLXyX5uUSe_zmw8iKD3y9W84OjIAjJ61CfpFsImyF9rIifmIrnAwnN9la6--95dWVL0RhQ6_4GnadhCFI5uQnJfEUWkQYswItA3azn5LKNnvH3Ze-MpBj4EqXTIOgo1wQ"
	//	method := "POST"
	//	client := &http.Client{}
	//	payload := strings.NewReader(body)
	//	req, _ := http.NewRequest(method, url, payload)
	//	req.Header.Add("Content-Type", "application/json")
	//	res, _ := client.Do(req)
	//
	//	defer res.Body.Close()
	//	body, err := io.ReadAll(res.Body)
	//
	//	c.SetStatusCode(http.StatusOK)
	//
	//	logger.Info("JdbAction response %v", string(body))
	//	_, _ = c.Write(body)
	//
	//	if err != nil {
	//		fmt.Println(err)
	//		return
	//	}
	//	fmt.Println(string(body))
	//}

}

func forward(tenant *model.PlatformTenant, x string, c *app.RequestContext, logger common.Logger) {
	url := tenant.Host + "/callback/game/call/jdb/action?x=" + x
	logger.Info("JdbAction url %v", url)
	method := "POST"
	client := &http.Client{}
	req, _ := http.NewRequest(method, url, nil)
	req.Header.Add("tenantCode", tenant.TenantCode)
	res, _ := client.Do(req)

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(res.Body)
	body, _ := io.ReadAll(res.Body)
	c.SetStatusCode(http.StatusOK)
	logger.Info("JdbAction response %v", string(body))
	var tempMap = make(map[string]interface{})
	err := json.Unmarshal(body, &tempMap)
	if err != nil {
		return
	}
	c.JSON(res.StatusCode, tempMap)
}

func getPlatformKey(c *app.RequestContext, currency string) *model.PlatformKey {
	platformKey := JdbRequestMap[currency]
	if platformKey == nil {
		platformKey = service.GetPlatformKeyService().FindOneOrDefault(c, currency, constant.JDB)
		JdbRequestMap[currency] = platformKey
	}
	return platformKey
}

func getUsername(platformKey *model.PlatformKey, x string) (string, bool) {
	body, _ := jdbAesDecrypt(platformKey.KeyJson, x)
	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(body), &m)
	if err != nil {
		return "", true
	}
	uid := m["uid"].(string)
	return uid, false
}

func getPlatformTenant(uid string, c *app.RequestContext) *model.PlatformTenant {
	var tenant *model.PlatformTenant
	for prefix, v := range JdbTenantMap {
		if strings.HasPrefix(uid, prefix) {
			tenant = v
		}
	}
	if tenant == nil {
		tenantList := service.GetPlatformTenantService().List(c)
		for _, item := range tenantList {
			JdbTenantMap[item.Prefix] = item
		}
	}
	for prefix, v := range JdbTenantMap {
		if strings.HasPrefix(uid, prefix) {
			tenant = v
		}
	}
	return tenant
}

func jdbAesDecrypt(keyJson string, encryptedText string) (string, error) {
	var jdbKey bo.JdbKey
	err := json.Unmarshal([]byte(keyJson), &jdbKey)
	if err != nil {
		return "", err
	}
	// Decode base64-encoded encryptedText
	decodedData, err := base64.RawURLEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", fmt.Errorf("invalid encrypted data: %w", err)
	}

	// Decrypt
	key := []byte(jdbKey.KEY)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("invalid decryption key: %w", err)
	}

	plainText := make([]byte, len(decodedData))
	iv := []byte(jdbKey.IV)
	blockMode := cipher.NewCBCDecrypter(block, iv)
	blockMode.CryptBlocks(plainText, decodedData)

	// Remove zero padding
	plainText = bytes.TrimRight(plainText, string([]byte{0}))
	return string(plainText), nil
}

func proxy() {
	proxyOneCurrency(USD)
}

func proxyOneCurrency(currency string) {
	http.HandleFunc("/jdb/"+currency+"/apiRequest.do", func(w http.ResponseWriter, r *http.Request) {
		platformKey := getPlatformKey(nil, currency)
		var jdbRequestUrl bo.JdbRequestUrl
		_ = json.Unmarshal([]byte(platformKey.UrlJson), &jdbRequestUrl)
		apiRequest := jdbRequestUrl.ApiRequest
		reverseProxy := JdbProxyMap[apiRequest]
		if reverseProxy != nil {
			reverseProxy.ServeHTTP(w, r)
			return
		}
		for i := 0; i < 10; i++ {
			lock := JdbMutex.TryLock()
			if lock {
				defer JdbMutex.Unlock()
				reverseProxy = common.NewProxy(apiRequest)
				JdbProxyMap[apiRequest] = reverseProxy
				reverseProxy.ServeHTTP(w, r)
			}
			time.Sleep(time.Millisecond * 100)
		}
	})
}
