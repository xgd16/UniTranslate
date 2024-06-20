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
	SaveConfig(serialNumber string, isUpdate bool, data *TranslatePlatform) (err error)
	// DelConfig 删除配置
	DelConfig(serialNumber string) (err error)
	// UpdateStatus 更新状态
	UpdateStatus(serialNumber string, status int) (err error)
}

type MySQLInitItem struct {
	TableName string   // 表名
	Table     string   // 创建表sql
	Index     []string // 创建索引sql
}
