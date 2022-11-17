package main

import (
	"log"
	"monopoly/internal"
)

// a entry point to run game monopoly
func main() {
	// 配置玩家信息
	if ok := internal.PlayerInit(); !ok {
		log.Println("玩家的配置不符合规定, 修改后继续!")
		return
	}

	// 进入游戏界面
	internal.AppMain()
}
