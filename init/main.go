package main

import (
	"fmt"
	"os"
	"runtime"
)

func main() {
	systemHostsPath := ""
	switch runtime.GOOS {
	case "windows":
		systemHostsPath = "C:/Windows/System32/drivers/etc/hosts"
	case "linux", "darwin","ios" :
		systemHostsPath = "hosts"
	case "android":
		systemHostsPath = "/system/etc/hosts"
	default:
		fmt.Println("抱歉，您的系统（" + runtime.GOOS + "）暂时不支持")
		return
	}

	data, err := os.ReadFile(systemHostsPath)
	if err != nil {
		fmt.Println("系统hosts文件 " + systemHostsPath + " 读取失败, 错误：" + err.Error())
		return
	}

	err = os.WriteFile("../rawHosts.txt", data, 0666)
	if err != nil {
		fmt.Println("rawHosts文件读取失败, 错误：" + err.Error())
		return
	}

	err = os.Chmod("../rawHosts.txt", 0111)
	if err != nil {
		fmt.Println("rawHosts文件权限设置失败，错误：" + err.Error())
		return
	}

	fmt.Println("初始化成功")
}