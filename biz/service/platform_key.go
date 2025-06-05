package service

import (
	"Hertz-Scaffold/biz/dal"
	"Hertz-Scaffold/conf"
	"github.com/cloudwego/hertz/pkg/app"
	"sync"
)

type PlatformKeyService interface {
	Find(c *app.RequestContext, currency string, code string) string
}

type PlatformKeyServiceProxy struct {
	common *CommonService
}

var (
	platformKeyService     PlatformKeyService
	platformKeyServiceOnce sync.Once
)

func GetPlatformKeyService() PlatformKeyService {
	platformKeyServiceOnce.Do(func() {
		platformKeyService = &PlatformKeyServiceProxy{
			common: &CommonService{},
		}
	})
	return platformKeyService
}

func (s *PlatformKeyServiceProxy) Find(c *app.RequestContext, currency string, code string) string {
	json, _ := dal.GetPlatformKeyDal().Find(c, currency, code)
	if json == "" {
		defaultCurrency := conf.AppConf.GetGameInfo().DefaultCurrency
		keyJson, _ := dal.GetPlatformKeyDal().Find(c, defaultCurrency, code)
		return keyJson
	} else {
		return json
	}
}
