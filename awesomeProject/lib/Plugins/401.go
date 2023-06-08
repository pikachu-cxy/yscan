package Plugins

import (
	"fmt"
	"net/http"
	"time"
)

func Force401Scan(model PluginService) {
	var basicusers = []string{"admin", "root", "test", "tomcat"}
	var basicpass = []string{"123456", "root", "admin", "tomcat", "test", "11111", "111111"}

	client := &http.Client{Timeout: time.Second}
	if req, err := client.Head(model.Host); err == nil {
		if req.StatusCode == 401 {
			req, err := http.NewRequest("HEAD", model.Host, nil)
			if err == nil {
				for _, user := range basicusers {
					for _, pass := range basicpass {
						req.SetBasicAuth(user, pass)
						if resp, err := client.Do(req); err == nil {
							if resp.StatusCode == 200 {
								println(fmt.Sprintf("[+]401 success! %s:%s", user, pass))
								return
							}
						}
					}
				}
			}
		}
	}
}
