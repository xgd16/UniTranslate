package handler

import (
	"github.com/gogf/gf/v2/container/gqueue"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"uniTranslate/src/global"
	"uniTranslate/src/types"
)

var SaveQueue = gqueue.New()

func SaveQueueHandler() {
	ctx := gctx.New()
	for {
		if v := SaveQueue.Pop(); v != nil {
			if err := global.StatisticalProcess.SaveCache(v.(*types.SaveData)); err != nil {
				g.Log().Error(ctx, "计数统计操作失败", v, err)
			}
		}
	}
}
