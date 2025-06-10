package service

import (
	"Hertz-Scaffold/biz/dal"
	"Hertz-Scaffold/biz/model"
	"github.com/cloudwego/hertz/pkg/app"
	"strings"
	"sync"
)

type PlatformTenantService interface {
	List(c *app.RequestContext) []*model.PlatformTenant
	GetPlatformTenant(c *app.RequestContext, platformCode string, username string) *model.PlatformTenant
}

type PlatformTenantServiceProxy struct {
	common *CommonService
}

var (
	platformTenantService     PlatformTenantService
	platformTenantServiceOnce sync.Once
	platformTenantMap         = make(map[string]*model.PlatformTenant)
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

func (s *PlatformTenantServiceProxy) GetPlatformTenant(c *app.RequestContext, platformCode string, username string) *model.PlatformTenant {
	var tenant *model.PlatformTenant
	tenant = getPlatformTenant(platformCode, username, tenant)
	if tenant == nil {
		tenantList := GetPlatformTenantService().List(c)
		for _, item := range tenantList {
			key := item.PlatformCode + ":" + item.Prefix
			platformTenantMap[key] = item
		}
	}
	tenant = getPlatformTenant(platformCode, username, tenant)
	return tenant
}

func getPlatformTenant(platformCode string, username string, tenant *model.PlatformTenant) *model.PlatformTenant {
	for _, v := range platformTenantMap {
		if v.PlatformCode == platformCode {
			if strings.HasPrefix(username, v.Prefix) {
				tenant = v
				break
			}
		}
	}
	return tenant
}
