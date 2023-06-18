package Plugins

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var user_agent = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.129 Safari/537.36,Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/27.0.1453.93 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.129 Safari/537.36,Mozilla/5.0 (Windows NT 6.2; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/30.0.1599.17 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.129 Safari/537.36,Mozilla/5.0 (X11; NetBSD) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/27.0.1453.116 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.129 Safari/537.36,Mozilla/5.0 (X11; CrOS i686 3912.101.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/27.0.1453.116 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.129 Safari/537.36,Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/27.0.1453.93 Safari/537.36",
}

func PathCrack(model PluginService) {
	var wg sync.WaitGroup
	requests := make(chan string)
	result := make([]string, 0)
	cl := make([]int64, 0)
	var black = []string{"not found", "未找到", "界面丢失了"}
	var statuscode = []string{"301", "200", "201", "403", "401", "302", "500"}
	path := "path.txt"
	threads := 100
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr, Timeout: 5 * time.Second}
	lines, _ := ReadLinesFromFile(path)
	//bar := progressbar.Default(int64(len(lines)), "pathBrute")
	println("正在进行目录扫描~ 请稍等---------------------------------")
	//请求原始页面
	//if !strings.HasPrefix(model.Host,"http"){}
	resp, err := client.Get(model.Host)
	if err != nil {
		println("请求网址出错~ 请检查网络状况无误后重试")
		return
	}
	if resp.StatusCode == 400 {
		uri, f := strings.CutPrefix(model.Host, "http://")
		if f {
			model.Host = fmt.Sprintf("https://%s", uri)
		}
	}
	//println("正在进行目录扫描~ 请稍等---------------------------------")
	for i := 0; i < threads; i++ {
		//time.Sleep(40 * time.Millisecond)
		go func() {
			for line := range requests {
				//随机请求头
				re, err := http.NewRequest("GET", line, nil)
				re.Header.Add("user-agent", user_agent[rand.Intn(5)])
				if err == nil {
					if req, err := client.Do(re); err == nil {
						for _, s := range statuscode {
							if strings.Contains(strconv.Itoa(req.StatusCode), s) {
								body, _ := ioutil.ReadAll(req.Body)
								for _, b := range black {
									if !strings.Contains(string(body), b) && req.ContentLength > 10 {
										if len(cl) == 0 {
											if resp.ContentLength != req.ContentLength {
												cl = append(cl, req.ContentLength)
												result = append(result, fmt.Sprintf("扫描路径为： %s,[%d] [%d]", line, req.ContentLength, req.StatusCode))
												//println(fmt.Sprintf("路径为： %s,[%d] [%d]", line, req.ContentLength, req.StatusCode))
											}
										} else {
											for _, c := range cl {
												if req.ContentLength != c {
													cl = append(cl, req.ContentLength)
													result = append(result, fmt.Sprintf("扫描路径为： %s,[%d] [%d]", line, req.ContentLength, req.StatusCode))
												}
											}
										}

									}
								}
							}
						}
					}
				}
				wg.Done()
			}
		}()
	}

	// 将目录添加到请求通道中
	for _, dir := range lines {
		wg.Add(1)
		path := model.Host + dir
		requests <- path
	}

	close(requests)
	wg.Wait()

	result1 := removeDuplicates(result)
	for _, r := range result1 {
		println(r)
	}
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
