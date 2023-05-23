package main

import (
	"awesomeProject/lib/File"
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
	var (
		useIcmp  bool
		outPut   string
		filePath string
		host     string
		silent   bool
		//ipList   []string
	)
	// 定义命令行参数
	flag.StringVar(&host, "host", "127.0.0.1", "输入IP/域名 地址，格式支持：192.168.21.1/24, 192.168.21.1-255, 192.168.21.1-192.168.21.255, www.example.com")
	flag.BoolVar(&useIcmp, "icmp", false, "是否确定进行icmp扫描")
	flag.StringVar(&filePath, "ipf", "", "指定IP段文件路径")
	flag.StringVar(&outPut, "output", "output.txt", "导出扫描结果到指定文件")
	flag.BoolVar(&silent, "silent", false, "是否输出结果至文件,默认不输出")
	flag.Parse()

	//port := flag.Int("port", 0, "端口号")
	//verbose := flag.Bool("verbose", false, "是否启用详细模式")
	start := time.Now()
	if useIcmp {
		if filePath == "" {
			if silent {
				File.CreateFile(outPut)
				exec.OutputSet(outPut)
				exec.IpIcmp(host, silent)
			}
			if silent == false {
				exec.IpIcmp(host, silent)
			}
		} else {
			if silent {
				File.CreateFile(outPut)
				ips, _ := File.ReadIPRangesFromFile(filePath)
				exec.OutputSet(outPut)
				for _, ip1 := range ips {
					exec.IpIcmp(ip1, silent)
				}
			}
		}

	}
	elapsed := time.Since(start)

	//输出执行时间。
	fmt.Println("扫描完成，共耗时: ", elapsed)
}
