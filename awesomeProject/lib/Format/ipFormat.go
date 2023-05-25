package Format

import (
	"awesomeProject/lib/exec"
	"awesomeProject/lib/pkg/internals/crackrunner"
	_ "awesomeProject/lib/pkg/internals/crackrunner"
	"awesomeProject/lib/pkg/runner"
	"encoding/binary"
	"fmt"
	"github.com/praetorian-inc/fingerprintx/pkg/plugins"
	"github.com/praetorian-inc/fingerprintx/pkg/scan"
	"net"
	"regexp"
	"strconv"
	"strings"
)

type IP interface {
	readIP() ([]byte, error)
}

//支持三种ip格式
//文本格式  暂定
//ip段格式 192.168.21.1/24
//ip范围  192.168.21.1-255 ||192.168.21.1-192.168.21.255(最后都转换为这种格式）
//ip 192.168.21.1

// IpCIDRFormat 检测ip段格式 192.168.21.1/24
func IpCIDRFormat(host string) ([]string, bool) {

	_, ipNet, err := net.ParseCIDR(host)
	if err != nil {
		//fmt.Println("无效的 IP 段:", err)
		return nil, false
	}
	// 获取IP地址列表
	ipList := make([]string, 0)
	for ip := ipNet.IP.Mask(ipNet.Mask); ipNet.Contains(ip); inc(ip) {
		ipList = append(ipList, net.IPv4(ip[0], ip[1], ip[2], ip[3]).String())
	}

	//fmt.Println(startIP.String() + "-" + endIP.String())
	return ipList, true
}

//ipFormat 判断ip格式
func ipFormat(host string) bool {
	pattern := `^(\d{1,3}\.){3}\d{1,3}$`

	matched, _ := regexp.MatchString(pattern, host)

	if matched {
		return true
	} else {
		return false
	}
}

// IsIPRange2 192.168.21.1-255
func IsIPRange2(s string) (string, bool) {
	// 以连字符分割字符串，判断是否有两个 IP 地址
	parts := strings.Split(s, "-")
	if len(parts) != 2 {
		return "", false
	}
	ip1 := net.ParseIP(parts[0])

	if ip1 == nil {
		return "", false
	}
	// 获取 IP 地址的前三位
	ip2 := parts[0][:strings.LastIndex(parts[0], ".")+1] + parts[1]
	if net.ParseIP(ip2) == nil {
		return "", false
	}
	ip := parts[0] + "-" + ip2
	println(parts[0] + " 123")
	println(ip)
	return ip, IpCompare(parts[0], ip2)
}

//192.168.21.1-192.168.21.255
func IsIPRange(s string) ([]string, bool) {

	// 以连字符分割字符串，判断是否有两个 IP 地址
	parts := strings.Split(s, "-")
	if len(parts) != 2 {
		return nil, false
	}
	if IpCompare(parts[0], parts[1]) == false {
		println("ip地址前后值不符合格式！")
	}

	// 判断两个 IP 地址的合法性
	start := net.ParseIP(parts[0])
	end := net.ParseIP(parts[1])
	if start == nil || end == nil {
		println("ip地址格式不正确！")
		return nil, false
	}
	// 将起始IP地址转换为无符号32位整数
	startInt := binary.BigEndian.Uint32(start.To4())

	// 将结束IP地址转换为无符号32位整数
	endInt := binary.BigEndian.Uint32(end.To4())
	// 获取IP地址列表
	ipList := make([]string, 0)
	for i := startInt; i <= endInt; i++ {
		ip := make(net.IP, 4)
		binary.BigEndian.PutUint32(ip, i)
		ipList = append(ipList, ip.String())
	}

	return ipList, true

}

//比对两个ip地址是否按顺序大小输入 如192.168.1.1 < 192.168.1.2
func IpCompare(ip1 string, ip2 string) bool {

	// 将 IP 地址转换为整数
	ip1Int := ipToInt(ip1)
	ip2Int := ipToInt(ip2)

	// 比较 IP 地址的整数值
	if ip1Int < ip2Int {

		//fmt.Println(ip1, "在", ip2, "之前   ")
		return true
	} else if ip1Int > ip2Int {
		fmt.Print(ip1, "在", ip2, "之后   ")
		return false
	} else {
		fmt.Print(ip1, "和", ip2, "相同   ")
		return false
	}

}

//将ip地址转化为整数
func ipToInt(ip string) uint32 {
	parts := strings.Split(ip, ".")
	a, _ := strconv.Atoi(parts[0])
	b, _ := strconv.Atoi(parts[1])
	c, _ := strconv.Atoi(parts[2])
	d, _ := strconv.Atoi(parts[3])
	return uint32((a << 24) | (b << 16) | (c << 8) | d)
}

func parseIPRange(ipRange string) []string {

	startIP, endIP := strings.Split(ipRange, "-")[0], strings.Split(ipRange, "-")[1]

	startIps := strings.Split(startIP, ".")
	endIps := strings.Split(endIP, ".")

	start, _ := strconv.Atoi(startIps[3])
	end, _ := strconv.Atoi(endIps[3])

	ipList := make([]string, end-start+1)

	for i := 0; i < len(ipList); i++ {
		ipList[i] = fmt.Sprintf("%s.%s.%s.%d", startIps[0], startIps[1], startIps[2], start+i)
	}

	return ipList
}

func IsDomainRange(domain string) bool {

	pattern := `^[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$` // 正则表达式模式

	matched, err := regexp.MatchString(pattern, domain)

	if err != nil {
		fmt.Println("域名正则匹配出错:", err)
		return false
	}

	if matched {
		//fmt.Println("域名匹配成功")
		return true
	} else {
		fmt.Println("域名不匹配")
		return false
	}
	return false
}

func IsUrl(url string) bool {
	pattern := `[a-zA-z]+://[^\s]*` // 正则表达式模式

	matched, err := regexp.MatchString(pattern, url)

	if err != nil {
		fmt.Println("正则匹配出错:", err)
		return false
	}

	if matched {
		fmt.Println("域名匹配成功")
		return true
	} else {
		fmt.Println("域名不匹配")
		return false
	}
	return false

}

func ChooseFormat(ip string) (s []string, s2 string) {

	p := make([]string, 0)
	if ipFormat(ip) {
		p = append(p, ip)
		return p, "ip"
	}
	ipscope, ipbool := IpCIDRFormat(ip)
	if ipbool {
		return ipscope, "ips"
	}
	ipscope, ipbool = IsIPRange(ip)
	if ipbool {
		return ipscope, "ips"
	}
	ips, ipbool1 := IsIPRange2(ip)
	if ipbool1 {
		ipscope, ipbool = IsIPRange(ips)
		return ipscope, "ips"
	}
	if IsDomainRange(ip) {
		p = append(p, ip)
		return p, "domain"
	}
	if IsUrl(ip) {
		p = append(p, ip)
		return p, "url"
	}

	return []string{}, ""
}

// 增加IP地址
func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func Choose(host string, port string, w bool) {

	hosts, format := ChooseFormat(host)
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
			inputs := exec.ScanPort(portsMap, host, w)
			//对系统端口进行指纹识别
			targetsList := make([]plugins.Target, 0)
			for _, input := range inputs {
				parsedTarget, _ := runner.ParseTarget(input)
				targetsList = append(targetsList, parsedTarget)
			}
			//fast模式 crackrunner.CreateScanConfigFast()
			results, _ := scan.ScanTargets(targetsList, scan.Config(runner.CreateScanConfig()))
			datas, _ := runner.Report(results)
			for _, data := range datas {
				println(data)
				//对识别的端口服务（ssh，mysql等）进行爆破
				options := &crackrunner.Options{Input: data, User: "root", Pass: "123456"}
				newRunner, _ := crackrunner.NewRunner(options)
				newRunner.Run()
			}
		}
	case "ips":
		//为了保证扫描效率，当无法ping通目标ip，则认为不存活
		ipAlive := exec.IpIcmp(hosts)
		portsMap, _ := exec.ParsePorts(port)
		targetsList := make([]plugins.Target, 0)
		for _, host := range ipAlive {
			inputs := exec.ScanPort(portsMap, host, w)
			for _, input := range inputs {
				parsedTarget, _ := runner.ParseTarget(input)
				targetsList = append(targetsList, parsedTarget)
			}
		}
		//fast模式 crackrunner.CreateScanConfigFast()
		results, _ := scan.ScanTargets(targetsList, scan.Config(runner.CreateScanConfigFast()))
		runner.Report(results)

	case "domain":
		host := hosts[0]
		//如果对方禁ping 通过DNS解析判断存活
		if exec.OnePing(host) || exec.DnsLookup(host) {
			//扫描端口
		}
		//爆破域名 深度扫描

	case "url":
		//访问连通性--指纹识别--poc探测
		//js爬取 深度扫描
	}
}
