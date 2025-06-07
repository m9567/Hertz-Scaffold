package dal

import (
	"Hertz-Scaffold/biz/model"
	"Hertz-Scaffold/biz/utils/common"
	"github.com/cloudwego/hertz/pkg/app"
	"sync"
)

type PlatformTenantDal struct {
	*CommonDAL
}

var (
	platformTenantDal     *PlatformTenantDal
	platformTenantDalOnce sync.Once
)

func GetPlatformTenantDal() *PlatformTenantDal {
	platformTenantDalOnce.Do(func() {
		platformTenantDal = &PlatformTenantDal{}
	})
	return platformTenantDal
}

func (ins PlatformTenantDal) List(c *app.RequestContext) []model.PlatformTenant {
	logger := common.GetCtxLogger(c)
	db, err := ins.GetTransaction(c)
	if err != nil {
		logger.Error(err.Error())
		return nil
	}
	var temps []model.PlatformTenant
	res := db.Table(model.PlatformTenant{}.TableName()).Find(&temps)
	if res.Error != nil {
		logger.Error(res.Error.Error())
		return nil
	}
	return temps
}
