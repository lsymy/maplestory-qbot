package helper

import (
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

func init() {
	engine := zero.New()

	engine.OnPrefix("/").Handle(func(ctx *zero.Ctx) {
		args := ctx.State["args"].(string)

		if args == "小猪ping" {
			ctx.SendChain(message.At(ctx.Event.UserID), message.Text("pong"))
		}
	})
}
