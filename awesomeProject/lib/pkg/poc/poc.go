package poc

import (
	"awesomeProject/lib/pkg/poc/common/check"
	"awesomeProject/lib/pkg/poc/common/output"
	xray_requests "awesomeProject/lib/pkg/poc/xray/requests"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

//const pocPath = "C:\\Users\\simple.chen\\Downloads\\yscan5\\awesomeProject\\pocs"

func Scan(targetURL string) {

	// 设置变量
	timeoutSecond := time.Duration(5) * time.Second

	// 初始化dnslog平台
	//common_structs.InitReversePlatform(*apiKey, *domain, timeoutSecond)
	//if common_structs.ReversePlatformType != xray_structs.ReverseType_Ceye {
	//	utils.WarningF("No Ceye api, use dnslog.cn")
	//}

	// 初始化http客户端
	xray_requests.InitHttpClient(10, "", timeoutSecond)

	// 计算xray的总发包量，初始化缓存
	xrayTotalReqeusts := 0

	// 加载poc
	//xrayPocs, nucleiPocs := LoadPocs(poc, pocPath)
	xrayPocs := getPocPath()
	for _, poc := range xrayPocs {
		ruleLens := len(poc.Rules)
		// 额外需要缓存connectionID
		if poc.Transport == "tcp" || poc.Transport == "udp" {
			ruleLens += 1
		}
		xrayTotalReqeusts += ruleLens
	}
	if xrayTotalReqeusts == 0 {
		xrayTotalReqeusts = 1
	}
	xray_requests.InitCache(xrayTotalReqeusts)

	// 初始化输出
	outputChannel, outputWg := output.InitOutput("", false, false)

	// 初始化check
	check.InitCheck(10, 100, false)

	// check开始
	check.Start(targetURL, xrayPocs, nil, outputChannel)
	check.Wait()

	// check结束
	close(outputChannel)
	check.End()
	outputWg.Wait()
}

func getPocPath() []Poc {
	dir, _ := os.Getwd()
	return loadTemplates(dir + "\\..\\..\\..\\pocs")
}

// 加载模板
func loadTemplates(pocpath string) []Poc {

	pocspath := FileForEachComplete(pocpath)
	pocs := make([]Poc, 0)
	// 遍历每个文件进行读取和解析
	for _, file := range pocspath {
		// 读取文件内容
		poc, err := ReadPocs(file)
		if err == nil {
			pocs = append(pocs, poc)
			fmt.Println("id: ", poc.Id)
		}
	}
	return pocs
}

// 根据keyword 过滤 返回pocs
func filterPoc(keyword string, pocpath string) []Poc {

	posc := loadTemplates(pocpath)
	fpoc := make([]Poc, 0)
	for _, poc := range posc {
		if strings.Contains(poc.Id, keyword) {
			fpoc = append(fpoc, poc)
			//println(poc.Id)
		}
	}
	return fpoc
}

func ShowPocList(keyword string) {
	if strings.ToLower(keyword) == "list" {
		//打印漏洞列表 todo
		return
	}
	//打印搜索的漏洞信息
	lists := strings.Split(keyword, ",")
	for _, list := range lists {
		list = strings.ToLower(list)
		pocs := filterPoc(list, "")
		for _, poc := range pocs {
			//输出需要美化
			println(poc.Info.Name + " " + poc.Info.Severity + " " + poc.Info.Description)
		}
	}
}

func ReadPocs(pocYaml string) (Poc, error) {
	var poc = Poc{}

	file, err := os.Open(pocYaml)
	if err != nil {
		return poc, err
	}
	defer file.Close()

	if err := yaml.NewDecoder(file).Decode(&poc); err != nil {
		return poc, err
	}
	return poc, nil
}

func FileForEachComplete(fileFullPath string) []string {
	files, err := ioutil.ReadDir(fileFullPath)
	if err != nil {
		log.Fatal(err)
	}
	var PocPath []string
	for _, file := range files {
		if file.IsDir() {
			path := strings.TrimSuffix(fileFullPath, "/") + "\\" + file.Name()
			subFile := FileForEachComplete(path)
			if len(subFile) > 0 {
				PocPath = append(PocPath, subFile...)
			}
		} else {
			if strings.HasSuffix(file.Name(), ".yaml") {
				PocPath = append(PocPath, fileFullPath+"\\"+file.Name())
				//println(fileFullPath + "\\" + file.Name())
			}
		}
	}
	return PocPath
}
