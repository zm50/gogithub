package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

var (
	historyPath = "./historyHosts.txt"
	rawPath = "./rawHosts.txt"
	urlsPath = "./urls.txt"
)

func main() {
	systemHostsPath := ""
	switch runtime.GOOS {
	case "windows":
		systemHostsPath = "C:/Windows/System32/drivers/etc/hosts"
	case "linux", "darwin","ios" :
		systemHostsPath = "/etc/hosts"
	case "android":
		systemHostsPath = "/system/etc/hosts"
	default:
		fmt.Println("抱歉，您的系统（" + runtime.GOOS + "）暂时不支持")
		return
	}

	err := updateHistoryHosts(systemHostsPath, historyPath)
	if err != nil {
		fmt.Println("历史hosts文件（" + historyPath + "） 更新失败, 错误：" + err.Error())
		return
	}

	err = updateSystemHosts(systemHostsPath, urlsPath, rawPath)
	if err != nil {
		fmt.Println("系统hosts文件（" + systemHostsPath + "） 更新 失败，错误:" + err.Error())
		return
	}

	fmt.Println("执行成功")
}


func updateSystemHosts(systemHostsPath, urlsPath, rawPath string) error {
	urlsBytes, err := os.ReadFile(urlsPath)
	if err != nil {
		return errors.New("urls文件读取失败, 错误：" + err.Error())
	}

	reader := bufio.NewReader(bytes.NewReader(urlsBytes))

	raw, err := os.ReadFile(rawPath)
	if err != nil {
		return errors.New("rawHosts文件读取失败, 错误：" + err.Error())
	}

	res := strings.Builder{}
	res.Write(raw)
	res.WriteString("\n####################################\n")
	c := colly.NewCollector()

	for {
		url, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return errors.New("缓冲读取失败, 错误：" + err.Error())
		}

		var line string
		c.OnHTML("div div div span", func(h *colly.HTMLElement) {
			if h.Attr("class") != "Whwtdhalf w15-0 lh45" {
				return
			}
			
			if s:=h.Attr("style"); s == "cursor:pointer;" {
				line = h.Text + " " + string(url) + "\n"
			}
		})


		err = c.Visit("https://ip.tool.chinaz.com/" + string(url))
		if err != nil {
			return errors.New("GitHub地址解析失败，错误：" + err.Error())
		}

		res.WriteString(line)
	}

	err = os.WriteFile(systemHostsPath, []byte(res.String()), 0666)
	if err != nil {
		return errors.New("系统hosts文件 " + systemHostsPath + " 修改失败， 错误：" + err.Error())
	}

	return err
}

func updateHistoryHosts(systemHostsPath, historyPath string) error {
	lasthosts, err := os.ReadFile(systemHostsPath)
	if err != nil {
		return errors.New("系统hosts文件 " + systemHostsPath + " 读取失败, 错误：" + err.Error())
	}
	historyFile, err := os.OpenFile(historyPath,os.O_WRONLY | os.O_APPEND, 0666)
	if err != nil {
		return errors.New("historyHosts文件打开失败, 错误：" + err.Error())
	}
	defer historyFile.Close()
	
	_, err = historyFile.WriteString("\n~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~\n更改时间：" +
	time.Now().Local().String() +
	"\n" +
	string(lasthosts))
	
	if err != nil {
		return errors.New("historyHosts文件写入失败, 错误：" + err.Error())
	}

	return nil
}