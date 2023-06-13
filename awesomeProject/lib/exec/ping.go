package exec

import (
	"awesomeProject/lib/File"
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

var number int
var output string
var ipAddress []string
var ipAlive []string

//var ip string

func OutputSet(output1 string) {
	output = output1
}
func OnePing(host string, output string) bool {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("ping", host, "-n", "1", "-w", "200")
	case "linux":
		cmd = exec.Command("ping", host, "-c", "1", "-W", "1")
	case "darwin":
		cmd = exec.Command("ping", host, "-c", "1", "-W", "200")
	case "freebsd":
		cmd = exec.Command("ping", "-c", "1", "-W", "200", host)
	case "openbsd":
		cmd = exec.Command("ping", "-c", "1", "-w", "200", host)
	case "netbsd":
		cmd = exec.Command("ping", "-c", "1", "-w", "2", host)
	default:
		cmd = exec.Command("ping", "-c", "1", host)
	}
	err := cmd.Run()
	if err != nil {
		return false
	} else {
		File.WriteFile(output, host+"  is alive!\n")
		fmt.Printf("%s is alive！\n", host)
	}
	return true
}

func somePing(host string, wg *sync.WaitGroup, output string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("ping", host, "-n", "1", "-w", "200")
	case "linux":
		cmd = exec.Command("ping", host, "-c", "1", "-W", "1")
	case "darwin":
		cmd = exec.Command("ping", host, "-c", "1", "-W", "200")
	case "freebsd":
		cmd = exec.Command("ping", "-c", "1", "-W", "200", host)
	case "openbsd":
		cmd = exec.Command("ping", "-c", "1", "-w", "200", host)
	case "netbsd":
		cmd = exec.Command("ping", "-c", "1", "-w", "2", host)
	default:
		cmd = exec.Command("ping", "-c", "1", host)
	}
	err := cmd.Run()
	if err != nil {
		if TcpScan(host, output) {
			number++
			ipAlive = append(ipAlive, host)
			//File.WriteFile(output, host+"  is alive!\n")
			//fmt.Printf("%s is alive！\n", host)
		}
	} else {
		number++
		ipAlive = append(ipAlive, host)
		File.WriteFile(output, host+"  is alive!\n")
		fmt.Printf("%s is alive！\n", host)
	}
	defer wg.Done()
}

func IpIcmp(ips []string, o string) []string {
	ipAlive = []string{}
	var wg sync.WaitGroup
	number = 0
	if ips != nil {
		for _, v := range ips {
			wg.Add(1)
			go somePing(v, &wg, o)
		}
		wg.Wait()

		fmt.Printf(ips[0] + " ip段存活ip数量为：" + strconv.Itoa(number) + "\n")

		//File.WriteFile(output, "段存活ip数量为："+strconv.Itoa(number)+"\n")

	}
	return ipAlive
}
func GateAwayIcmp(ips []string, output string) {
	var wg sync.WaitGroup
	task := make(chan string, 0)
	if ips != nil {
		for i := 0; i < 1000; i++ {
			go func() {
				for v := range task {
					ip := net.ParseIP(v)
					//只探测内网地址
					if ip.IsPrivate() {
						var cmd *exec.Cmd
						switch runtime.GOOS {
						case "windows":
							cmd = exec.Command("ping", v, "-n", "1", "-w", "400")
						case "linux":
							cmd = exec.Command("ping", v, "-c", "1", "-W", "1")
						case "darwin":
							cmd = exec.Command("ping", v, "-c", "1", "-W", "400")
						case "freebsd":
							cmd = exec.Command("ping", "-c", "1", "-W", "200", v)
						case "openbsd":
							cmd = exec.Command("ping", "-c", "1", "-w", "200", v)
						case "netbsd":
							cmd = exec.Command("ping", "-c", "1", "-w", "2", v)
						default:
							cmd = exec.Command("ping", "-c", "1", v)
						}
						err := cmd.Run()
						if err != nil {
							if TcpScan(v, output) {
								number++
								ipAlive = append(ipAlive, v)
								//File.WriteFile(output, host+"  is alive!\n")
								//fmt.Printf("%s is alive！\n", host)
							}
						} else {
							number++
							ipAlive = append(ipAlive, v)
							File.WriteFile(output, v+"  is alive!\n")
							fmt.Printf("%s is alive！\n", v)
						}
					}
					wg.Done()
				}
			}()
		}
		for _, ip := range ips {
			wg.Add(1)
			task <- ip
		}
		close(task)
		wg.Wait()
	}
}

func DomainIcmp(ips []string, o string) []string {
	ipAlive = []string{}
	var wg sync.WaitGroup
	number = 0
	if ips != nil {
		for _, v := range ips {
			wg.Add(1)
			go OnePing(v, o)
		}
		wg.Wait()
		//File.WriteFile(output, "段存活ip数量为："+strconv.Itoa(number)+"\n")
	}
	return ipAlive
}

func DnsLookup(host string) bool {
	_, err := net.LookupIP(host)
	if err != nil {
		fmt.Printf("DNS lookup failed: %v\n", err)
		return false
	}

	fmt.Printf("Exist IP addresses for %s:\n", host)
	return true
}

func TcpScan(host string, output string) bool {
	ports := []int{21, 22, 80, 135, 139, 443, 445, 3389, 8080}
	for _, port := range ports {

		address := fmt.Sprintf("%s:%d", host, port)

		conn, err := net.DialTimeout("tcp", address, time.Second)
		if err == nil {
			ms := "[+tcp] " + host + ":" + strconv.Itoa(port) + " " + "is alive! (目标可能禁用了icmp协议)"
			File.WriteFile(output, ms+"\n")
			println(ms)
			conn.Close()
			return true
		}
	}
	return false
}

func IcmpAlive(ip string, o string) {
	GateAwayList := GetGateAway(ip)
	GateAwayIcmp(GateAwayList, o)
	//IpIcmp(GateAwayList, o)
}

func GetGateAway(ip string) []string {
	GateAwayList := make([]string, 0)
	GateAwayListOne := make([]bool, 3)
	strArr := strings.Split(ip, ".")
	for index := 2; index > 0; index-- {
		if strArr[index] == "*" && index == 2 {
			for i := 0; i < 255; i++ {
				GateAwayList = append(GateAwayList, fmt.Sprintf("%s.%s.%d.1", strArr[0], strArr[1], i))
				GateAwayList = append(GateAwayList, fmt.Sprintf("%s.%s.%d.254", strArr[0], strArr[1], i))
			}
			GateAwayListOne[2] = true
		}
		if strArr[index] == "*" && index == 1 {
			if GateAwayListOne[2] {
				GateAwayList = nil
				for j := 0; j < 255; j++ {
					for i := 0; i < 255; i++ {
						GateAwayList = append(GateAwayList, fmt.Sprintf("%s.%d.%d.1", strArr[0], i, j))
						GateAwayList = append(GateAwayList, fmt.Sprintf("%s.%d.%d.254", strArr[0], i, j))
					}
				}
			} else {
				for i := 0; i < 255; i++ {
					GateAwayList = append(GateAwayList, fmt.Sprintf("%s.%d.%s.1", strArr[0], i, strArr[2]))
					GateAwayList = append(GateAwayList, fmt.Sprintf("%s.%d.%s.254", strArr[0], i, strArr[2]))
				}
			}
			GateAwayListOne[1] = true
		}
		if strArr[index] == "*" && index == 0 {
			//GateAwayListOne[0] = true
			//192 ,172, 10
			var adds = [3]int{192, 172, 10}
			if GateAwayListOne[2] && GateAwayListOne[1] {
				GateAwayList = nil
				for j := 0; j < 255; j++ {
					for i := 0; i < 255; i++ {
						for _, add := range adds {
							GateAwayList = append(GateAwayList, fmt.Sprintf("%d.%d.%d.1", add, i, j))
							GateAwayList = append(GateAwayList, fmt.Sprintf("%d.%d.%d.254", add, i, j))
						}
					}
				}
			} else if GateAwayListOne[2] {
				GateAwayList = nil
				for i := 0; i < 255; i++ {
					for _, add := range adds {
						GateAwayList = append(GateAwayList, fmt.Sprintf("%d.%s.%d.1", add, strArr[1], i))
						GateAwayList = append(GateAwayList, fmt.Sprintf("%d.%s.%d.254", add, strArr[1], i))
					}
				}
			} else if GateAwayListOne[1] {
				//println(123)
				GateAwayList = nil
				for i := 0; i < 255; i++ {
					for _, add := range adds {
						GateAwayList = append(GateAwayList, fmt.Sprintf("%d.%d.%s.1", add, i, strArr[2]))
						GateAwayList = append(GateAwayList, fmt.Sprintf("%d.%d.%s.254", add, i, strArr[2]))
					}
				}

			} else {
				for _, add := range adds {
					GateAwayList = append(GateAwayList, fmt.Sprintf("%d.%s.%s.1", add, strArr[1], strArr[2]))
					GateAwayList = append(GateAwayList, fmt.Sprintf("%d.%s.%s.254", add, strArr[1], strArr[2]))
				}
			}
		}

	}
	return GateAwayList
}

// udp
func UdpScan(host string, output string) bool {
	ports := []int{53, 69, 123, 161, 514, 1900, 4500, 5353}

	for _, port := range ports {
		/*
			// 构建 UDP 地址
			address := &net.UDPAddr{
				IP:   net.ParseIP(host),
				Port: port,
			}

		*/
		address := fmt.Sprintf("%s:%d", host, port)

		conn, _ := net.DialTimeout("udp", address, time.Second)
		if conn != nil {
			// 发送测试数据
			_, err := conn.Write([]byte("hello"))
			if err != nil {
				return false
			}
			// 设置超时时间
			conn.SetReadDeadline(time.Now().Add(time.Second))
			// 接收响应
			buf := make([]byte, 1024)
			_, err = conn.Read(buf)
			if err != nil {
				// 判断错误类型
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					// 超时表示端口开放
					ms := "[+udp] " + host + ":" + strconv.Itoa(port) + " " + "is alive! (目标禁用了icmp协议)"
					File.WriteFile(output, ms+"\n")
					println(ms)
					return true
				}
				// 其他错误表示端口关闭
			}
		}
		conn.Close()
	}
	return false
}
