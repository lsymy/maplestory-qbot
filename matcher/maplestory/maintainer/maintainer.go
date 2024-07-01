package maintainer

import (
	"MSBot/config"
	"MSBot/db"
	rule "MSBot/rules"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

type maintenance_info struct {
	id     int
	detail string
}

func init() {
	engine := zero.New()

	engine.OnCommand("小猪修改周日", rule.CheckRule).Handle(func(ctx *zero.Ctx) {
		msg := ctx.Event.Message.CQCode()

		re := regexp.MustCompile(`\[[^\]]+\]`)
		cqstring := re.FindString(msg)

		err := saveImageForCQImage(cqstring, "周日")
		if err != nil {
			ctx.SendChain(message.Text("信息维护失败！请确认格式"))
			return
		}

		ctx.SendChain(message.Text("信息维护成功！输入小猪周日可查看最新更新信息"))
	})

	engine.OnCommand("小猪修改活动", rule.CheckRule).Handle(func(ctx *zero.Ctx) {
		msg := ctx.Event.Message.CQCode()

		re := regexp.MustCompile(`\[[^\]]+\]`)
		cqstring := re.FindString(msg)

		err := saveImageForCQImage(cqstring, "活动")
		if err != nil {
			ctx.SendChain(message.Text("信息维护失败！请确认格式"))
			return
		}

		ctx.SendChain(message.Text("信息维护成功！输入小猪活动可查看最新更新信息"))
	})

	engine.OnCommand("小猪修改维护", rule.CheckRule).Handle(func(ctx *zero.Ctx) {
		args := ctx.State["args"].(string)

		database := db.GetDB()

		// 查询要修改的记录
		row := database.QueryRow("SELECT id, detail FROM maintenance_info WHERE id = 1")

		var info maintenance_info
		err := row.Scan(&info.id, &info.detail)
		if err != nil {
			// 处理错误
			return
		}

		// 修改记录的 detail 字段内容
		info.detail = args

		// 执行 UPDATE 语句来更新记录
		updatedRow, err := database.Exec("UPDATE maintenance_info SET detail = ? WHERE id = ?", info.detail, info.id)
		if err != nil {
			// 处理错误
			return
		}

		// 检查是否成功更新
		rowsAffected, err := updatedRow.RowsAffected()
		if err != nil {
			// 处理错误
			return
		}

		fmt.Printf("成功更新 %d 条记录\n", rowsAffected)

		ctx.SendChain(message.Text("信息维护成功！输入小猪维护可查看最新更新信息"))

	})
}

func saveImageForCQImage(cqstring string, filename string) (err error) {
	// 提取url
	start := strings.Index(cqstring, "url=") + 4
	end := strings.Index(cqstring, "]")
	url := cqstring[start:end]

	fileName := filename + ".png"
	downloadPath := config.LocalResourceAddress + "ms/"

	// 发起HTTP GET请求
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error downloading file:", err)
		return err
	}
	defer resp.Body.Close()

	// 创建目标文件
	filePath := filepath.Join(downloadPath, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return err
	}
	defer file.Close()

	// 将HTTP响应体内容写入文件
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return err
	}

	return nil
}
