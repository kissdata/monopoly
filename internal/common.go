package internal

import (
	"bufio"
	"io"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Player struct {
	Name       string
	gender     rune // m/f
	createDate string
	Money      int
	Position   Grid

	logoPath string // logo位置
}

func (me *Player) SetLogoPath(jpg string) {
	me.logoPath = filepath.Join(RootDir, jpg)
}

func (me *Player) GetLogoPath() string {
	return me.logoPath
}

func (dice *Dice) SetLogoPath(jpg string) {
	dice.logoPath = filepath.Join(RootDir, jpg)
}

func (dice *Dice) GetLogoPath() string {
	return dice.logoPath
}

func (dice *Dice) Tossing() int {
	rand.Seed(time.Now().UnixNano())
	dice.Number = 1 + rand.Intn(6) // [1, 7)
	return dice.Number
}

var (
	AppTitle   = "大富翁"
	ConfigPath = "./configs/player.txt" // 以单点开头
	PlayerLogo = "./internal/assets/tinytiger.jpg"
	PCroleLogo = "./internal/assets/computer.png"
	DiceJPG    = "./internal/assets/dice.jpg"
	Dice2JPG   = "./internal/assets/pickdice.jpg" // 骰子旋转模拟

	RootDir  string
	PlayerA  Player
	Computer Player
)

// 功能: 找到项目的cmd目录
//
//	@return [string] 找不到时返回空串
func findDircmd() string {
	exeDir, _ := filepath.Abs("") // 可执行程序所在目录
	f := exeDir
	for i := 0; filepath.Base(f) != "cmd"; i++ {
		if i > 3 {
			break
		}
		f = filepath.Dir(f) // 往父级翻目录
	}
	if filepath.Base(f) != "cmd" {
		log.Println("项目中的cmd目录不存在, 请检查!")
		return ""
	}
	return f
}

// 功能: 解析配置文件内容
//
//	@return [string] 返回唯一有效的那行玩家数据
func parseConfigFile(fileAbs string) (content string) {
	// 逐行读取, 去除注释
	file, err := os.Open(ConfigPath)
	if err != nil {
		return ""
	}
	defer file.Close()
	line := bufio.NewReader(file)
	for {
		var data []byte
		if data, _, err = line.ReadLine(); err == io.EOF {
			break
		}
		// 忽略空行和 #开头的注释
		if len(data) == 0 || data[0] == 35 {
			continue
		}
		content = string(data)
		log.Println("player info in file: ", content)
		break
	}
	return content
}

// 功能: 删除字符串中的多余空格，有多个空格时，仅保留一个空格
func DeleteExtraSpace(src string) string {
	s2 := make([]byte, len(src))
	copy(s2, src)

	regstr := "\\s{2,}"              // 两个及两个以上空格的正则表达式
	reg, _ := regexp.Compile(regstr) // 编译正则表达式
	spc_index := reg.FindStringIndex(string(s2))
	for len(spc_index) > 0 {
		s2 = append(s2[:spc_index[0]+1], s2[spc_index[1]:]...) // 删除多余空格
		spc_index = reg.FindStringIndex(string(s2))            // 继续搜
	}
	return string(s2)
}

// 功能: 玩家初始化, 玩家信息在txt文件里
func PlayerInit() bool {
	var f string
	if f = findDircmd(); f == "" {
		return false
	}
	RootDir = filepath.Dir(f)
	ConfigPath = filepath.Join(RootDir, ConfigPath)

	_, err := os.Stat(ConfigPath)
	if err != nil {
		log.Println("there is no config file! error: ", err)
		return false
	}

	var content string
	if content = parseConfigFile(ConfigPath); len(content) < 5 {
		return false
	}
	content = DeleteExtraSpace(content)
	infoArr := strings.Split(content, " ")
	PlayerA.createDate, PlayerA.Name = infoArr[0], infoArr[1]
	if sexy := []rune(strings.Trim(infoArr[2], " ")); len(sexy) != 1 {
		log.Println("玩家性别存储错误, 只能是m/f")
		PlayerA.gender = 'f'
	} else {
		PlayerA.gender = sexy[0]
	}
	PlayerA.Money, _ = strconv.Atoi(infoArr[3])
	log.Printf("player: %+v", PlayerA)
	PlayerA.SetLogoPath(PlayerLogo)

	// 电脑玩家
	Computer.createDate = time.Now().Format("2006-01-02")
	Computer.gender = 'm'
	Computer.Name = "Robot"
	Computer.Money = 5000
	Computer.SetLogoPath(PCroleLogo)

	return true
}
