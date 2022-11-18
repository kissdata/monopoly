package internal

import (
	"log"
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

// 功能: 模拟骰子旋转
//
//	@param dur[int]    给定持续的时间, 单位秒
func (dice *Dice) Run(rend *sdl.Renderer, dur int) {
	dice.SetLogoPath(Dice2JPG)
	diceImg, err := img.Load(dice.GetLogoPath())
	if err != nil {
		log.Println("没有骰子图像, 不显示")
		return
	}
	diceImg.SetColorKey(false, sdl.MapRGB(diceImg.Format, 255, 255, 255))

	texture, err := rend.CreateTextureFromSurface(diceImg)
	if err != nil {
		log.Println("Failed to create texture, err: ", err)
		return
	}
	defer texture.Destroy()

	var clips [6]sdl.Rect // 存储拆分的图
	for i := 0; i < 6; i++ {
		// x: 25 -125- 20 -125- 20 -125-
		clips[i].X, clips[i].Y = int32(25+(i%3)*145), int32(25+(i/3)*143)
		clips[i].W, clips[i].H = 125, 125
	}

	solidArea := &sdl.Rect{X: 260, Y: 150, W: 112, H: 112} // 骰子的固定显示位

	myTimer := time.NewTimer(time.Duration(dur) * time.Second)
	i := 0
	run := true
	for run {
		select {
		case <-myTimer.C:
			run = false
		default:
			rend.Clear()                             // 先清空
			rend.Copy(texture, &clips[i], solidArea) // 再复制
			rend.Present()                           // 最后显示

			sdl.Delay(150)
			i = (i + 1) % 6
		}
	}

}

// 欢迎页
func HomePage() {
	window, err := sdl.CreateWindow(AppTitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		902, // 一行11个方块, 每个宽80, 间距2
		612, // 一行9 个方块, 每个高66, 间隔2
		sdl.WINDOW_SHOWN)
	if err != nil {
		log.Println("Failed to create window, error: ", err)
		return
	}
	defer window.Destroy()

	rend, _ := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	defer rend.Destroy()

	var d Dice
	d.Run(rend, 2)
}
