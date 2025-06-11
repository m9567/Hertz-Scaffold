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
	"strings"
)

var (
	jdbApiRequest = "apiRequest.do"
)

func init() {
	//todo 改造前/callback/game/call/jdb/action
	common.Register(constant.CallbackAPIModule, constant.MethodPost, "/game/call/jdb/action", JdbAction)
	//todo /callback/usd/game/call/jdb/action
	common.Register(constant.CallbackAPIModule, constant.MethodPost, "/game/:currency/call/jdb/action", JdbAction)

	common.Register(constant.InnerApIModel, constant.MethodPost, "/game/:currency/jdb/:path", JdbRequest)

}

func JdbAction(ctx context.Context, c *app.RequestContext) {
	currency := c.Param("currency")
	x := c.Query("x")
	platformKey := service.GetPlatformKeyService().GetPlatformKey(c, currency, constant.JDB)
	if platformKey == nil {
		return
	}

	username, done := getJdbUsername(platformKey, x)
	if done {
		return
	}
	tenant := service.GetPlatformTenantService().GetPlatformTenant(c, constant.JDB, username)
	if tenant == nil {
		return
	}

	b, f, _ := strings.Cut(c.FullPath(), "/:currency")
	url := b + f + "?"
	c.QueryArgs().VisitAll(func(k, v []byte) {
		if strings.HasSuffix(url, "?") {
			url = url + string(k) + "=" + string(v)
		} else {
			url = url + "&" + string(k) + "=" + string(v)
		}
	})
	statusCode, tempMap := common.ForwardJson(tenant, url, "")
	c.JSON(statusCode, tempMap)

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

func getJdbUsername(platformKey *model.PlatformKey, param interface{}) (string, bool) {
	body, _ := jdbAesDecrypt(platformKey.KeyJson, param.(string))
	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(body), &m)
	if err != nil {
		return "", true
	}
	uid := m["uid"].(string)
	return uid, false
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

func JdbRequest(ctx context.Context, c *app.RequestContext) {
	currency := c.Param("currency")
	path := c.Param("path")
	var params = make(map[string]string)
	c.PostArgs().VisitAll(func(k, v []byte) {
		params[string(k)] = string(v)
	})
	platformKey := service.GetPlatformKeyService().GetPlatformKey(c, currency, constant.JDB)
	var urlMap = make(map[string]string)
	json.Unmarshal([]byte(platformKey.UrlJson), &urlMap)
	statusCode, tempMap := common.ForwardFormUrl("", urlMap[path], params)
	c.JSON(statusCode, tempMap)

}
