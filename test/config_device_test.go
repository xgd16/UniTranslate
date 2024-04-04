package test

import (
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
	"github.com/gogf/gf/v2/test/gtest"
	"testing"
	"uniTranslate/src/devices"
	"uniTranslate/src/types"
)

func TestXDbConfigDevice(t *testing.T) {
	DeviceConfigDeviceTest(devices.NewXDbConfigDevice())
}

func TestMySQLConfigDevice(t *testing.T) {
	DeviceConfigDeviceTest(devices.NewMySQLConfigDevice())
}

func DeviceConfigDeviceTest(device types.ConfigDeviceInterface) {
	gtest.AssertNil(device.Init())
	data := &types.TranslatePlatform{
		Platform: "test",
		Status:   1,
		Level:    1,
	}
	data.InitMd5()

	err := device.SaveConfig(
		data.Md5,
		data,
	)

	_, err = device.GetConfig()
	gtest.AssertNil(err)

	_, ok, err := device.GetTranslateInfo(data.Md5)
	gtest.AssertNil(err)

	gtest.Assert(ok, true)
}
