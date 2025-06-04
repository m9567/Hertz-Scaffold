package handler

import (
	"Hertz-Scaffold/biz/bo"
	"Hertz-Scaffold/biz/constant"
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
	"strings"
)

func init() {
	common.Register(constant.DefaultAPIModule, constant.MethodPost, ":currency/jdb/action", JdbAction)
}

func JdbAction(ctx context.Context, c *app.RequestContext) {
	logger := common.GetCtxLogger(c)
	//
	request := bo.JdbRequest{}
	err := c.Bind(&request)
	if err != nil {
		return
	}
	x := c.Query("x")

	logger.Info("JdbAction1 %v x %v", request, x)

	keyJson := service.GetPlatformKeyService().Find(c, request.Currency, constant.JDB)
	body, _ := jdbAesDecrypt(keyJson, x)
	logger.Info("JdbAction %v", body)

	m := make(map[string]interface{})
	err = json.Unmarshal([]byte(body), &m)
	if err != nil {
		return
	}
	uid := m["uid"].(string)
	// native
	if strings.HasPrefix(uid, "9wyl") {
		url := "https://testapi.abcvip.website/callback/game/call/jdb/action?x=GxhqU_KBGxPLYEWVYjk4RPrjxj0JI3Y7LFmm0Pik2Nh6Gs-F7lB5eF66CfOw8b06JOWe5f_k_3o-qxpXFAt5SPLXyX5uUSe_zmw8iKD3y9W84OjIAjJ61CfpFsImyF9rIifmIrnAwnN9la6--95dWVL0RhQ6_4GnadhCFI5uQnJfEUWkQYswItA3azn5LKNnvH3Ze-MpBj4EqXTIOgo1wQ"
		method := "POST"
		client := &http.Client{}
		req, _ := http.NewRequest(method, url, nil)
		res, _ := client.Do(req)

		defer res.Body.Close()
		body, _ := io.ReadAll(res.Body)

		c.SetStatusCode(http.StatusOK)
		logger.Info("JdbAction response %v", string(body))
		_, _ = c.Write(body)
	}
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
