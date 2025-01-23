package global

import (
	"context"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/xgd16/gf-x-tool/xstorage"
)

// SystemConfigStruct 系统配置结构体
type ServerConfigStruct struct {
	Name                  string `json:"name"`
	Address               string `json:"address"`
	CacheMode             string `json:"cacheMode"`
	CachePlatform         bool   `json:"cachePlatform"`
	ConfigDeviceMode      string `json:"configDeviceMode"`
	RecordDeviceMode      string `json:"recordDeviceMode"`
	ConfigDeviceDb        string `json:"configDeviceDb"`
	CacheRefreshOnStartup bool   `json:"cacheRefreshOnStartup"`
	Key                   string `json:"key"`
	KeyMode               int    `json:"keyMode"`
	CacheWriteToStorage   bool   `json:"cacheWriteToStorage"`
	RequestRecordKeepDays int    `json:"requestRecordKeepDays"`
	ApiEditConfig         bool   `json:"apiEditConfig"`
}

type GraylogConfigStruct struct {
	Open bool   `json:"open"`
	Host string `json:"host"`
	Post string `json:"post"`
}

// SystemConfig 系统配置信息
var SystemConfig *gjson.Json

// ServerConfig 系统配置信息
var ServerConfig *ServerConfigStruct

// InitSystemConfig 初始化系统配置信息
func InitSystemConfig() (err error) {
	// 创建全局上下文对象
	cfg, err := g.Cfg().Data(Ctx)
	if err != nil {
		panic("初始化系统配置错误: " + err.Error())
	}
	SystemConfig = gjson.New(cfg, true)

	if err = SystemConfig.Get("server").Scan(&ServerConfig); err != nil {
		return
	}

	// 初始化 GfCache
	GfCache = gcache.New()
	return
}

// XDB 文件式存储
var XDB = xstorage.CreateXDB()

// GfCache 全局缓存对象
var GfCache *gcache.Cache

// Ctx 全局上下文对象
var Ctx context.Context
