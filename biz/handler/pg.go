package handler

import (
	"Hertz-Scaffold/biz/constant"
	"Hertz-Scaffold/biz/model"
	"Hertz-Scaffold/biz/service"
	"Hertz-Scaffold/biz/utils/common"
	"context"
	"encoding/json"
	"github.com/cloudwego/hertz/pkg/app"
	"strings"
)

func init() {
	common.Register(constant.CallbackAPIModule, constant.MethodPost, "/game/call/pg/VerifySession", pgCallbackCommon)
	common.Register(constant.CallbackAPIModule, constant.MethodPost, "/:currency/game/call/pg/VerifySession", pgCallbackCommon)

	common.Register(constant.CallbackAPIModule, constant.MethodPost, "/game/call/pg/Cash/Get", pgCallbackCommon)
	common.Register(constant.CallbackAPIModule, constant.MethodPost, "/:currency/game/call/pg/Cash/Get", pgCallbackCommon)

	common.Register(constant.CallbackAPIModule, constant.MethodPost, "/game/call/pg/Cash/TransferInOut", pgCallbackCommon)
	common.Register(constant.CallbackAPIModule, constant.MethodPost, "/:currency/game/call/pg/Cash/TransferInOut", pgCallbackCommon)

	common.Register(constant.CallbackAPIModule, constant.MethodPost, "/game/call/pg/Cash/Adjustment", pgCallbackCommon)
	common.Register(constant.CallbackAPIModule, constant.MethodPost, "/:currency/game/call/pg/Cash/Adjustment", pgCallbackCommon)

	common.Register(constant.CallbackAPIModule, constant.MethodPost, "/game/call/pg/Cash/UpdateBetDetail", pgCashUpdateBetDetail)
	common.Register(constant.CallbackAPIModule, constant.MethodPost, "/:currency/game/call/pg/Cash/UpdateBetDetail", pgCashUpdateBetDetail)
}

func pgGetUsername(_ *model.PlatformKey, params map[string]string) (string, bool) {
	playerName := params["player_name"]
	if playerName == "" {
		playerName = common.GetUserNameByToken(params["operator_player_session"])
	}
	if playerName == "" {
		return "", true
	} else {
		return playerName, false
	}
}

func pgCallbackCommon(ctx context.Context, c *app.RequestContext) {
	logger := common.GetCtxLogger(c)
	currency := c.Query("currency")
	var params = make(map[string]string)
	c.PostArgs().VisitAll(func(k, v []byte) {
		params[string(k)] = string(v)
	})

	platformKey := service.GetPlatformKeyService().GetPlatformKey(c, currency, constant.PG)
	logger.Info(platformKey.KeyJson)

	username, done := pgGetUsername(platformKey, params)
	if done {
		return
	}

	tenant := service.GetPlatformTenantService().GetPlatformTenant(c, constant.PG, username)
	if tenant == nil {
		return
	}
	b, f, _ := strings.Cut(c.FullPath(), "/:currency")
	url := b + f

	statusCode, tempMap := common.ForwardFormUrl(tenant, url, params)
	c.JSON(statusCode, tempMap)
}

func pgCashUpdateBetDetail(ctx context.Context, c *app.RequestContext) {
	res := make(map[string]interface{})
	json.Unmarshal([]byte("{\"data\":{\"is_success\":true},\"error\":null}"), &res)
	c.JSON(200, res)
}
