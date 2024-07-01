package at

import (
	zero "github.com/wdvxdr1123/ZeroBot"
)

func init() {
	engine := zero.New()

	engine.OnMessage().Handle(func(ctx *zero.Ctx) {
		// if ctx.Event.IsToMe {
		// 	msg := strings.ToLower(ctx.MessageString())
		// }
	})

}
