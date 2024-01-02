package getter

import (
	"io"
	"net/http"
	"proxypool/pkg/models"
	"regexp"
	"unknwon.dev/clog/v2"
)

func closeReaderIO(body io.ReadCloser) {
	err := body.Close()
	if err != nil {
		clog.Error("请求响应 close 错误： ", err)
	}
}

func IP89() (result []*models.IP) {
	clog.Info("开始爬取网站 89ip start")

	// 抓取的正则
	ExprIP := regexp.MustCompile(`((25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))\.){3}(25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))\:([0-9]+)`)

	// 爬取的网址
	pollURL := "http://www.89ip.cn/tqdl.html?api=1&num=100&port=&address=%E7%BE%8E%E5%9B%BD&isp="

	response, err := http.Get(pollURL)

	if err != nil {
		return nil
	}

	if response.StatusCode != 200 {
		return nil
	}

	defer closeReaderIO(response.Body)

	body, _ := io.ReadAll(response.Body)

	bodyHtml := string(body)

	ipList := ExprIP.FindAllString(bodyHtml, 100)

	clog.Info("ip list: %v", ipList)

	for _, ipString := range ipList {
		ip := models.NewIp()
		ip.Data = ipString
		ip.Type1 = "http"
		ip.Source = "89ip"
		clog.Info("[89ip] ip = %s", ip.Data)

		result = append(result, ip)
	}

	clog.Info("89 ip 爬取完成")

	return
}
