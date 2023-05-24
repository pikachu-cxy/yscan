package main

import (
	"awesomeProject/lib/File"
	"awesomeProject/lib/Format"
	"awesomeProject/lib/exec"
	"flag"
	"fmt"
	"net"
	"strings"
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
		//useIcmp  bool
		outPut   string
		filePath string
		host     string
		//silent   bool
		//url      string
		port string
	)
	// 定义命令行参数
	flag.StringVar(&host, "host", "", "输入IP/域名/url地址，格式支持：192.168.21.1/24, 192.168.21.1-255, 192.168.21.1-192.168.21.255, www.example.com,http://www.example.com")
	//flag.BoolVar(&useIcmp, "icmp", false, "是否确定进行icmp扫描")
	//flag.StringVar(&url, "u", "", "输入url地址，执行指纹识别，poc探测")
	flag.StringVar(&port, "port", "top100", "输入需要扫描的端口,支持如下参数：full（全端口扫描）,top100,top1000 ,HttPorts(常见http端口）,1-65535（自定义端口范围）")
	flag.StringVar(&filePath, "hf", "", "指定需扫描的IP/域名/url 文件路径")
	flag.StringVar(&outPut, "output", "output.txt", "导出扫描结果到指定文件")
	//flag.BoolVar(&silent, "silent", false, "是否输出结果至文件,默认不输出")
	flag.Parse()

	//port := flag.Int("port", 0, "端口号")
	//verbose := flag.Bool("verbose", false, "是否启用详细模式")
	start := time.Now()
	//如只在命令行输入资产，则认为扫描资产数量不大，仅在命令行输出扫描结果
	if filePath == "" && host != "" {
		choose(host, port)
	}
	if filePath != "" && host == "" {
		File.CreateFile(outPut)
		ips, _ := File.ReadIPRangesFromFile(filePath)
		exec.OutputSet(outPut)
		for _, host := range ips {
			choose(host, port)
		}
	}

	elapsed := time.Since(start)

	//输出执行时间。
	fmt.Println("扫描完成，共耗时: ", elapsed)
}

func choose(host string, port string) {
	hosts, format := Format.ChooseFormat(host)
	switch strings.ToLower(format) {
	case "ip":
		//先判断是内网环境/外网环境
		//内网可以arp udp tcp http
		host = hosts[0]
		ip := net.ParseIP(host)
		//内网地址
		if ip.IsPrivate() {
			//内网判断方法 arp,udp常用端口，http常用端口

		} else {
			//外网判断方法 udp常用端口，http常用端口

		}
		if exec.OnePing(host) {
			portsMap, _ := exec.ParsePorts(port)
			exec.ScanPort(portsMap, host)
		}
	case "ips":
		//为了保证扫描效率，当无法ping通目标ip，则认为不存活
		ipAlive := exec.IpIcmp(hosts)
		portsMap, _ := exec.ParsePorts(port)
		for _, host := range ipAlive {
			exec.ScanPort(portsMap, host)
		}

	case "domain":
		//爆破域名 深度扫描
		//js爬取 深度扫描
		host := hosts[0]
		//如果对方禁ping 通过DNS解析判断存活
		if exec.OnePing(host) || exec.DnsLookup(host) {
			//todo
		}

	case "url":
		//访问连通性--指纹识别--poc探测
	}
}
