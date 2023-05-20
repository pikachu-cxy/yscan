package exec

import (
	"awesomeProject/lib/Format"
	"fmt"
	"os/exec"
	"runtime"
	"sync"
)

func osping(host string, wg *sync.WaitGroup) {
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
		fmt.Printf("%s is alive！\n", host)
	}
}

func IpIcmp(ip *string) {

	var wg sync.WaitGroup

	ips := sometone(*ip)

	if ips != nil {
		for _, v := range ips {
			wg.Add(1)
			go osping(v, &wg)

		}
		wg.Wait()

	}

}

func sometone(ip string) []string {
	ips := Format.ChooseFormat(ip)

	if ips == nil {
		return nil
	}
	return ips
}
