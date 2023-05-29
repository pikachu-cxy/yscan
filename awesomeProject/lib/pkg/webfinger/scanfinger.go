package webfinger

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"github.com/spaolacci/murmur3"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

func isKey(body string, key []string) bool {
	var x bool
	x = true
	for _, k := range key {
		if strings.Contains(strings.ToLower(body), strings.ToLower(k)) {
			x = x && true
		} else {
			x = x && false
		}
	}
	return x
}
func isReg(body string, key []string) bool {
	var x bool
	x = true
	for _, k := range key {
		re := regexp.MustCompile(k)
		if re.Match([]byte(body)) {
			x = x && true
		} else {
			x = x && false
		}
	}
	return x
}

func WebFinger(host string) []string {
	cmslist, err := LoadFingerPrint("finger.json")
	if err != nil {
		println("json jie xi chu cuo")
	}
	cmss := ScanFinger(host, cmslist)
	return cmss
}
func ScanFinger(url string, cmsList *PackFinger) []string {

	cmss := make([]string, 0)
	body, _ := HttpRequest(url)
	bodystring := string(body)
	favhash := getfavicon(strings.ToLower(bodystring), url)
	for _, cms := range cmsList.Fingerprint {
		if cms.Location == "body" {
			if cms.Method == "keyword" {
				if isKey(bodystring, cms.Keyword) {
					cmss = append(cmss, cms.CMS)
				}
			}
			if cms.Method == "faviconhash" {
				if favhash == cms.Keyword[0] {
					cmss = append(cmss, cms.CMS)
				}
			}
			if cms.Method == "regular" {
				if isReg(bodystring, cms.Keyword) {
					cmss = append(cmss, cms.CMS)
				}
			}
		}
		if cms.Location == "header" {
			if cms.Method == "keyword" {
				if isKey(bodystring, cms.Keyword) {
					cmss = append(cmss, cms.CMS)
				}
			}
		}
	}
	result := RemoveDuplicates(cmss)
	return result
}

func xegexpjs(reg string, resp string) (reslut1 [][]string) {
	reg1 := regexp.MustCompile(reg)
	if reg1 == nil {
		log.Println("regexp err")
		return nil
	}
	result1 := reg1.FindAllStringSubmatch(resp, -1)
	return result1
}

func getfavicon(httpbody string, turl string) string {
	faviconpaths := xegexpjs(`href="(.*?favicon....)"`, httpbody)
	var faviconpath string
	u, err := url.Parse("http://" + turl)
	if err != nil {
		panic(err)
	}
	turl = u.Scheme + "://" + u.Host
	if len(faviconpaths) > 0 {
		fav := faviconpaths[0][1]
		if fav[:2] == "//" {
			faviconpath = "http:" + fav
		} else {
			if fav[:4] == "http" {
				faviconpath = fav
			} else {
				faviconpath = turl + "/" + fav
			}

		}
	} else {
		faviconpath = turl + "/favicon.ico"
	}
	return favicohash(faviconpath)
}
func favicohash(host string) string {
	timeout := time.Duration(8 * time.Second)
	var tr *http.Transport

	tr = &http.Transport{
		MaxIdleConnsPerHost: -1,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		DisableKeepAlives:   true,
	}

	client := http.Client{
		Timeout:   timeout,
		Transport: tr,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse /* 不进入重定向 */
		},
	}
	resp, err := client.Get(host)
	if err != nil {
		//log.Println("favicon client error:", err)
		return "0"
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			//log.Println("favicon file read error: ", err)
			return "0"
		}
		faviconMMH3 := fmt.Sprintf("%d", FaviconHash(body))
		return faviconMMH3
	} else {
		return "0"
	}
}
func FaviconHash(data []byte) int32 {
	stdBase64 := base64.StdEncoding.EncodeToString(data)
	stdBase64 = InsertInto(stdBase64, 76, '\n')
	hasher := murmur3.New32WithSeed(0)
	hasher.Write([]byte(stdBase64))
	return int32(hasher.Sum32())
}
func InsertInto(s string, interval int, sep rune) string {
	var buffer bytes.Buffer
	before := interval - 1
	last := len(s) - 1
	for i, char := range s {
		buffer.WriteRune(char)
		if i%interval == before && i != last {
			buffer.WriteRune(sep)
		}
	}
	buffer.WriteRune(sep)
	return buffer.String()
}

func RemoveDuplicates(arr []string) []string {
	// 创建一个空的 map 用于存储元素和出现的次数
	seen := make(map[string]bool)
	result := []string{} // 存储去重后的结果

	for _, element := range arr {
		if !seen[element] {
			// 如果元素在 map 中不存在，则添加到结果列表中
			result = append(result, element)
			seen[element] = true
		}
	}

	return result
}
