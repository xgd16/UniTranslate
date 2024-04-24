package types

// ConfigDeviceInterface 配置
type ConfigDeviceInterface interface {
	// Init 系统启动时初始化驱动使用
	Init() (err error)
	// GetConfig 获取配置信息
	GetConfig(refresh bool) (mapData map[string]*TranslatePlatform, err error)
	// GetTranslateInfo 获取翻译平台信息
	GetTranslateInfo(serialNumber string) (platform *TranslatePlatform, ok bool, err error)
	// SaveConfig 存储配置
	SaveConfig(serialNumber string, data *TranslatePlatform) (err error)
}
