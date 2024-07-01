package sleep

import (
	"time"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

var qqGroup int64 = 775223884
var qqGroupTest int64 = 118338724

var sleepImg = "/Users/lsymy/go_pro/bot/pig-bot/resource/others/sleep.png"

func init() {
	go sleepTask()
}

func sleepTask() {
	for {
		// 计算下一次9:30的时间
		now := time.Now()

		next := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		if now.After(next) {
			// 如果现在的时间已经过了今天的9:30，那么计算明天的9:30
			next = next.Add(24 * time.Hour)
		}

		// 计算从现在开始到下一次9:30的时间差
		duration := next.Sub(now)

		// 创建一个定时器
		timer := time.NewTimer(duration)

		// 等待定时器触发
		<-timer.C

		// 执行你的任务
		zero.RangeBot(func(id int64, ctx *zero.Ctx) bool {
			ctx.SendGroupMessage(qqGroup, message.Image("file:///"+sleepImg))
			return false
		})

	}
}
