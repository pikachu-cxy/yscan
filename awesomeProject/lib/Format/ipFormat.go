package Format

import (
	"awesomeProject/lib/File"
	p "awesomeProject/lib/Plugins"
	"awesomeProject/lib/exec"
	"awesomeProject/lib/pkg/afrog/config"
	afrogresult "awesomeProject/lib/pkg/afrog/result"
	afrogrunner "awesomeProject/lib/pkg/afrog/runner"
	"awesomeProject/lib/pkg/afrog/utils"
	"awesomeProject/lib/pkg/crack"
	httpxrunner "awesomeProject/lib/pkg/httpx/runner"
	"awesomeProject/lib/pkg/internals/crackrunner"
	_ "awesomeProject/lib/pkg/internals/crackrunner"
	"awesomeProject/lib/pkg/runner"
	"awesomeProject/lib/pkg/scan"
	"bufio"
	"crypto/tls"
	"encoding/binary"
	"fmt"
	crack2 "github.com/niudaii/crack/pkg/crack"
	"github.com/praetorian-inc/fingerprintx/pkg/plugins"
	"github.com/projectdiscovery/goflags"
	_ "github.com/projectdiscovery/goflags"
	"github.com/projectdiscovery/gologger"
	_ "github.com/projectdiscovery/utils/slice"
	"github.com/tidwall/gjson"
	"net"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type service struct {
	ip        string
	port      string
	protocol  int
	tls       bool
	transport string
	version   string
	metadata  []string
}

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

// ipFormat 判断ip格式
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
	pattern := `(^(\d{1,3}\.){3}\d{1,3})-\d{1,3}$`
	matched, _ := regexp.MatchString(pattern, s)
	if !matched {
		return "", false
	}
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
	//println(parts[0] + " 123")
	//println(ip)
	return ip, IpCompare(parts[0], ip2)
}

// 192.168.21.1-192.168.21.255
func IsIPRange(s string) ([]string, bool) {
	pattern := `(^(\d{1,3}\.){3}\d{1,3})-(\d{1,3}\.){3}\d{1,3}$`
	matched, _ := regexp.MatchString(pattern, s)
	if !matched {
		return nil, false
	}
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

// 比对两个ip地址是否按顺序大小输入 如192.168.1.1 < 192.168.1.2
func IpCompare(ip1 string, ip2 string) bool {

	// 将 IP 地址转换为整数
	ip1Int, err := ipToInt(ip1)
	if err != nil {
		return false
	}
	ip2Int, err := ipToInt(ip2)
	if err != nil {
		return false
	}

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

// 将ip地址转化为整数
func ipToInt(ip string) (uint32, error) {
	parts := strings.Split(ip, ".")
	a, err := strconv.Atoi(parts[0])
	if err != nil {
		return 1, err
	}
	b, err := strconv.Atoi(parts[1])
	if err != nil {
		return 1, err
	}
	c, err := strconv.Atoi(parts[2])
	if err != nil {
		return 1, err
	}
	d, err := strconv.Atoi(parts[3])
	if err != nil {
		return 1, err
	}
	return uint32((a << 24) | (b << 16) | (c << 8) | d), nil
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
		//fmt.Println("域名不匹配")
		return false
	}
	return false
}

func IsUrl(url string) bool {
	pattern := `[a-zA-z]+://[^\s]*` // 正则表达式模式

	matched, err := regexp.MatchString(pattern, url)

	if err != nil {
		//fmt.Println("正则匹配出错:", err)
		return false
	}

	if matched {
		//fmt.Println("url匹配成功")
		return true
	} else {
		//fmt.Println("url不匹配")
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
	if IsUrl(ip) {
		p = append(p, ip)
		return p, "url"
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

func ipIsAlive(host string, output string) bool {
	ip := net.ParseIP(host)
	//内网地址
	if ip.IsPrivate() {
		//内网判断方法 arp,udp常用端口，http常用端口
		if exec.OnePing(host, output) {
			return true
		} else {
			if exec.TcpScan(host, output) {
				return true
			}
		}
	} else {
		//外网判断方法 udp常用端口，http常用端口
		if exec.OnePing(host, output) {
			//ping 通 不继续下面逻辑，直接返回true
			return true
		} else {
			if exec.TcpScan(host, output) {
				return true
			}
		}
	}
	return false
}

func isHTTPPort(url string) bool {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr, Timeout: 2 * time.Second}
	uri := fmt.Sprintf("http://%s", url)
	resp, err := client.Head(uri)

	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return true
}

func Choose(host string, port string, w bool, dict bool, o string, poc bool, searchPoc string, Plugins string, noping bool) {
	hosts, format := ChooseFormat(host)
	var ps p.PluginService
	ps.Host = host
	switch strings.ToLower(format) {
	case "ip":
		//先判断是内网环境/外网环境
		//内网可以arp udp tcp http
		host = hosts[0]
		if noping || ipIsAlive(host, o) {
			portsMap, _ := exec.ParsePorts(port)
			inputs := exec.Scan(portsMap, host, w, o)
			//对系统端口进行指纹识别
			targetsList := make([]plugins.Target, 0)
			for _, input := range inputs {
				parsedTarget, _ := runner.ParseTarget(input)
				targetsList = append(targetsList, parsedTarget)
			}
			println("正在进行指纹识别~ 请稍等------------------------------")
			//fast模式 crackrunner.CreateScanConfigFast()
			results, _ := scan.ScanTargets(targetsList, scan.Config(runner.CreateScanConfig()))
			datas, _ := runner.Report(results)

			for _, data := range datas {
				File.WriteFile(o, data+"\n")
			}
			//浪费了性能 应该识别为http后传给httpx
			for _, input := range inputs {
				//http指纹识别 or poc探测
				if isHTTPPort(input) {
					if Plugins != "" {
						ps.Host = fmt.Sprintf("http://%s", input)
						ParsePlugins(Plugins, ps)
					}
					checkData(input, o, poc, searchPoc)
				}
			}
			if dict {
				brute(datas, "user.txt", "pass.txt")
			} else {
				//brute(datas, "", "")
			}
		}
	case "ips":
		//为了保证扫描效率，当无法ping通目标ip，则认为不存活
		targetsList := make([]plugins.Target, 0)
		inputs := make([]string, 0)
		if noping {
			portsMap, _ := exec.ParsePorts(port)
			for _, host := range hosts {
				inputs = exec.Scan(portsMap, host, w, o)
				//finger识别开始
				for _, input := range inputs {
					//println(input)
					parsedTarget, _ := runner.ParseTarget(input)
					targetsList = append(targetsList, parsedTarget)
				}
			}
		} else {
			ipAlive := exec.IpIcmp(hosts, o)

			portsMap, _ := exec.ParsePorts(port)
			for _, host := range ipAlive {
				inputs = exec.Scan(portsMap, host, w, o)
				//finger识别开始
				for _, input := range inputs {
					//println(input)
					parsedTarget, _ := runner.ParseTarget(input)
					targetsList = append(targetsList, parsedTarget)
				}
			}
		}
		for _, input := range inputs {
			if isHTTPPort(input) {
				if Plugins != "" {
					ps.Host = fmt.Sprintf("http://%s", input)
					ParsePlugins(Plugins, ps)
				}
				checkData(input, o, poc, searchPoc)
			}
		}

		println("正在进行系统端口指纹识别~ 请稍等------------------------------")
		//fast模式 crackrunner.CreateScanConfigFast()
		results, _ := scan.ScanTargets(targetsList, scan.Config(runner.CreateScanConfigFast()))
		datas, _ := runner.Report(results)
		for _, data := range datas {
			File.WriteFile(o, data+"\n")
		}
		if dict {
			brute(datas, "user.txt", "pass.txt")
		} else {
			//brute(datas, "", "")
		}
	case "domain":
		host := hosts[0]
		//如果对方禁ping 通过DNS解析判断存活
		if exec.OnePing(host, o) || exec.DnsLookup(host) {
			//扫描端口
			portsMap, _ := exec.ParsePorts(port)
			targetsList := make([]plugins.Target, 0)
			//input 格式www.baidu.com:110
			inputs := exec.Scan(portsMap, host, w, o)
			for _, input := range inputs {
				parsedTarget, _ := runner.ParseTarget(input)
				targetsList = append(targetsList, parsedTarget)
			}
			//http 指纹识别所需时间较长 1*time.Millisecond 不足
			results, _ := scan.ScanTargets(targetsList, scan.Config(runner.CreateScanConfigFast()))
			datas, _ := runner.Report(results)
			for _, data := range datas {
				println(data)
			}
			if Plugins != "" {
				ParsePlugins(Plugins, ps)
			}
			checkData(host, o, poc, searchPoc)

		}

	case "url":
		//访问连通性--指纹识别--poc探测
		//todo js爬取 深度扫描
		//WebFinger(host)
		//支持的插件
		if Plugins != "" {
			ParsePlugins(Plugins, ps)
		}

		if searchPoc != "" {
			WebPoc(host, searchPoc, o, false)
			return
		}
		httpRunner(hosts, o)
		//technologies []string todo poc扫描
		technologies := httpxrunner.Techs
		var techs string
		for _, tech := range technologies {
			techs = tech + "," + techs
		}
		techs = strings.TrimRight(techs, ",")
		if poc {
			WebPoc(host, techs, o, false)
		}
	}
}

func ParseTech(searchPoc string) []string {
	techs := strings.Split(strings.ToLower(searchPoc), ",")
	return techs
}

func brutePath(host string, hosts goflags.StringSlice, o string, path string, threads int) {
	var wg sync.WaitGroup
	urls := make([]string, 0)
	urls = append(urls, host)
	requests := make(chan string)
	//排除了404和400状态码显示
	lines, _ := ReadLinesFromFile(path)
	for i := 0; i < threads; i++ {
		//wg.Add(1)
		go func() {
			//defer wg.Done()
			for line := range requests {
				options := httpxrunner.ParseOptions(hosts, o, line)
				//println(options)
				httpxRunner, _ := httpxrunner.New(options)
				httpxRunner.RunEnumeration()
				httpxRunner.Close()
				wg.Done()
			}
		}()
	}
	// 将目录添加到请求通道中
	for _, dir := range lines {
		wg.Add(1)
		requests <- dir
	}
	close(requests)
	wg.Wait()
}

func ReadLinesFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

func checkData(data string, o string, poc bool, searchPoc string) {
	//urls := make([]string, 0)

	if searchPoc != "" {
		if searchPoc == "list" {
			return
		} else {
			WebPoc(data, searchPoc, o, false)
			return
		}
	}

	urls := make([]string, 0)
	urls = append(urls, data)
	httpRunner(urls, o)
	if poc {
		technologies := httpxrunner.Techs
		if technologies != nil {
			var techs string
			for _, tech := range technologies {
				techs = tech + "," + techs
			}
			techs = strings.TrimRight(techs, ",")
			WebPoc(data, techs, o, false)
		}
	}
}

func httpRunner(hosts []string, o string) {
	options := httpxrunner.ParseOptions(hosts, o, "")
	//println(options)
	httpxRunner, _ := httpxrunner.New(options)
	httpxRunner.RunEnumeration()
	httpxRunner.Close()
}

func httpRunnerBrute(hosts []string, o string, path string) {
	options := httpxrunner.ParseOptions(hosts, o, path)
	//println(options)
	httpxRunner, _ := httpxrunner.New(options)
	httpxRunner.RunEnumeration()
	httpxRunner.Close()
}

func brute(datas []string, userDict string, passDict string) {
	for _, data := range datas {
		//println(data)
		//if strings.Contains(data, "Name") {
		//	println(data)
		//}
		// gjson解析指纹识别结果
		ip := gjson.Get(data, "ip")
		port := gjson.Get(data, "port")
		protocol := gjson.Get(data, "protocol")
		uri := ""
		/*
			if protocol.String() == "http" {

			}

		*/
		if crack.SupportProtocols[protocol.String()] {
			//整理组合结果为127.0.0.1:3306|mysql 格式，传给crack
			uri = ip.String() + ":" + port.String() + "|" + protocol.String()
		} else {
			continue
		}
		println("正在进行端口服务爆破~ 请稍等------------------------------")
		//密码爆破 指纹识别到单个结果 即开始爆破
		if protocol.String() == "redis" || protocol.String() == "memcached" {
			options := crackrunner.Options{Input: uri}
			//fmt.Printf("%v", options)
			option := crackrunner.ParseOptions(&options)
			//fmt.Printf("%v", option)
			//设置爆破参数 线程 超时
			crackOptions := setOptions(50, 1)
			newRunner, err := crackrunner.NewRunner(option, crack2.Options(crackOptions))
			if err != nil {
				gologger.Fatal().Msgf("Could not create runner: %v", err)
			}
			newRunner.Run(protocol.String())
		} else {
			options := crackrunner.Options{Input: uri, UserFile: userDict, PassFile: passDict}
			//fmt.Printf("%v", options)
			option := crackrunner.ParseOptions(&options)
			//fmt.Printf("%v", option)
			//设置爆破参数 线程 超时
			crackOptions := setOptions(50, 1)
			newRunner, err := crackrunner.NewRunner(option, crack2.Options(crackOptions))
			if err != nil {
				gologger.Fatal().Msgf("Could not create runner: %v", err)
			}
			newRunner.Run(protocol.String())
		}
	}
}

func setOptions(thread int, timeout int) crack.Options {
	crackOptions := crack.Options{
		Threads:  thread,
		Timeout:  timeout,
		Delay:    0,
		CrackAll: false,
		Silent:   false,
	}
	return crackOptions
}

func WebPoc(host string, techs string, output string, pl bool) {

	options, err := config.NewOptions(host, techs, pl)
	r, err := afrogrunner.NewRunner(options)
	if err != nil {
		gologger.Error().Msg(err.Error())
		os.Exit(0)
	}
	if err != nil {
		gologger.Error().Msgf("Could not create runner: %s\n", err)
		os.Exit(0)
	}
	var (
		lock   = sync.Mutex{}
		number uint32
	)
	r.OnResult = func(result *afrogresult.Result) {
		if result.IsVul {
			lock.Lock()
			atomic.AddUint32(&number, 1)
			result.PrintColorResultInfoConsole(utils.GetNumberText(int(number)), output)

			if !options.DisableOutputHtml {
				r.Report.SetResult(result)
				r.Report.Append(utils.GetNumberText(int(number)))
			}
			/*
				if len(options.Json) > 0 || len(options.JsonAll) > 0 {
					r.JsonReport.SetResult(result)
					r.JsonReport.Append()
				}
			*/
			lock.Unlock()
		}
	}
	if err := r.Run(); err != nil {
		gologger.Error().Msgf("poc runner run err: %s\n", err)
		os.Exit(0)
	}
}
