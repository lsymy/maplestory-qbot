package token

import (
	"time"
)

type TokenBucket struct {
	capacity   int       // 令牌桶容量
	tokens     int       // 当前令牌数
	refillRate int       // 令牌恢复速率
	lastRefill time.Time // 上次令牌恢复时间
}

var Tb *TokenBucket

func init() {
	Tb = NewTokenBucket(6, 1)
}

// 创建一个新的令牌桶
func NewTokenBucket(capacity, refillRate int) *TokenBucket {
	return &TokenBucket{
		capacity:   capacity,
		tokens:     capacity,
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

// 检查是否有足够的令牌可用，并减少令牌数量
func (tb *TokenBucket) Take() bool {
	now := time.Now()
	elapsed := now.Sub(tb.lastRefill).Minutes()     // 计算自上次取令牌以来的经过时间
	tokensToAdd := int(elapsed) * tb.refillRate * 2 // 计算需要恢复的令牌数量

	if tokensToAdd > 0 {
		tb.lastRefill = now // 更新上次令牌恢复时间
		tb.tokens += tokensToAdd
		if tb.tokens > tb.capacity {
			tb.tokens = tb.capacity // 限制令牌数不超过容量
		}
	}

	if tb.tokens > 0 {
		tb.tokens--
		return true // 有足够的令牌
	}

	return false // 令牌不足
}
