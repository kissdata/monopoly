/*
use go-sdl2, which is interoperability between go and SDL(C lib)
*/
package internal

import (
	"container/list"
	"log"
	"math/rand"
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

func DrawBoard(wi *sdl.Window) *list.List {
	surface, _ := wi.GetSurface()
	surface.FillRect(nil, 0)

	boardList := list.New()

	rand.Seed(time.Now().Unix())
	for i := 0; i < 38; i++ {
		grid := Grid{sdl.Rect{}, i, "", true, "",
			rand.Intn(200)}

		if i == 0 || i == 19 {
			grid.CanBuy, grid.Price = false, 0 // 起始点和牢房不可购买
		}
		// 确认每个方格的位置
		switch {
		case i == 0:
			grid.Rect = sdl.Rect{X: 1, Y: 1, W: 80, H: 66}
		case i == 19:
			grid.Rect = sdl.Rect{X: 821, Y: 545, W: 80, H: 66}
		case i > 0 && i < 11:
			grid.Rect = sdl.Rect{X: int32(1 + 82*i), Y: int32(1), W: 80, H: 66}
		case i >= 11 && i < 18:
			t := i - 10
			grid.Rect = sdl.Rect{X: 821, Y: int32(1 + 68*t), W: 80, H: 66}
		case i > 19 && i < 30:
			t := i - 20
			grid.Rect = sdl.Rect{X: int32(1 + 82*t), Y: int32(545), W: 80, H: 66}
		case i >= 30 && i < 37:
			t := i - 29
			grid.Rect = sdl.Rect{X: 1, Y: int32(1 + 68*t), W: 80, H: 66}
		}

		boardList.PushBack(grid)
	}

	for e := boardList.Front(); e != nil; e = e.Next() {
		grid := e.Value.(Grid)
		p := &grid.Rect
		if e.Value.(Grid).Id == 0 {
			surface.FillRect(p, 0xffE9967A) // 起点方块
		} else if grid.Id == 19 {
			surface.FillRect(p, 0xffF8F8Ff) // JAIL
		} else {
			surface.FillRect(p, 0xffEEE8AA)
		}
		log.Println(e.Value)
	}

	wi.UpdateSurface()
	return boardList
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

// 随机产生一个数字, 显示对应的点数
func StopDice(render *sdl.Renderer, dice *Dice) int {
	number := dice.Tossing()
	textureDice, _ := img.LoadTexture(render, dice.GetLogoPath())

	var clips [6]sdl.Rect
	for i := 0; i < 6; i++ {
		// x: 40 -120- 50 -120- 50 -120- 40
		clips[i].X, clips[i].Y = int32(40+(i%3)*170), int32(40+(i/3)*170)
		clips[i].W, clips[i].H = 120, 120
	}

	solidArea := &sdl.Rect{X: 260, Y: 150, W: 120, H: 120} // 骰子的固定显示位

	render.Copy(textureDice, &clips[number-1], solidArea)

	return number
}

func AppMain() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		log.Panic(err)
		return
	}
	defer sdl.Quit()
	// 不写的话, 首次加载图像资源会有延迟
	if err := img.Init(img.INIT_JPG | img.INIT_PNG); err != nil {
		log.Panic(err)
	}
	defer img.Quit()

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

	render, _ := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
	defer render.Destroy()

	var dice Dice
	dice.SetLogoPath(DiceJPG)

	running := true
	exchange := 1
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				log.Println("Game Over")
				running = false
			case *sdl.MouseButtonEvent:
				if t.State == sdl.PRESSED {
					if exchange == 1 {
						number := StopDice(render, &dice)
						PlayerA.MovePlayer(number, window)
						window.UpdateSurface()
						exchange++
					} else {
						number := StopDice(render, &dice)
						Computer.MovePlayer(number, window)
						window.UpdateSurface()
						exchange--
					}
				}
			}
		}
		sdl.Delay(500)
	}
}
