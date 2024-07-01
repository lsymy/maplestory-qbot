package rules

import (
	"fmt"

	token "MSBot/token"

	zero "github.com/wdvxdr1123/ZeroBot"
)

func CheckRule(ctx *zero.Ctx) bool {
	// 检测令牌是否足够
	if !token.Tb.Take() {
		fmt.Println("token fail")
		return false
	}

	if zero.OnlyGroup(ctx) {
		return true
	}

	return false
}
