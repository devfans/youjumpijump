package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"runtime/debug"
	"time"

	jump "github.com/faceair/youjumpijump"
)

var similar *jump.Similar

var r = jump.NewRequest()

type ScreenshotRes struct {
	Value     string `json:"value"`
	SessionID string `json:"sessionId"`
	Status    int    `json:"status"`
}

func screenshot(ip string) (*ScreenshotRes, image.Image) {
	_, body, err := r.Get(fmt.Sprintf("http://%s/screenshot", ip))
	if err != nil {
		panic(err)
	}

	res := new(ScreenshotRes)
	err = json.Unmarshal(body, res)
	if err != nil {
		panic(err)
	}

	pngValue, err := base64.StdEncoding.DecodeString(res.Value)
	if err != nil {
		panic(err)
	}

	src, err := png.Decode(bytes.NewReader(pngValue))
	if err != nil {
		panic(err)
	}
	return res, src
}

func main() {
	defer func() {
		jump.Debugger()
		if e := recover(); e != nil {
			log.Printf("%s: %s", e, debug.Stack())
			fmt.Print("程序已崩溃，按任意键退出")
			var c string
			fmt.Scanln(&c)
		}
	}()

	var ip string
	fmt.Print("请输入 WebDriverAgentRunner 监听的 IP 和端口 (例如 192.168.9.94:8100):")
	_, err := fmt.Scanln(&ip)
	if err != nil {
		log.Fatal(err)
	}

	var inputRatio float64
	fmt.Print("请输入跳跃系数(推荐值 3.856，可适当调整区间 3.600- 4.000): 精确到千分位，目前我用的3.856 能坚持到 1800-> : ")
	_, err = fmt.Scanln(&inputRatio)
	if err != nil {
		log.Fatal(err)
	}

	similar = jump.NewSimilar(inputRatio)

	for {
		jump.Debugger()

		res, src := screenshot(ip)

		f, _ := os.OpenFile("jump.png", os.O_WRONLY|os.O_CREATE, 0600)
		png.Encode(f, src)
		f.Close()

		start, end := jump.Find(src)
		 isjump := true
		//game over 检测不到游戏 命令行提示 可以手动调用系统的关闭重开方法，待更新
		if start == nil {
			log.Print("请重新开始游戏，检测不到游戏画面")
			isjump = false
			//break
		} else if end == nil {
			isjump = false
			log.Print("请重新开始游戏，检测不到游戏画面")
			//break
		}
		if isjump {
			nowDistance := jump.Distance(start, end)
			similarDistance, nowRatio := similar.Find(nowDistance)

			log.Printf("起点位置:%v 结束位置:%v 距离:%.2f 系数:%.2f 变换后的系数:%v 跳跃时间:%.2fms =xxx%v", start, end, nowDistance, similarDistance, nowRatio, nowDistance*nowRatio,nowDistance * nowRatio / 1000)
			// 跳跃时间
			var jumptime float64
			//fmt.Print("请输入跳跃时间(可适当调整):")
			//_, err = fmt.Scanln(&jumptime)
			//fmt.Println("跳跃时间====",jumptime)
			//主要修改了跳跃系数
			jumptime = nowDistance/123.54/inputRatio;


			fmt.Println("距离时间=====",nowDistance/123.54/inputRatio);
			_, _, err = r.PostJSON(fmt.Sprintf("http://%s/session/%s/wda/touchAndHold", ip, res.SessionID), map[string]interface{}{
				"x":        200,
				"y":        200,
				"duration": jumptime,

					///nowDistance * nowRatio / 1000,
			})
			if err != nil {
				panic(err)
			}

			//go func() {
			//	time.Sleep(time.Duration(nowDistance*nowRatio/1000+50) * time.Millisecond)
			//	_, src := screenshot(ip)
			//	finally, _ := jump.Find(src)
			//
			//	f, _ = os.OpenFile("jump.test.png", os.O_WRONLY|os.O_CREATE, 0600)
			//	png.Encode(f, src)
			//	f.Close()
			//
			//	if finally != nil {
			//		finallyDistance := jump.Distance(start, finally)
			//		finallyRatio := (nowDistance * nowRatio) / finallyDistance
			//
			//		if finallyRatio > nowRatio/2 && finallyRatio < nowRatio*2 {
			//			similar.Add(finallyDistance, finallyRatio)
			//		}
			//	}
			//}()

			time.Sleep(time.Millisecond * 1500)
		}

	}
}
