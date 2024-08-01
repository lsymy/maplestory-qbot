package character

import (
	rule "MSBot/rules"
	"regexp"

	"MSBot/db"
	"database/sql"

	"fmt"
	"net/http"

	"time"

	"encoding/json"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"

	"bytes"

	"github.com/wcharczuk/go-chart/v2" //exposes "chart"
	"github.com/wcharczuk/go-chart/v2/drawing"

	"github.com/FloatTech/floatbox/web"
)

type Character struct {
	CharacterData CharacterData `json:"CharacterData"`
}

type CharacterData struct {
	CharacterImageURL  string      `json:"CharacterImageURL"`
	Class              string      `json:"Class"`
	ClassRank          float64     `json:"ClassRank"`
	EXP                float64     `json:"EXP"`
	EXPPercent         float64     `json:"EXPPercent"`
	GlobalRanking      float64     `json:"GlobalRanking"`
	LegionLevel        float64     `json:"LegionLevel"`
	LegionCoinsPerDay  float64     `json:"LegionCoinsPerDay"`
	Level              float64     `json:"Level"`
	Name               string      `json:"Name"`
	Server             string      `json:"Server"`
	ServerClassRanking float64     `json:"ServerClassRanking"`
	ServerRank         float64     `json:"ServerRank"`
	GraphData          []GraphData `json:"GraphData"`
}

type GraphData struct {
	AvatarURL        string `json:"AvatarURL"`
	ClassID          int    `json:"ClassID"`
	ClassRankGroupID int    `json:"ClassRankGroupID"`
	CurrentEXP       int    `json:"CurrentEXP"`
	DateLabel        string `json:"DateLabel"`
	EXPDifference    int    `json:"EXPDifference"`
	EXPToNextLevel   int    `json:"EXPToNextLevel"`
	ImportTime       int    `json:"ImportTime"`
	Level            int    `json:"Level"`
	Name             string `json:"Name"`
	ServerID         int    `json:"ServerID"`
	ServerMergeID    int    `json:"ServerMergeID"`
	TotalOverallEXP  int    `json:"TotalOverallEXP"`
}

type character_qq struct {
	id             int
	character_name string
	qqid           int
	from_group_id  int
}

func init() {

	engine := zero.New()

	engine.OnPrefix("小猪查询", rule.CheckRule).Handle(func(ctx *zero.Ctx) {
		args := ctx.State["args"].(string)
		// info := characterSearch(args)
		info, e := characterSearch(args)

		if e == 0 {
			ctx.SendChain(message.Text("未查询到角色信息"))
			return
		}

		if e == 1 {
			ctx.SendChain(message.Text("请求失败"))
			return
		}

		if e == 2 {
			ctx.SendChain(message.Text("解析失败"))
			return
		}

		sendCharacterInfo(info, ctx)
	})

	engine.OnPrefix("小猪绑定", rule.CheckRule).Handle(func(ctx *zero.Ctx) {
		characterName := ctx.State["args"].(string)
		if characterName == "" {
			ctx.SendChain(message.Text("输id啊"))
			return
		}

		database := db.GetDB()

		fmt.Println(ctx.Event.GroupID)

		err := UpdateOrInsert(database, ctx.Event.UserID, characterName, ctx.Event.GroupID)
		if err != nil {
			ctx.SendChain(message.Text(err))
			return
		}

		ctx.SendChain(message.Text("绑定成功"))
	})

	engine.OnFullMatch("izhu", rule.CheckRule).Handle(func(ctx *zero.Ctx) {
		database := db.GetDB()
		fmt.Println(ctx.Event.UserID)
		sqlStr := "select * from character_qq where qqid = ?;"
		row := database.QueryRow(sqlStr, ctx.Event.UserID)
		var user character_qq
		row.Scan(&user.id, &user.character_name, &user.qqid, &user.from_group_id)

		fmt.Println(user.character_name)
		// 查询
		info, e := characterSearch(user.character_name)
		if e == 0 {
			ctx.SendChain(message.Text("未查询到角色信息"))
			return
		}

		if e == 1 {
			ctx.SendChain(message.Text("请求失败"))
			return
		}

		if e == 2 {
			ctx.SendChain(message.Text("解析 api json 失败, 确认请求站点是否正常"))
			return
		}

		sendCharacterInfo(info, ctx)
	})

	engine.OnPrefix("izhu", rule.CheckRule).Handle(func(ctx *zero.Ctx) {
		// characterName := ctx.State["args"].(string)
		str := ctx.Event.RawMessage
		re := regexp.MustCompile(`qq=(\d+)`)
		match := re.FindStringSubmatch(str)
		if len(match) > 1 {
			qqNumber := match[1]
			fmt.Println("QQ number:", qqNumber)

			database := db.GetDB()

			sqlStr := "select * from character_qq where qqid = ?;"
			row := database.QueryRow(sqlStr, qqNumber)
			var user character_qq
			row.Scan(&user.id, &user.character_name, &user.qqid)

			// 查询
			info, e := characterSearch(user.character_name)
			if e == 0 {
				ctx.SendChain(message.Text("未查询到角色信息"))
				return
			}

			if e == 1 {
				ctx.SendChain(message.Text("请求失败"))
				return
			}

			if e == 2 {
				ctx.SendChain(message.Text("解析 api json 失败, 确认请求站点是否正常"))
				return
			}

			sendCharacterInfo(info, ctx)
		} else {
			return
		}
	})
}

func characterSearch(characterName string) (CharacterData, int) {
	fmt.Println("角色请求开始")

	fmt.Println(characterName)
	fmt.Println("https://api.maplestory.gg/v2/public/character/gms/" + characterName)
	req, err := http.NewRequest("GET", "https://api.maplestory.gg/v2/public/character/gms/"+characterName, nil)
	if err != nil {
		fmt.Println("角色请求失败")
		return CharacterData{}, 0
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("请求失败:", err)
		return CharacterData{}, 1
	}
	defer resp.Body.Close()

	var character Character
	if err := json.NewDecoder(resp.Body).Decode(&character); err != nil {
		fmt.Println("解析 api JSON 失败:", err)
		return CharacterData{}, 2
	}

	return CharacterData{
		CharacterImageURL:  character.CharacterData.CharacterImageURL,
		Name:               character.CharacterData.Name,
		Class:              character.CharacterData.Class,
		Level:              character.CharacterData.Level,
		LegionLevel:        character.CharacterData.LegionLevel,
		GlobalRanking:      character.CharacterData.GlobalRanking,
		EXP:                character.CharacterData.EXP,
		EXPPercent:         character.CharacterData.EXPPercent,
		ServerClassRanking: character.CharacterData.ServerClassRanking,
		LegionCoinsPerDay:  character.CharacterData.LegionCoinsPerDay,
		Server:             character.CharacterData.Server,
		GraphData:          character.CharacterData.GraphData,
	}, 200
}

func sendCharacterInfo(info CharacterData, ctx *zero.Ctx) {
	// 创建一个新的数据序列
	xValues := make([]float64, len(info.GraphData))
	yValues := make([]float64, len(info.GraphData))

	var prevEntry GraphData // 定义变量来保存上一个 entry 的信息
	// 将日期转换成横坐标值
	for i, entry := range info.GraphData {
		t, _ := time.Parse("2006-01-02", entry.DateLabel)
		xValues[i] = float64(t.Unix()) // 将时间转换为Unix时间戳，再转换为float64类型

		if i != 0 {
			if prevEntry.EXPDifference < 0 {
				yValues[i] = (float64(entry.CurrentEXP) + float64(prevEntry.EXPToNextLevel)) / 1e9
			} else {
				yValues[i] = float64(prevEntry.EXPDifference) / 1e9 // 将经验值转换成以B为单位
			}
		} else {
			yValues[i] = 0
		}

		prevEntry = entry
	}

	fmt.Println(yValues)

	// 定义变量来记录大于 0 的值的数量
	count := 0
	// 定义变量来记录大于 0 的值的总和
	sumDiffExp := 0.0

	// 迭代切片
	for _, value := range yValues {
		// 判断元素是否大于 0
		if value > 0 {
			sumDiffExp += value
			count++
		}
	}

	currentExp := info.EXP / 1e9
	currentUpLevelExp := info.EXP / (info.EXPPercent / 100) / 1e9

	avgDaysExp := sumDiffExp / float64(count)
	upLevelNeedExp := currentUpLevelExp - currentExp
	// 当前需要升级经验 / 平均日经验 = 天数
	estimateUpLevelDay := upLevelNeedExp / avgDaysExp
	fmt.Println(estimateUpLevelDay)

	characterInfo := message.Text(
		"角色名: "+info.Name+"\n",
		"职业: "+info.Class+"("+fmt.Sprintf("%d", int(info.ServerClassRanking))+")"+"\n",
		"等级: "+fmt.Sprintf("%d", int(info.Level))+"\n",
		"服务器: "+info.Server+"\n",
		"联盟等级: "+fmt.Sprintf("%d", int(info.LegionLevel))+"\n",
		"联盟币(日): "+fmt.Sprintf("%d", int(info.LegionCoinsPerDay))+"\n",
		"当前经验: "+fmt.Sprintf("%.2f", info.EXP/1e9)+"b"+"\n",
		"当前经验百分比: "+fmt.Sprintf("%.2f%%", info.EXPPercent)+"\n",
		fmt.Sprintf("%d", count)+"日均经验预计升级天数: "+fmt.Sprintf("%.2f天", estimateUpLevelDay),
	)

	roleImage, roleErr := web.GetData(info.CharacterImageURL)

	if len(xValues) != 0 && len(yValues) != 0 && roleErr == nil {
		graph := chart.Chart{
			XAxis: chart.XAxis{
				Name: "日期",
				ValueFormatter: func(v interface{}) string {
					unixTimeStamp := v.(float64) // 假设这是一个Unix时间戳
					t := time.Unix(int64(unixTimeStamp), 0)
					formattedDate := t.Format("01-02") // 格式化为 mm-dd 格式
					return formattedDate
				},
			},
			YAxis: chart.YAxis{
				Name: "Daily EXP",
			},
			Series: []chart.Series{
				chart.ContinuousSeries{
					XValues: xValues[1:],
					YValues: yValues[1:],
					Style: chart.Style{
						FillColor: drawing.ColorFromHex("#bfdcf6").WithAlpha(50),
					},
				},
			},
		}

		// 将图表渲染为 PNG 图片
		buffer := bytes.NewBuffer([]byte{})
		err := graph.Render(chart.PNG, buffer)
		if err != nil {
			fmt.Println(xValues)
			fmt.Println(yValues)
			fmt.Println("渲染图表失败:", err)
			return
		}

		ctx.SendChain(message.ImageBytes(roleImage), characterInfo, message.ImageBytes(buffer.Bytes()))
	} else {
		ctx.SendChain(message.ImageBytes(roleImage), characterInfo)
	}

}

// UpdateOrInsert 更新或插入数据
func UpdateOrInsert(db *sql.DB, qqid int64, character_name string, groupid int64) error {
	// 首先检查数据库中是否已经存在该记录
	var existingValue string
	err := db.QueryRow("SELECT qqid FROM character_qq WHERE qqid = ?", qqid).Scan(&existingValue)

	switch {
	case err == sql.ErrNoRows:
		// 如果没有找到记录，则执行插入操作
		_, err := db.Exec("INSERT INTO character_qq (qqid, character_name, from_group_id) VALUES (?, ?, ?)", qqid, character_name, groupid)
		if err != nil {
			return err
		}
	case err != nil:
		// 如果查询过程中发生错误，则返回错误
		return err
	default:
		// 如果找到了记录，则执行更新操作
		_, err := db.Exec("UPDATE character_qq SET character_name = ?, from_group_id = ? WHERE qqid = ?", character_name, groupid, qqid)
		if err != nil {
			return err
		}
	}

	return nil
}
