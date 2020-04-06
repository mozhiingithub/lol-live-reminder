package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var (
	i   = flag.Int("i", 15, "interval爬虫时间间隔，单位为分。默认为15分钟。") // 爬取网页的时间间隔
	url = flag.String("l", "", "url指定赛事网址。")                // 预约赛事网址
)

func main() {

	// 获取flag参数
	flag.Parse()

	// 未指定网址，直接终止程序
	if "" == *url {
		errorCheck(errors.New("未输入赛事网址"))
	}

	// 根据给定网址，找出指定赛事的前一场的网址

	// 匹配出指定赛事网页的标识码
	// 例：https://m.wanplus.com/schedule/58607.html
	// 标识码为：58607
	code := regexp.MustCompile("[\\d]+").FindString(*url)
	if "" == code { // 匹配失败
		errorCheck(errors.New("未能匹配当前赛事标识码，请检查网址"))
	}
	codeInt, e := strconv.Atoi(code)
	errorCheck(e)

	// 生成前一场的网址,前一场的标识码为指定赛事标识码减1
	exURL := fmt.Sprintf("https://m.wanplus.com/schedule/%d.html", codeInt-1)

	// 在规定的时间间隔内，爬取指定网页源码
	for {
		// 获取源码
		resp, e := http.Get(exURL)
		errorCheck(e)
		bs, e := ioutil.ReadAll(resp.Body)
		errorCheck(e)
		s := string(bs)

		// 判断“已结束”是否出现在源码中
		if strings.Contains(s, "已结束") { // 出现，表明上一场比赛已结束
			break
		} else { // 未出现，表明上一场比赛还未结束或未开始，则过一个时间间隔后再判断
			time.Sleep(time.Duration(*i) * time.Minute)
		}
	}

	// 上一场比赛已结束，发出提醒
	log.Println("预约赛事即将开赛。")

	// 判断程序所在平台，使用响应的方式弹出提醒文本
	sys := runtime.GOOS
	var cmd string
	switch sys {
	case "linux":
		{
			cmd = "xdg-open"
		}
	case "windows":
		{
			cmd = "notepad"
		}
	default:
		{
			return
		}
	}
	exec.Command(cmd, "remind.txt").Start()

}

// 检查错误
func errorCheck(e error) {
	if nil != e {
		panic(e)
	}
}
