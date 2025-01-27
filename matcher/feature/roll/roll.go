package mofish

import (
	"fmt"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"

	"crypto/rand"
	"math/big"
)

func init() {
	engine := zero.New()

	engine.OnPrefix("/").Handle(func(ctx *zero.Ctx) {
		args := ctx.State["args"].(string)

		if args == "roll" {
			max := big.NewInt(100) // 生成0到100之间的随机数
			num, err := rand.Int(rand.Reader, max)
			if err != nil {
				fmt.Println("生成随机数时出错：", err)
				return
			}

			fmt.Println(ctx.Event.MessageID)
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(num))
		}
	})
}
