package help

import (
	"MSBot/config"
	"MSBot/db"
	rule "MSBot/rules"
	"fmt"
	"time"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

var MsFullMatchKeywords = map[string]string{
	"小猪菜单":     "菜单.png",
	"小猪菜单2":    "菜单2.png",
	"小猪火花":     "火花.png",
	"小猪毕业副手":   "毕业副手.png",
	"小猪boss":   "boss.png",
	"小猪link":   "link.png",
	"小猪练级":     "lianji2.png",
	"小猪练级2":    "lianji.png",
	"小猪练级3":    "lianji3.png",
	"小猪练级4":    "lianji4.png",
	"小猪核心":     "核心.png",
	"小猪超技":     "超技.png",
	"小猪冒险家超技":  "冒险家超技.png",
	"小猪骑士团超技":  "骑士团超技.png",
	"小猪内潜":     "内潜.png",
	"小猪新内潜":    "新职业内潜.png",
	"小猪职业内潜":   "职业内潜.png",
	"小猪三级link": "3jlink.png",
	"小猪3级link": "3jlink.png",
	"小猪角色卡":    "角色卡.png",
	"小猪潜能":     "潜能.png",
	"小猪潜能2":    "潜能2.png",
	"小猪职业名称":   "职业名称.png",
	"小猪564":    "564.png",
	"小猪神子问答":   "神子问答.png",
	"小猪远征":     "远征.png",
	"小猪远征技能":   "远征技能.png",
	"小猪周日":     "周日.png",
	"小猪dmt":    "dmt.png",
	"小猪魔方":     "魔方.png",
	"小猪远征攻略":   "远征攻略.png",
	"小猪周常":     "周常.png",
	"小猪岛球":     "岛球.png",
	"小猪怪怪卡":    "怪怪.png",
	"小猪忍者城堡":   "忍者城堡.png",
	"小猪au":     "au.png",
	"小猪BOSS":   "BOSS2.png",
	"小猪斗燃":     "斗燃.png",
	"小猪刷图":     "刷图.png",
	"小猪战斗":     "战斗.png",
	"小猪100":    "100.png",
	"小猪托德":     "托德.png",
	"小猪装备":     "装备.png",
	"小猪活动":     "活动.png",
	"小猪buff":   "buff.png",
	"小猪创世":     "chuangshi.png",
}

type maintenance_info struct {
	id     int
	detail string
}

func init() {
	engine := zero.New()

	for keyword, image := range MsFullMatchKeywords {
		keyword := keyword
		image := image
		engine.OnFullMatch(keyword, rule.CheckRule).Handle(func(ctx *zero.Ctx) {
			sendImage(ctx, image)
		})
	}

	engine.OnFullMatch("小猪wiki", rule.CheckRule).Handle(func(ctx *zero.Ctx) {
		ctx.SendChain(message.Text("https://maplestory.fandom.com/wiki"))
	})

	engine.OnFullMatch("小猪联盟计算", rule.CheckRule).Handle(func(ctx *zero.Ctx) {
		ctx.SendChain(message.Text("https://xenogents.github.io/LegionSolver/"))
	})

	engine.OnFullMatch("小猪核心计算", rule.CheckRule).Handle(func(ctx *zero.Ctx) {
		ctx.SendChain(message.Text("https://soundmark.github.io/maple-nodes/"))
	})

	engine.OnFullMatch("小猪上星", rule.CheckRule).Handle(func(ctx *zero.Ctx) {
		ctx.SendChain(message.Text("https://brendonmay.github.io/starforceCalculator/"))
	})

	engine.OnFullMatch("小猪航海", rule.CheckRule).Handle(func(ctx *zero.Ctx) {
		ctx.SendChain(message.Text("一图: 香水×2+肥皂×2" + "\n" + "三图: 优质皮+眼镜+肥皂×2 + 肥皂×4×6次"))
	})

	duration := 30 * time.Minute
	engine.OnFullMatch("小猪ask提醒", rule.CheckRule).Handle(func(ctx *zero.Ctx) {
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("30分钟后将提醒你收菜，ask过程中进行切线、进商城等操作会导致奖励变为meso，建议在挂机和刷图时ask."))
		time.AfterFunc(duration, func() {
			// 在这里编写将消息发送给 QQ 的代码
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ask收菜!"))
		})
	})

	engine.OnFullMatch("小猪维护", rule.CheckRule).Handle(func(ctx *zero.Ctx) {
		database := db.GetDB()

		sqlStr := "select * from maintenance_info where id = 1;"
		row := database.QueryRow(sqlStr)

		var info maintenance_info
		row.Scan(&info.id, &info.detail)

		ctx.SendChain(message.Text(info.detail))
	})

	engine.OnFullMatch("小猪141", rule.CheckRule).Handle(func(ctx *zero.Ctx) {
		ctx.SendChain(message.Text("waste treatment plant 3 ---10j\nshaded dump site----------15j\nthe swamp of despair------ 25j\ndeep mire ------------------ 30j\n35j去f1野猪 或者\narmor pig land ------------- 37j\nmilitary camp 1------------- 42j\nsilent swamp --------------- 46j\nstairway to the sky 1 ------- 51j\nshaft 4  ---------------------- 64j\nsahel 2  ----------------------70j\nlab-area c-2 -----------------75j\nminar forest:west border ----84j\nsky nest 2   3  ---------------87j   5星\nzak\nforgotten passage   ---------102j  28星\nthe cave of trials 2  3 --------114j  55星"))
	})

	engine.OnFullMatch("小猪黎明脸", rule.CheckRule).Handle(func(ctx *zero.Ctx) {
		ctx.SendChain(message.Text("Greu vulture\nSkelosaurus\n[*]Skelosaurus\nAdvanced Knight A"))
	})
}

func sendImage(ctx *zero.Ctx, image string) {
	fmt.Println(config.LocalResourceHost + "ms/" + image)
	ctx.SendChain(message.Image(config.LocalResourceHost + "ms/" + image))
}
