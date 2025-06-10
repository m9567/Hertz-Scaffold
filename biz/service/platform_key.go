package service

import (
	"Hertz-Scaffold/biz/dal"
	"Hertz-Scaffold/biz/model"
	"Hertz-Scaffold/conf"
	"github.com/cloudwego/hertz/pkg/app"
	"sync"
)

type PlatformKeyService interface {
	FindOneOrDefault(c *app.RequestContext, currency string, code string) *model.PlatformKey
	GetPlatformKey(c *app.RequestContext, currency string, code string) *model.PlatformKey
}

type PlatformKeyServiceProxy struct {
	common *CommonService
}

var (
	platformKeyService     PlatformKeyService
	platformKeyServiceOnce sync.Once
	platformRequestMap     = make(map[string]*model.PlatformKey)
)

func GetPlatformKeyService() PlatformKeyService {
	platformKeyServiceOnce.Do(func() {
		platformKeyService = &PlatformKeyServiceProxy{
			common: &CommonService{},
		}
	})
	return platformKeyService
}

func (s *PlatformKeyServiceProxy) GetPlatformKey(c *app.RequestContext, currency string, platformCode string) *model.PlatformKey {
	key := platformCode + ":" + currency
	platformKey := platformRequestMap[key]
	if platformKey == nil {
		platformKey = GetPlatformKeyService().FindOneOrDefault(c, currency, platformCode)
		platformRequestMap[key] = platformKey
	}
	return platformKey
}

func (s *PlatformKeyServiceProxy) FindOneOrDefault(c *app.RequestContext, currency string, code string) *model.PlatformKey {
	platformKey, _ := dal.GetPlatformKeyDal().FindOne(c, currency, code)
	if platformKey != nil {
		return platformKey
	}
	defaultCurrency := conf.AppConf.GetGameInfo().DefaultCurrency
	platformKey, _ = dal.GetPlatformKeyDal().FindOne(c, defaultCurrency, code)
	return platformKey
}
