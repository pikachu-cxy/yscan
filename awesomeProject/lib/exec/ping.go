package exec

import (
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"strconv"
	"sync"
)

var number int
var output string
var ipAddress []string
var ipAlive []string

//var ip string

func OutputSet(output1 string) {
	output = output1
}
func OnePing(host string) bool {
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
		fmt.Printf("%s is can't connect！\n", host)
		return false
	} else {
		fmt.Printf("%s is alive！\n", host)
	}
	return true
}

func somePing(host string, wg *sync.WaitGroup) {
	defer wg.Done()
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
		//return false
	} else {
		number++
		ipAlive = append(ipAlive, host)
		//File.WriteFile(output, host+"  is alive!\n")
		fmt.Printf("%s is alive！\n", host)
	}
}

func IpIcmp(ips []string) []string {
	ipAlive = []string{}
	var wg sync.WaitGroup
	number = 0
	if ips != nil {
		for _, v := range ips {
			wg.Add(1)
			go somePing(v, &wg)
		}
		wg.Wait()
		fmt.Printf("ip段存活ip数量为：" + strconv.Itoa(number) + "\n")

		//File.WriteFile(output, "段存活ip数量为："+strconv.Itoa(number)+"\n")

	}
	return ipAlive
}

func DomainIcmp(ips []string) []string {
	ipAlive = []string{}
	var wg sync.WaitGroup
	number = 0
	if ips != nil {
		for _, v := range ips {
			wg.Add(1)
			go OnePing(v)
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

//内网arp探测
func arpConn(targetIPs []string) {

}
