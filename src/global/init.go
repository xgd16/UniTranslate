package global

type RunModeType uint8

const (
	HttpMode RunModeType = iota // HTTP 运行方式
	CmdMode                     // 终端运行方式
)

// RunMode 运行方式
var RunMode = HttpMode
