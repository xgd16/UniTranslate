package queue

import (
	"uniTranslate/src/service/queue/handler"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

var register = map[string]func(){
	"CountRecord":   handler.CountRecordQueueHandler,
	"RequestRecord": handler.RequestRecordQueueHandler,
	"Save":          handler.SaveQueueHandler,
}

func Service() {
	ctx := gctx.New()
	for k, fn := range register {
		go func(k string, fn func()) {
			defer func() {
				if err := recover(); err != nil {
					g.Log().Error(ctx, "队列启动失败", err)
				}
			}()
			fn()
		}(k, fn)
		g.Log().Infof(ctx, "[%s] 队列已启动", k)
	}
}
