package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

type Config struct {
	HistoryPath string
	RawPath string
	UrlsPath string
}

func main() {
	config, err := loadConfig()
	if err != nil {
		fmt.Println("加载配置文件失败，错误：" + err.Error())
		return
	}

	systemHostsPath, err :=  getSystemHostsPath()
	if err != nil {
		fmt.Println("获取系统hosts文件路径失败，错误：" + err.Error())
		return
	}

	err = updateHistoryHosts(systemHostsPath, config.HistoryPath)
	if err != nil {
		fmt.Println("历史hosts文件（" + config.HistoryPath + "） 更新失败, 错误：" + err.Error())
		return
	}

	err = updateSystemHosts(systemHostsPath, config.UrlsPath, config.RawPath)
	if err != nil {
		fmt.Println("系统hosts文件（" + systemHostsPath + "） 更新 失败，错误:" + err.Error())
		return
	}

	fmt.Println("执行成功")
}

func loadConfig() (*Config, error) {
	var config = &Config{}
	data, err := os.ReadFile("./config.json")
	if err != nil {
		return nil, errors.New("配置文件读取失败，错误：" + err.Error())
	}
	err = json.Unmarshal(data, config)
	if err != nil {
		return nil, errors.New("配置文件反序列化失败，错误：" + err.Error())
	}

	return config, nil
}

func getSystemHostsPath() (string, error) {
	systemHostsPath := ""
	switch runtime.GOOS {
	case "windows":
		systemHostsPath = "C:/Windows/System32/drivers/etc/hosts"
	case "linux", "darwin","ios" :
		systemHostsPath = "/etc/hosts"
	case "android":
		systemHostsPath = "/system/etc/hosts"
	default:
		return "", errors.New("抱歉，您的系统（" + runtime.GOOS + "）暂时不支持")
	}
	return systemHostsPath, nil
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
	
	_, err = historyFile.WriteString("更改时间：" +
	time.Now().Local().String() +
	"\n" +
	string(lasthosts) + 
	"\n~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~\n")
	
	if err != nil {
		return errors.New("historyHosts文件写入失败, 错误：" + err.Error())
	}

	return nil
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

	var info []byte
	switch runtime.GOOS {
	case "linux":
		info, err = exec.Command("bash", "-c", "/etc/init.d/networking restart").Output()
		fmt.Println(string(info))
	case "windows":
		info, err = exec.Command("ipconfig", "/flushdns").Output()
		fmt.Println(string(info))
	case "darwin":
		info, err = exec.Command("dscacheutil", "-flushcache").Output()
		fmt.Println(string(info))
	}
	if err != nil {
		return errors.New("DNS刷新失败, 错误：" + err.Error())
	}

	return err
}
