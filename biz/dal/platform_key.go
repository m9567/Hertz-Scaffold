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

func (ins *PlatformKeyDal) Find(c *app.RequestContext, currency string, code string) (keyJson string, err error) {
	logger := common.GetCtxLogger(c)
	db, err := ins.GetTransaction(c)
	if err != nil {
		logger.Error(err.Error())
		return "", err
	}

	tmp := &model.PlatformKey{}
	res := db.Table(model.PlatformKey{}.TableName()).Where("currency = ? AND code = ?", currency, code).First(&tmp)
	if res.Error != nil {
		logger.Error(res.Error.Error())
		return "", res.Error
	}
	return tmp.KeyJson, nil
}
