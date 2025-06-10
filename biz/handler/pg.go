package handler

import (
	"Hertz-Scaffold/biz/bo"
	"Hertz-Scaffold/biz/constant"
	"Hertz-Scaffold/biz/model"
	"Hertz-Scaffold/biz/service"
	"Hertz-Scaffold/biz/utils/common"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"strings"
)

func init() {
	common.Register(constant.CallbackAPIModule, constant.MethodPost, "/game/call/pg/VerifySession", VerifySession)
	common.Register(constant.CallbackAPIModule, constant.MethodPost, "/:currency/game/call/pg/VerifySession", VerifySession)
}

type PG struct {
	code string
}

func (x *PG) getUsername(_ *model.PlatformKey, params interface{}) (string, bool) {
	return "", false
}

func VerifySession(ctx context.Context, c *app.RequestContext) {
	logger := common.GetCtxLogger(c)
	var params = make(map[string]string)
	c.PostArgs().VisitAll(func(k, v []byte) {
		params[string(k)] = string(v)
	})
	request := bo.PgBase{}
	err := c.Bind(&request)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	//platformKey := service.GetPlatformKeyService().GetPlatformKey(c, request.Currency, constant.PG)

	//username, done := CallbackRegistry[constant.PG].GetUsername(platformKey, params)
	username, done := params["player_name"], false
	if done {
		return
	}

	tenant := service.GetPlatformTenantService().GetPlatformTenant(c, constant.PG, username)
	if tenant == nil {
		return
	}
	b, f, _ := strings.Cut(c.FullPath(), "/:currency")
	url := b + f

	body, _ := c.Body()

	statusCode, tempMap := common.ForwardJson(tenant, url, string(body))
	c.JSON(statusCode, tempMap)

}
