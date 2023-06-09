package Plugins

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"
)

/*
jsfinder:
"([a-zA-Z0-9_\\-]{1,}\\.(?:php|asp|aspx|jsp|json|action|html|js|txt|xml)(?:\\?[^\"|']{0,}|))",
"([a-zA-Z0-9_\\-/]{1,}/[a-zA-Z0-9_\\-/]{1,}\\.(?:[a-zA-Z]{1,4}|action)(?:[\\?|/][^\"|']{0,}|))",
"((?:/|\\.\\./|\\./)[^\"'><,;| *()(%%$^/\\\\\\[\\]][^\"'><,;|()]{1,})",
"((?:[a-zA-Z]{1,10}://|//)[^\"'/]{1,}\\.[a-zA-Z]{2,}[^\"']{0,})",
*/
var (
	leak_info_patterns = map[string][]string{

		"UrlFind": []string{`['\"](([a-zA-Z0-9]+:)?\\/\\/)?[a-zA-Z0-9\\-\\.]*?\\.(xin|com|cn|net|com.cn|vip|top|cc|shop|club|wang|xyz|luxe|site|news|pub|fun|online|win|red|loan|ren|mom|net.cn|org|link|biz|bid|help|tech|date|mobi|so|me|tv|co|vc|pw|video|party|pics|website|store|ltd|ink|trade|live|wiki|space|gift|lol|work|band|info|click|photo|market|tel|social|press|game|kim|org.cn|games|pro|men|love|studio|rocks|asia|group|science|design|software|engineer|lawyer|fit|beer|我爱你|中国|公司|网络|在线|网址|网店|集团|中文网)(\\:\\d{1,5})?(\\/.*?)?['\"]`,
			"(https{0,1}:[-a-zA-Z0-9（）@:%_\\+.~#?&//=]{2,250}?[-a-zA-Z0-9（）@:%_\\+.~#?&//=]{3}[.]js)",
			"[\"'‘“`]\\s{0,6}(/{0,1}[-a-zA-Z0-9（）@:%_\\+.~#?&//=]{2,250}?[-a-zA-Z0-9（）@:%_\\+.~#?&//=]{3}[.]js)",
			"=\\s{0,6}[\",',’,”]{0,1}\\s{0,6}(/{0,1}[-a-zA-Z0-9（）@:%_\\+.~#?&//=]{2,250}?[-a-zA-Z0-9（）@:%_\\+.~#?&//=]{3}[.]js)",
			"[\"'‘“`]\\s{0,6}(https{0,1}:[-a-zA-Z0-9()@:%_\\+.~#?&//={}]{2,250}?)\\s{0,6}[\"'‘“`]",
			"=\\s{0,6}(https{0,1}:[-a-zA-Z0-9()@:%_\\+.~#?&//={}]{2,250})",
			"[\"'‘“`]\\s{0,6}([#,.]{0,2}/[-a-zA-Z0-9()@:%_\\+.~#?&//={}]{2,250}?)\\s{0,6}[\"'‘“`]",
			"\"([-a-zA-Z0-9()@:%_\\+.~#?&//={}]+?[/]{1}[-a-zA-Z0-9()@:%_\\+.~#?&//={}]+?)\"",
			"href\\s{0,6}=\\s{0,6}[\"'‘“`]{0,1}\\s{0,6}([-a-zA-Z0-9()@:%_\\+.~#?&//={}]{2,250})|action\\s{0,6}=\\s{0,6}[\"'‘“`]{0,1}\\s{0,6}([-a-zA-Z0-9()@:%_\\+.~#?&//={}]{2,250})",
		},
		"domain":               []string{`['"](([a-zA-Z0-9]+:)?\/\/)?[a-zA-Z0-9\-\.]*?\.(xin|com|cn|net|com.cn|vip|top|cc|shop|club|wang|xyz|luxe|site|news|pub|fun|online|win|red|loan|ren|mom|net.cn|org|link|biz|bid|help|tech|date|mobi|so|me|tv|co|vc|pw|video|party|pics|website|store|ltd|ink|trade|live|wiki|space|gift|lol|work|band|info|click|photo|market|tel|social|press|game|kim|org.cn|games|pro|men|love|studio|rocks|asia|group|science|design|software|engineer|lawyer|fit|beer|我爱你|中国|公司|网络|在线|网址|网店|集团|中文网)(\:\d{1,5})?(\/)?['"]`},
		"IP":                   []string{"([^0-9]((127\\.0\\.0\\.1)|(10\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3})|(172\\.((1[6-9])|(2\\d)|(3[01]))\\.\\d{1,3}\\.\\d{1,3})|(192\\.168\\.\\d{1,3}\\.\\d{1,3})))", `['"](([a-zA-Z0-9]+:)?\/\/)?\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\:\d{1,5}(\/.*?)?['"]`},
		"IDCard":               []string{`[^0-9]((\d{8}(0\d|10|11|12)([0-2]\d|30|31)\d{3}$)|(\d{6}(18|19|20)\d{2}(0[1-9]|10|11|12)([0-2]\d|30|31)\d{3}(\d|X|x)))[^0-9]`},
		"phone":                []string{`[^0-9A-Za-z](1(3([0-35-9]\d|4[1-8])|4[14-9]\d|5([\d]\d|7[1-79])|66\d|7[2-35-8]\d|8\d{2}|9[89]\d)\d{7})[^0-9A-Za-z]`},
		"SpringBoot":           []string{`((local.server.port)|(:{\"mappings\":{\")|({\"_links\":{\"self\":))`},
		"swagger":              []string{"((swagger-ui.html)|(\"swagger\":)|(Swagger UI)|(swaggerUi))"},
		"druid":                []string{"((Druid Stat Index)|(druid monitor))"},
		"oss":                  []string{"([A|a]ccess[K|k]ey[I|i][d|D]|[A|a]ccess[K|k]ey[S|s]ecret)"},
		"jdbc-connect":         []string{`(jdbc:[a-z:]+://[A-Za-z0-9\.\-_:;=/@?,&]+)`},
		"Github Access Token":  []string{`([a-z0-9_-]*:[a-z0-9_\-]+@github\.com*)`},
		"Authorization Header": []string{`((basic [a-z0-9=:_\+\/-]{5,100})|(bearer [a-z0-9_.=:_\+\/-]{5,100}))`},
		"key":                  []string{"(session_key|sessionKey|secret|access_token)"},
		"Bucket":               []string{"(InvalidBucketName|NoSuchBucket|<Key>)"},
		"path":                 []string{"/['\"](?:\\/|\\.\\.\\/|\\.\\/)[^\\/\\>\\< \\)\\(\\{\\}\\,\\'\\\"\\\\]([^\\>\\< \\)\\(\\{\\}\\,\\'\\\"\\\\])*?['\"]", `['"][^\/\>\< \)\(\{\}\,\'\"\\][\w\/]*?\/[\w\/]*?['"]`},
	}

	filter = []string{
		"www\\.w3\\.org",
		"example\\.com",
		".*css.?$",
		".*png.?$",
		".*jpg.?$",
		".*jpeg.?$",
		".*gif.?$",
		".*svg.?$",
		".*vue.?$",
		".*scss.?$",
		//"\\.css\\?|\\.jpeg\\?|\\.jpg\\?|\\.png\\?|.gif\\?|www\\.w3\\.org|example\\.com|\\<|\\>|\\{|\\}|\\[|\\]|\\||\\^|;|/js/|\\.src|\\.replace|\\.url|\\.att|\\.href|location\\.href|javascript:|location:|application/x-www-form-urlencoded|\\.createObject|:location|\\.path|\\*#__PURE__\\*|\\*\\$0\\*|\\n",
		//".*\\.css$|.*\\.scss$|.*,$|.*\\.jpeg$|.*\\.jpg$|.*\\.png$|.*\\.gif$|.*\\.ico$|.*\\.svg$|.*\\.vue$|.*\\.ts$",
	}
)

func filterfunc(info string) bool {
	for _, b := range filter {
		match, _ := regexp.MatchString(b, info)
		if match == true {
			return false
		}
	}
	return true
}

func infoFinder(body string, url string) {
	info := make([]string, 0)
	var match bool
	for _, leak := range leak_info_patterns {
		for _, k := range leak {
			infos := regexp.MustCompile(k).FindAllStringSubmatch(body, -1)
			for _, inf := range infos {
				//sort 去重 排除静态资源/黑名单网址
				match = filterfunc(inf[0])
				if match == true {
					if !strings.Contains(inf[0], "\n") {
						info = append(info, inf[0])
					}
				}

			}
		}
	}

	result := removeDuplicates(info)
	println(fmt.Sprintf("jsfinder爬取到 %s 的信息如下：", url))
	for _, r := range result {
		println(r)
	}

	return
}

func removeDuplicates(arr []string) []string {
	// 创建一个 map 用于记录已经出现的字符串
	seen := make(map[string]bool)

	// 创建一个结果数组
	result := []string{}

	for _, str := range arr {
		// 如果字符串未出现过，则将其添加到结果数组中，并将其记录到 map 中
		if !seen[str] {
			result = append(result, str)
			seen[str] = true
		}
	}
	return result
}

func JsFinder(model PluginService) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr, Timeout: time.Second}
	//req, err := http.NewRequest("GET", model.Host, nil)

	resp, err := client.Get(model.Host)
	if err == nil {
		body, _ := ioutil.ReadAll(resp.Body)
		strbody := string(body)
		infoFinder(strbody, model.Host)
		resp.Body.Close()
	} else {
		println(err.Error())
	}

}
