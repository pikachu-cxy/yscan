package crackrunner

import (
	"awesomeProject/lib/File"
	"fmt"
	"github.com/niudaii/crack/pkg/crack"
	"github.com/projectdiscovery/gologger"
)

type Runner struct {
	Options     *Options
	CrackRunner *crack.Runner
}

func NewRunner(options *Options, crackOptions crack.Options) (*Runner, error) {

	crackRunner, err := crack.NewRunner(&crackOptions)
	if err != nil {
		return nil, fmt.Errorf("crack.NewRunner() err, %v", err)
	}
	return &Runner{
		Options:     options,
		CrackRunner: crackRunner,
	}, nil
}

func (r *Runner) Run(module string) {
	// 解析目标
	addrs := crack.ParseTargets(r.Options.Targets)

	//addrs = crack.FilterModule(addrs, r.options.Module)
	addrs = crack.FilterModule(addrs, module)
	if len(addrs) == 0 {
		gologger.Info().Msgf("目标为空")
		return
	}
	// 存活探测
	//gologger.Info().Msgf("存活探测")
	addrs = r.CrackRunner.CheckAlive(addrs)
	//gologger.Info().Msgf("存活数量: %v", len(addrs))

	// 服务爆破
	results := r.CrackRunner.Run(addrs, r.Options.UserDict, r.Options.PassDict)
	if len(results) > 0 {
		gologger.Info().Msgf("爆破成功: %v", len(results))
		for _, result := range results {
			gologger.Print().Msgf("[+]%v -> %v %v", result.Protocol, result.Addr, result.UserPass)
			//将爆破成功结果写入文件
			File.WriteFile("output.txt", "[+]"+result.Protocol+" -> "+result.Addr+result.UserPass+"\n")
		}
	}
}
