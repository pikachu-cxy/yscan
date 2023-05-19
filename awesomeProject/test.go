package main

import (
	"flag"
	"fmt"
)

func main() {
	// 定义命令行参数
	ip := flag.String("ip", "", "IP 地址")
	port := flag.Int("port", 0, "端口号")
	verbose := flag.Bool("verbose", false, "是否启用详细模式")

	// 解析命令行参数
	flag.Parse()

	// 打印用户输入的参数值
	fmt.Println("IP 地址:", *ip)
	fmt.Println("端口号:", *port)
	fmt.Println("详细模式:", *verbose)

	// 其他处理逻辑...
}
