package main

import (
	"awesomeProject/lib/File"
	"awesomeProject/lib/Format"
	p "awesomeProject/lib/Plugins"
	"awesomeProject/lib/exec"
	"flag"
	"fmt"
	_ "github.com/praetorian-inc/fingerprintx/pkg/plugins"
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
				by PikaChu
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
		port      string
		dict      bool
		path      bool
		pocs      bool
		searchPoc string
		icmp      string
		Plugins   string
		ps        p.PluginService
	)
	// 定义命令行参数
	flag.StringVar(&host, "host", "", "输入待扫描的IP/域名/url地址，格式支持：192.168.21.1/24, 192.168.21.1-255, 192.168.21.1-192.168.21.255, www.example.com,http://www.example.com")
	flag.StringVar(&icmp, "icmp", "", "icmp扫描,探测内网某子网段存活，例：-icmp 192.168.*.1,-icmp 172.*.*.1")
	flag.StringVar(&port, "port", "top100", "输入需要扫描的端口,支持如下参数：full（全端口扫描）,top100,top1000 ,HttPorts(常见http端口）,1-65535（自定义端口范围）")
	flag.StringVar(&filePath, "hf", "", "指定需扫描的IP/域名/url 文件路径")
	flag.StringVar(&outPut, "output", "output.txt", "导出扫描结果到指定文件")
	flag.BoolVar(&dict, "dict", false, "是否对端口服务(ssh,ftp..)进行爆破，默认不爆破")
	flag.BoolVar(&path, "path", false, "是否进行目录爆破,默认不爆破")
	flag.BoolVar(&pocs, "pocs", false, "是否进行自动poc探测,默认不探测")
	//flag.StringVar(&subdomain, "domain", false, "是否爆破子域名，默认不爆破")
	flag.StringVar(&searchPoc, "s", "", "跳过指纹识别,对目标网址指定poc探测,例：-s shiro,seeyon,weblogic,thinkphp")
	flag.StringVar(&Plugins, "plugins", "", "针对url进行测试，可指定使用某个插件(401,jsfinder...),-plugins list 查看程序目前支持的所有插件详情")
	//flag.StringVar(&proxy, "proxy", "", "指定使用的代理 http://127.0.0.1:8080")
	flag.Parse()

	start := time.Now()

	File.CreateFile(outPut)

	if icmp != "" {
		exec.IcmpAlive(icmp, outPut)
	}

	//如只在命令行输入资产，则认为扫描资产数量不大
	if filePath == "" && host != "" {
		//File.CreateFile(outPut)
		Format.Choose(host, port, true, dict, outPut, path, pocs, searchPoc, Plugins)
	}
	if filePath != "" && host == "" {
		//File.CreateFile(outPut)
		ips, _ := File.ReadIPRangesFromFile(filePath)
		for _, host := range ips {
			//choose 一次只读取一条资产
			Format.Choose(host, port, true, dict, outPut, path, pocs, searchPoc, Plugins)
		}
	}
	if host == "" && filePath == "" {
		ps.Host = host
		Format.ParsePlugins(Plugins, ps)
		return
	}

	elapsed := time.Since(start)

	//输出执行时间。
	fmt.Println("全部扫描完成，共耗时: ", elapsed)

}
