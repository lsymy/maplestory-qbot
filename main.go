package main

import (
	"fmt"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/driver"

	_ "MSBot/helper"
	_ "MSBot/matcher"

	"MSBot/config"
	"MSBot/db"

	"net/http"
)

func init() {
	res := db.InitDB()

	if res != nil {
		return
	}
}

func main() {

	imageFs := http.FileServer(http.Dir(config.LocalResourceAddress))
	http.Handle("/", imageFs)

	go func() {
		fmt.Println("go func http")
		http.ListenAndServe(":8089", nil)
	}()

	zero.RunAndBlock(&zero.Config{
		NickName:      []string{"Pig Bot"},
		CommandPrefix: "/",
		SuperUsers:    []int64{405252510},
		Driver: []zero.Driver{
			// 与 bot 的 ws 主动被动通信方式二选一
			// driver.NewWebSocketClient("ws://0.0.0.0:6700/", ""),
			driver.NewWebSocketServer(16, "ws://0.0.0.0:6701/", ""),
		},
	}, nil)
}
