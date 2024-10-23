package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/parnurzeal/gorequest"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: false})
	logrus.SetLevel(logrus.DebugLevel)
}

func main() {
	result := map[string]interface{}{}
	response, data, errors := gorequest.New().TLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
		Get("https://blog.ysicing.net/feed/json").Timeout(15*time.Second).
		Set("User-Agent", "Profile Bot; +https://github.com/ysicing/ysicing").EndStruct(&result)
	if errors != nil || response.StatusCode > http.StatusBadRequest {
		logrus.Fatalf("request failed, errors: %v, response: %v, data: %v", errors, response, data)
	}
	buf := &bytes.Buffer{}
	buf.WriteString("\n\n")
	cstSh, _ := time.LoadLocation("Asia/Shanghai")
	updated := time.Now().In(cstSh).Format("2006-01-02 15:04:05")
	buf.WriteString("### 我在[Solitudes](https://ysicing.me)的近期动态\n\n每天自动刷新，最近更新时间：`" + updated + "`\n\n")
	for id, event := range result["items"].([]interface{}) {
		if id > 6 {
			break
		}
		evt := event.(map[string]interface{})
		url := evt["url"].(string)
		title := evt["title"].(string)
		buf.WriteString("* " + " [" + title + "](" + url + ")\n")
	}
	buf.WriteString("*  [随便看看, 随机推荐](https://ysicing.me/random/)\n")
	buf.WriteString("\n\n")
	fmt.Println(buf.String())
	readme, err := os.ReadFile("README.md")
	if nil != err {
		logrus.Fatalf("read README.md failed: %s", data)
	}
	startFlag := []byte("<!--events start -->")
	beforeStart := readme[:bytes.Index(readme, startFlag)+len(startFlag)]
	newBeforeStart := make([]byte, len(beforeStart))
	copy(newBeforeStart, beforeStart)
	endFlag := []byte("<!--events end -->")
	afterEnd := readme[bytes.Index(readme, endFlag):]
	newAfterEnd := make([]byte, len(afterEnd))
	copy(newAfterEnd, afterEnd)
	newReadme := append(newBeforeStart, buf.Bytes()...)
	newReadme = append(newReadme, newAfterEnd...)
	if err := os.WriteFile("README.md", newReadme, 0644); nil != err {
		logrus.Fatalf("write README.md failed: %s", data)
	}
}
