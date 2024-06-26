package devices

import (
	"errors"
	"uniTranslate/src/types"
)

// RecordHandler 记录处理
var RecordHandler types.RecordInterface

// ConfigDevice 配置驱动
var ConfigDevice types.ConfigDeviceInterface

// GetConfigDevice 获取驱动配置
func GetConfigDevice() (device types.ConfigDeviceInterface, err error) {
	if ConfigDevice == nil {
		err = errors.New("配置获取驱动尚未初始化")
		return
	}
	device = ConfigDevice
	return
}

// MustGetConfigDevice 忽略错误获取驱动配置
func MustGetConfigDevice() (device types.ConfigDeviceInterface) {
	device, err := GetConfigDevice()
	if err != nil {
		panic(err)
	}
	return
}
