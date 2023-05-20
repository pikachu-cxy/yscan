package exec

import (
	"awesomeProject/lib/Format"
	"fmt"
	"os/exec"
)

func IpIcmp(ip *string) {

	ips, _ := sometone(*ip)
	if ips != nil {
		for _, v := range ips {
			//fmt.Printf("%d and %s",i,v)
			cmd := exec.Command("ping", v, "-n", "1", "-w", "200") // Windows系统上使用 "-n" 参数指定发送的 ICMP 请求次数
			err := cmd.Run()
			if err != nil {
				fmt.Printf("执行命令出错：%v\n", err)
				return
			} else {
				fmt.Printf("%s is alive！", *ip)
			}

		}
	}

}

func sometone(ip string) (s []string, bool2 bool) {
	ips := Format.ChooseFormat(ip)

	if ips == nil {
		return nil, false
	}
	return ips, true
}
