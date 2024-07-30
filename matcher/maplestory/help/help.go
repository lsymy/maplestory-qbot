package help

import (
	"MSBot/config"
	"MSBot/db"
	rule "MSBot/rules"
	"fmt"
	"strings"
	"time"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

type maintenance_info struct {
	id     int
	detail string
}

type content struct {
	id      int
	keyword string
	content string
}

func init() {
	engine := zero.New()

	engine.OnRegex(`小猪(.*)`).Handle(func(ctx *zero.Ctx) {
		rematched := ctx.State["regex_matched"].([]string)
		keyword := strings.ReplaceAll(rematched[1], " ", "")

		database := db.GetDB()

		fmt.Println(keyword)
		// 查询要修改的记录
		row := database.QueryRow("SELECT id, keyword, content FROM learn_content WHERE keyword = ?", keyword)
		var info content
		err := row.Scan(&info.id, &info.keyword, &info.content)
		fmt.Println(info.content)
		fmt.Println(info.keyword)

		if err != nil {
			// 不存在
			fmt.Println("不存在")
			return
		}

		ctx.SendChain(message.Image(config.LocalResourceHost + "ms/" + info.content))
	})

	engine.OnFullMatch("小猪wiki", rule.CheckRule).Handle(func(ctx *zero.Ctx) {
		ctx.SendChain(message.Text("https://maplestory.fandom.com/wiki"))
	})

	engine.OnFullMatch("小猪计算器", rule.CheckRule).Handle(func(ctx *zero.Ctx) {
		ctx.SendChain(message.Text("联盟计算器：https://xenogents.github.io/LegionSolver/" + "\n" + "核心计算器：https://soundmark.github.io/maple-nodes/" + "\n" + "上星计算器：https://brendonmay.github.io/starforceCalculator/" + "\n"))
	})

	engine.OnFullMatch("小猪航海", rule.CheckRule).Handle(func(ctx *zero.Ctx) {
		ctx.SendChain(message.Text("一图(15): 香水×2+肥皂×2" + "\n" + "三图(12): 优质皮+眼镜+肥皂×2 + 肥皂×4×6次"))
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

	engine.OnFullMatch("小猪幻影森林", rule.CheckRule).Handle(func(ctx *zero.Ctx) {
		ctx.SendChain(message.Image(config.LocalResourceHost+"ms/"+"幻影森林1.png"), message.Image(config.LocalResourceHost+"ms/"+"幻影森林2.png"))
	})
}
