package main

import (
	"awesomeProject/lib/exec"
	"flag"
	"fmt"
	"time"
)

func main() {
	const logo = ` 
██╗   ██╗███████╗ ██████╗ █████╗ ███╗   ██╗
╚██╗ ██╔╝██╔════╝██╔════╝██╔══██╗████╗  ██║
 ╚████╔╝ ███████╗██║     ███████║██╔██╗ ██║
  ╚██╔╝  ╚════██║██║     ██╔══██║██║╚██╗██║
   ██║   ███████║╚██████╗██║  ██║██║ ╚████║
   ╚═╝   ╚══════╝ ╚═════╝╚═╝  ╚═╝╚═╝  ╚═══╝
					by cxy
--------------------------------------------
`

	fmt.Print(logo)
	// 定义命令行参数
	ip := flag.String("ip", "127.0.0.1", "输入IP 地址，格式支持：192.168.21.1/24, 192.168.21.1-255, 192.168.21.1-192.168.21.255")
	flag.Parse()

	//port := flag.Int("port", 0, "端口号")
	//verbose := flag.Bool("verbose", false, "是否启用详细模式")
	start := time.Now()
	exec.IpIcmp(ip)
	elapsed := time.Since(start)

	//输出执行时间。
	fmt.Println("扫描完成，共耗时: ", elapsed)
}
