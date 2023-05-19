package Format

import (
	"fmt"
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
//ip范围  192.168.21.1-255 ||192.168.21.1-192.168.21.255
//ip 192.168.21.1

// IpCIDRFormat 检测ip段格式 192.168.21.1/24
func IpCIDRFormat(host string) ([]string, bool) {

	ip, ipNet, err := net.ParseCIDR(host)
	if err != nil {
		fmt.Println("无效的 IP 段:", err)
		return []string{}, false
	}

	// 获取 IP 地址范围的开始和结束 IP
	startIP := ip.Mask(ipNet.Mask)
	endIP := make(net.IP, len(startIP))
	copy(endIP, startIP)

	for i := range endIP {
		endIP[i] |= ^ipNet.Mask[i]
	}

	fmt.Println("IP 地址范围：")
	fmt.Println("开始 IP:", startIP)
	fmt.Println("结束 IP:", endIP)
	return []string{}, true
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
func IsIPRange2(s string) bool {
	// 以连字符分割字符串，判断是否有两个 IP 地址
	parts := strings.Split(s, "-")
	if len(parts) != 2 {
		return false
	}
	ip1 := net.ParseIP(parts[0])

	if ip1 == nil {
		return false
	}
	// 获取 IP 地址的前三位
	ip2 := parts[0][:strings.LastIndex(parts[0], ".")+1] + parts[1]
	if net.ParseIP(ip2) == nil {
		return false
	}
	return ipCompare(parts[0], parts[1])
}

//192.168.21.1-192.168.21.255
func IsIPRange(s string) bool {
	// 以连字符分割字符串，判断是否有两个 IP 地址
	parts := strings.Split(s, "-")
	if len(parts) != 2 {
		return false
	}

	// 判断两个 IP 地址的合法性
	ip1 := net.ParseIP(parts[0])
	ip2 := net.ParseIP(parts[1])
	if ip1 == nil || ip2 == nil {
		return false
	}
	return ipCompare(parts[0], parts[1])
	//return true
}

//比对两个ip地址是否按顺序大小输入 如192.168.1.1 < 192.168.1.2
func ipCompare(ip1 string, ip2 string) bool {

	// 将 IP 地址转换为整数
	ip1Int := ipToInt(ip1)
	ip2Int := ipToInt(ip2)

	// 比较 IP 地址的整数值
	if ip1Int < ip2Int {

		fmt.Println(ip1, "在", ip2, "之前")
		return true
	} else if ip1Int > ip2Int {
		fmt.Println(ip1, "在", ip2, "之后")
		return false
	} else {
		fmt.Println(ip1, "和", ip2, "相同")
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
