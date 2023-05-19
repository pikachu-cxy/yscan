package exec

import (
	"fmt"
	"os/exec"
)

func IpIcmp(ip *string) {

	cmd := exec.Command("ping", *ip, "-n", "1", "-w", "200") // Windows系统上使用 "-n" 参数指定发送的 ICMP 请求次数
	err := cmd.Run()
	if err != nil {
		fmt.Printf("执行命令出错：%v\n", err)
		return
	} else {
		fmt.Printf("%s is alive！", *ip)
	}
}
