/*
use go-sdl2, which is interoperability between go and SDL(C lib)
*/
package internal

import (
	"log"
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type Dice struct {
	Number   int      // 1 ~ 6
	Position sdl.Rect // 固定位置

	logoPath string
}

// 方格类
type Grid struct {
	Rect    sdl.Rect
	Id      int
	Name    string
	CanBuy  bool   // 可购买
	Belongs string // 所属玩家
	Price   int
}

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

	// 固定起始点
	if me.Name == "Robot" {
		me.Position.Rect = sdl.Rect{X: 42, Y: 26, W: 36, H: 36}
	} else {
		me.Position.Rect = sdl.Rect{X: 4, Y: 26, W: 36, H: 36}
	}

	surface, _ := wi.GetSurface()
	playerImg.BlitScaled(nil, surface, &me.Position.Rect)
	wi.UpdateSurface()
	return nil
}

// 功能: 移动玩家
func (me *Player) MovePlayer(stepNumber int, wi *sdl.Window) {
	playerImg, _ := img.Load(me.GetLogoPath())
	defer playerImg.Free()

	surface, _ := wi.GetSurface()
	// 原位置以填充原色的方式恢复
	if me.Position.Id == 0 {
		surface.FillRect(&me.Position.Rect, 0xffE9967A)
	} else {
		surface.FillRect(&me.Position.Rect, 0xffEEE8AA)
	}

	i := 0
	switch {
	case me.Position.Rect.Y < 50: // 原先在顶行
		for i < stepNumber {
			if me.Position.Rect.X += int32(82); me.Position.Rect.X > 800 {
				break
			}
			i++
		}
		me.Position.Rect.Y += int32(68 * (stepNumber - i))

	case me.Position.Rect.X > 821: // 右侧
		for i < stepNumber {
			if me.Position.Rect.Y += int32(68); me.Position.Rect.Y > 545 {
				break
			}
			i++
		}
		me.Position.Rect.X -= int32(82 * (stepNumber - i))
	case me.Position.Rect.Y > 545: // 底行
		for i < stepNumber {
			if me.Position.Rect.X -= int32(82); me.Position.Rect.X < 81 {
				break
			}
			i++
		}
		me.Position.Rect.Y -= int32(68 * (stepNumber - i))

	case me.Position.Rect.X < 80: // 左侧
		for i < stepNumber {
			if me.Position.Rect.Y -= int32(68); me.Position.Rect.Y < 61 {
				break
			}
			i++
		}
		me.Position.Rect.X += int32(82 * (stepNumber - i))
	}

	playerImg.BlitScaled(nil, surface, &me.Position.Rect)
	wi.UpdateSurface()
}

func (dice *Dice) Prepare(rend *sdl.Renderer) {
	dice.SetLogoPath(DiceJPG)
	diceImg, err := img.Load(dice.GetLogoPath())
	if err != nil {
		log.Println("没有骰子图像, 不显示")
		return
	}
	diceImg.SetColorKey(true, sdl.MapRGB(diceImg.Format, 255, 255, 255))

	texture, err := rend.CreateTextureFromSurface(diceImg)
	if err != nil {
		log.Println("Failed to create texture, err: ", err)
		return
	}
	defer texture.Destroy()

	// 原图拆分
	var clips [6]sdl.Rect
	for i := 0; i < 6; i++ {
		// x: 42 -112- 56 -112- 56 -112- 42
		clips[i].X, clips[i].Y = int32(42+(i%3)*168), int32(40+(i/3)*168)
		clips[i].W, clips[i].H = 112, 112
	}

	solidArea := &sdl.Rect{X: 260, Y: 150, W: 112, H: 112} // 骰子的固定显示位

	i := 0
	for i < 6 {
		rend.Clear()                             // 先清空
		rend.Copy(texture, &clips[i], solidArea) // 再复制
		rend.Present()                           // 最后显示

		time.Sleep(1 * time.Duration(time.Second))
		i++
	}

}

func AppMain() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		log.Panic(err)
		return
	}
	defer sdl.Quit()

	HomePage()
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

	var dice Dice
	for i := 0; i < 20; i++ {
		PlayerA.MovePlayer(dice.Tossing(), window)
		time.Sleep(1 * time.Duration(time.Second))
		Computer.MovePlayer(dice.Tossing(), window)
		time.Sleep(1 * time.Duration(time.Second))
	}

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
