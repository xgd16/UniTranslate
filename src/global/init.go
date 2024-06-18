package global

type RunModeType uint8

const (
	HttpMode RunModeType = iota
	CmdMode
)

var RunMode = HttpMode
