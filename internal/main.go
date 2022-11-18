/*
use go-sdl2, which is interoperability between go and SDL(C lib)
*/
package internal

import (
	"log"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

func DrawBoard(wi *sdl.Window) {
	surface, _ := wi.GetSurface()
	surface.FillRect(nil, 0)

	// 起点方块
	rect := sdl.Rect{X: 1, Y: 1, W: 80, H: 66}
	surface.FillRect(&rect, 0xffE9967A)
	// JAIL
	rect = sdl.Rect{X: 821, Y: 545, W: 80, H: 66}
	surface.FillRect(&rect, 0xffF8F8Ff)

	for i := 1; i < 11; i++ {
		rect = sdl.Rect{X: int32(1 + 82*i), Y: int32(1), W: 80, H: 66}
		surface.FillRect(&rect, 0xffEEE8AA)

	}
	for i := 1; i < 8; i++ {
		rect = sdl.Rect{X: 821, Y: int32(1 + 68*i), W: 80, H: 66}
		surface.FillRect(&rect, 0xffEEE8AA)
	}
	for i := 9; i > 0; i-- {
		rect = sdl.Rect{X: int32(1 + 82*i), Y: int32(545), W: 80, H: 66}
		surface.FillRect(&rect, 0xffEEE8AA)
	}
	for i := 8; i > 0; i-- {
		rect = sdl.Rect{X: 1, Y: int32(1 + 68*i), W: 80, H: 66}
		surface.FillRect(&rect, 0xffEEE8AA)
	}

	wi.UpdateSurface()
}

// 功能: 玩家放置于起点方格
func (me *Player) Prepare(wi *sdl.Window) error {
	playerImg, err := img.Load(me.GetLogoPath())
	if err != nil {
		log.Println("角色logo加载失败, error: ", err)
		return err
	}
	defer playerImg.Free()

	surface, _ := wi.GetSurface()
	playerImg.BlitScaled(nil, surface, &sdl.Rect{X: 21, Y: 21, W: 36, H: 36})
	wi.UpdateSurface()
	return nil
}

// 功能: 移动玩家
func (me Player) MovePlayer(stepNumber int) {

}

func AppMain() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		log.Panic(err)
		return
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow(AppTitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		902, // 一行11个方块, 每个宽80, 间距2
		612, // 一行9 个方块, 每个高66, 间隔2
		sdl.WINDOW_SHOWN)
	if err != nil {
		log.Println("Failed to create window, error: ", err)
		return
	}
	defer window.Destroy()

	DrawBoard(window)

	// 确认玩家角色能显示
	err1, err2 := PlayerA.Prepare(window), Computer.Prepare(window)
	if err1 != nil || err2 != nil {
		log.Println("assets目录找不到玩家logo, err: ", err1, err2)
		return
	}

	PlayerA.MovePlayer(1)

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				log.Println("Game Over")
				running = false
				break
			}
		}
	}
}
