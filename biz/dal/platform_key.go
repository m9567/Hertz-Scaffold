package dal

import (
	"Hertz-Scaffold/biz/model"
	"Hertz-Scaffold/biz/utils/common"
	"github.com/cloudwego/hertz/pkg/app"
	"sync"
)

type PlatformKeyDal struct {
	*CommonDAL
}

var (
	platformKeyDal     *PlatformKeyDal
	platformKeyDalOnce sync.Once
)

func GetPlatformKeyDal() *PlatformKeyDal {
	platformKeyDalOnce.Do(func() {
		platformKeyDal = &PlatformKeyDal{}
	})
	return platformKeyDal
}

func (ins *PlatformKeyDal) FindOne(c *app.RequestContext, currency string, code string) (dto *model.PlatformKey, err error) {
	logger := common.GetCtxLogger(c)
	db, err := ins.GetTransaction(c)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	temp := &model.PlatformKey{}
	res := db.Table(model.PlatformKey{}.TableName()).Where("currency = ? AND code = ?", currency, code).First(&temp)
	if res.Error != nil {
		logger.Error(res.Error.Error())
		return nil, res.Error
	}
	return temp, nil
}
