package learn

import (
	"MSBot/config"
	"MSBot/db"
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

func init() {
	engine := zero.New()

	engine.OnRegex(`(.*)小猪学习(.*)`).Handle(func(ctx *zero.Ctx) {

		rematched := ctx.State["regex_matched"].([]string)
		img := rematched[0]
		keyword := strings.ReplaceAll(rematched[1], " ", "")

		re := regexp.MustCompile(`\[[^\]]+\]`)
		cqstring := re.FindString(img)

		if keyword == "" || cqstring == "" {
			ctx.SendChain(message.Text("学习失败！请确认格式 [keyword 小猪学习 [Image]"))
			return
		}

		fmt.Println(keyword)
		fmt.Println(cqstring)
		err := saveImageForCQImage(cqstring, keyword)
		if err != nil {
			ctx.SendChain(message.Text("学习失败！请确认格式 [keyword 小猪学习 [Image]"))
			return
		}
		res := saveContentToDB(cqstring, keyword)
		if res == true {
			ctx.SendChain(message.Text("小猪学习成功。输入小猪" + keyword + "即可查看"))
			return
		} else {
			ctx.SendChain(message.Text("学习失败！"))
			return
		}
	})
}

type content struct {
	id      int
	keyword string
	content string
}

func saveContentToDB(path string, keyword string) bool {
	database := db.GetDB()

	// 查询要修改的记录
	row := database.QueryRow("SELECT id, keyword, content FROM learn_content WHERE keyword = ?", keyword)

	var info content
	err := row.Scan(&info.id, &info.keyword, &info.content)
	if err != nil {
		// 不存在
		insert := `
		INSERT INTO learn_content (keyword, content)
		VALUES (?, ?)`
		_, err = database.Exec(insert, keyword, keyword+".png")
		return err == nil // 如果插入操作成功，则返回 true
	}
	updateQuery := "UPDATE learn_content SET content = ? WHERE keyword = ?"
	_, err = database.Exec(updateQuery, keyword+".png", keyword)
	if err != nil {
		return false
	}

	return true
}

func saveImageForCQImage(cqstring string, filename string) (err error) {
	// 提取url
	start := strings.Index(cqstring, "url=") + 4
	end := strings.Index(cqstring, "]")
	url := cqstring[start:end]

	fileName := filename + ".png"
	downloadPath := config.LocalResourceAddress + "/ms/"

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
