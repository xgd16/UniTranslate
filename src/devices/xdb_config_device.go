package devices

import (
	"uniTranslate/src/global"
	"uniTranslate/src/types"

	"github.com/gogf/gf/v2/frame/g"
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
	err = t.initCountRecord()
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

func (t *XDbConfigDevice) initCountRecord() (err error) {

	m := g.Model("count_record")

	for k := range t.xdb.GetGJson().Get("translate").MapStrVar() {
		count, err1 := m.Clone().Where("serialNumber", k).Count()
		if err1 != nil {
			err = err1
			return
		}
		if count > 0 {
			continue
		}
		if _, err = m.Clone().InsertIgnore(g.Map{
			"serialNumber": k,
		}); err != nil {
			return
		}
	}
	return
}
