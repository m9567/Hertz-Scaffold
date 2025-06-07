package service

import (
	"Hertz-Scaffold/biz/dal"
	"Hertz-Scaffold/biz/model"
	"github.com/cloudwego/hertz/pkg/app"
	"sync"
)

type PlatformTenantService interface {
	List(c *app.RequestContext) []*model.PlatformTenant
}

type PlatformTenantServiceProxy struct {
	common *CommonService
}

var (
	platformTenantService     PlatformTenantService
	platformTenantServiceOnce sync.Once
)

func GetPlatformTenantService() PlatformTenantService {
	platformTenantServiceOnce.Do(func() {
		platformTenantService = &PlatformTenantServiceProxy{
			common: &CommonService{},
		}
	})
	return platformTenantService
}

func (s *PlatformTenantServiceProxy) List(c *app.RequestContext) []*model.PlatformTenant {
	return dal.GetPlatformTenantDal().List(c)
}
