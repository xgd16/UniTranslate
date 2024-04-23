package cron

import (
	"context"
	"uniTranslate/src/global"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcron"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
)

func Service() {
	gcron.Add(gctx.New(), "0 */2 * * *", func(ctx context.Context) {
		defer func() {
			if err := recover(); err != nil {
				g.Log().Error(ctx, "清理请求记录失败", err)
			}
		}()
		g.Log().Infof(ctx, "每2小时执行一次 清理请求记录")
		if err := clearRequestRecord(ctx); err != nil {
			g.Log().Error(ctx, "清理请求记录失败", err)

		}
	})
}

func clearRequestRecord(ctx context.Context) (err error) {
	// 清理请求记录
	clearTime := gtime.Now().AddDate(0, 0, -global.RequestRecordKeepDays)
	delT, err := g.Model("request_record").Where("createTime <= ?", clearTime).Delete()
	delCount, _ := delT.RowsAffected()
	g.Log().Infof(ctx, "清理请求记录完成, 删除 %d 条记录 <= %s", delCount, clearTime.Format("Y-m-d H:i:s"))
	return
}
