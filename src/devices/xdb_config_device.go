package devices

import (
	"uniTranslate/src/global"
	"uniTranslate/src/types"

	"github.com/xgd16/gf-x-tool/xstorage"
)

type XDbConfigDevice struct {
	xdb *xstorage.XDB
}

const defaultKeyName = "translate"

// NewXDbConfigDevice 创建一个XDbConfigDevice
func NewXDbConfigDevice() *XDbConfigDevice {
	return &XDbConfigDevice{
		xdb: global.XDB,
	}
}

func (t *XDbConfigDevice) Init() (err error) {
	return
}

func (t *XDbConfigDevice) GetConfig(refresh bool) (mapData map[string]*types.TranslatePlatform, err error) {
	if refresh {
		t.xdb.Init()
	}
	err = t.xdb.GetGJson().Get(defaultKeyName).Scan(&mapData)
	return
}

func (t *XDbConfigDevice) GetTranslateInfo(serialNumber string) (platform *types.TranslatePlatform, ok bool, err error) {
	err = t.xdb.Get(defaultKeyName, serialNumber).Scan(&platform)
	if err != nil {
		return
	}
	ok = platform != nil
	return
}

func (t *XDbConfigDevice) SaveConfig(serialNumber string, data *types.TranslatePlatform) (err error) {
	err = t.xdb.Set(defaultKeyName, serialNumber, data)
	return
}
