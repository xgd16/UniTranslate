package handler

import (
	"github.com/gogf/gf/v2/container/gqueue"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"uniTranslate/src/devices"
	"uniTranslate/src/types"
)

var CountRecordQueue = gqueue.New()

func CountRecordQueueHandler() {
	ctx := gctx.New()
	for {
		if v := CountRecordQueue.Pop(); v != nil {
			if err := devices.RecordHandler.CountRecord(v.(*types.CountRecordData)); err != nil {
				g.Log().Error(ctx, "计数统计操作失败", v, err)
			}
		}
	}
}
