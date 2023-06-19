package check

import (
	"awesomeProject/lib/pkg/poc"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"awesomeProject/lib/pkg/poc/common/errors"
	common_structs "awesomeProject/lib/pkg/poc/common/structs"
	"awesomeProject/lib/pkg/poc/utils"
	xray_structs "awesomeProject/lib/pkg/poc/xray/structs"
	nuclei_structs "github.com/WAY29/pocV/pkg/nuclei/structs"

	"github.com/panjf2000/ants"
)

var (
	EmptyLinks = []string{}

	Ticker  *time.Ticker
	Pool    *ants.PoolWithFunc
	Verbose bool

	WaitGroup sync.WaitGroup

	OutputChannel chan common_structs.Result

	ResultPool = sync.Pool{
		New: func() interface{} {
			return new(common_structs.PocResult)
		},
	}
)

// 初始化协程池
func InitCheck(threads, rate int, verbose bool) {
	var err error

	rateLimit := time.Second / time.Duration(rate)
	Ticker = time.NewTicker(rateLimit)
	Pool, err = ants.NewPoolWithFunc(threads, check)
	if err != nil {
		utils.CliError("Initialize goroutine pool error: "+err.Error(), 2)
	}

	Verbose = verbose
}

// 将任务放入协程池
func Start(target string, xrayPocMap []poc.Poc, nucleiPocMap map[string]nuclei_structs.Poc, outputChannel chan common_structs.Result) {
	// 设置outputChannel
	OutputChannel = outputChannel

	for _, poc1 := range xrayPocMap {
		WaitGroup.Add(1)
		Pool.Invoke(&xray_structs.Task{
			Target: target,
			Poc:    poc1,
		})
	}
	for _, poc := range nucleiPocMap {
		WaitGroup.Add(1)
		Pool.Invoke(&nuclei_structs.Task{
			Target: target,
			Poc:    poc,
		})
	}

}

// 等待协程池
func Wait() {
	WaitGroup.Wait()
}

// 释放协程池
func End() {
	Pool.Release()
}

// 核心代码，poc检测
func check(taskInterface interface{}) {
	var (
		oRequest *http.Request = nil

		err     error
		pocName string
	)

	defer WaitGroup.Done()
	<-Ticker.C

	switch taskInterface.(type) {
	case *xray_structs.Task:
		task, ok := taskInterface.(*xray_structs.Task)
		if !ok {
			wrappedErr := errors.Newf(errors.ConvertInterfaceError, "Can't convert task interface: %#v", err)
			utils.ErrorP(wrappedErr)
			return
		}
		target, poc := task.Target, task.Poc

		pocName = poc.Id
		if poc.Transport != "tcp" && poc.Transport != "udp" {
			oRequest, _ = http.NewRequest("GET", target, nil)
		}

		isVul, err := executeXrayPoc(oRequest, target, &xray_structs.Poc{})
		if err != nil {
			utils.ErrorP(err)
			return
		}

		pocResult := ResultPool.Get().(*common_structs.PocResult)
		pocResult.Str = fmt.Sprintf("%s (%s)", target, pocName)
		pocResult.Success = isVul
		pocResult.URL = target
		pocResult.PocName = poc.Id
		pocResult.PocLink = poc.Info.Reference
		pocResult.PocAuthor = poc.Info.Author
		pocResult.PocDescription = poc.Info.Description

		OutputChannel <- pocResult

	case *nuclei_structs.Task:
		var (
			desc    string
			author  string
			authors []string
		)

		task, ok := taskInterface.(*nuclei_structs.Task)
		if !ok {
			wrappedErr := errors.Newf(errors.ConvertInterfaceError, "Can't convert task interface: %#v", err)
			utils.ErrorP(wrappedErr)
			return
		}
		target, poc := task.Target, task.Poc
		authors, ok = poc.Info.Authors.Value.([]string)
		if !ok {
			author = "Unknown"
		} else {
			author = strings.Join(authors, ", ")
		}

		results, isVul, err := executeNucleiPoc(target, &poc)
		if err != nil {
			utils.ErrorP(err)
			return
		}

		for _, r := range results {
			if r.ExtractorName != "" {
				desc = r.TemplateID + ":" + r.ExtractorName
			} else if r.MatcherName != "" {
				desc = r.TemplateID + ":" + r.MatcherName
			}

			pocResult := ResultPool.Get().(*common_structs.PocResult)
			pocResult.Str = fmt.Sprintf("%s (%s) ", r.Matched, r.TemplateID)
			pocResult.Success = isVul
			pocResult.URL = r.Matched
			pocResult.PocName = r.TemplateID
			pocResult.PocLink = EmptyLinks
			pocResult.PocAuthor = author
			pocResult.PocDescription = desc

			OutputChannel <- pocResult
		}
	}

}

func PutPocResult(result *common_structs.PocResult) {
	result.Str = ""
	result.Success = false
	result.URL = ""
	result.PocName = ""
	result.PocLink = nil
	result.PocDescription = ""
	result.PocAuthor = ""

	ResultPool.Put(result)
}
