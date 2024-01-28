package getter

import (
	"github.com/Aiicy/htmlquery"
	"proxypool/pkg/models"
	"unknwon.dev/clog/v2"
)

func IP3306() (result []*models.IP) {
	clog.Info("[IP3306]  start GET ip proxy")

	pollURL := "http://www.ip3366.net/free/?stype=1&page=1"

	doc, _ := htmlquery.LoadURL(pollURL)
	trNode, err := htmlquery.Find(doc, "//div[@id='list']//table//tbody//tr")
	clog.Info("[IP3306] start up")

	if err != nil {
		clog.Info("[IP3306]] parse pollUrl error")
		clog.Warn(err.Error())
	}

	clog.Info("[IP3306] len(trNode) = %d ", len(trNode))

	for i := 1; i < len(trNode); i++ {
		tdNode, _ := htmlquery.Find(trNode[i], "//td")
		ip := htmlquery.InnerText(tdNode[0])
		port := htmlquery.InnerText(tdNode[1])
		Type := htmlquery.InnerText(tdNode[3])
		speed := htmlquery.InnerText(tdNode[5])

		IP := models.NewIp()
		IP.Data = ip + ":" + port

		if Type == "HTTPS" {
			IP.Type1 = "https"
			IP.Type2 = ""

		} else if Type == "HTTP" {
			IP.Type1 = "http"
		}
		IP.Source = "ip3366.net"
		IP.Speed = extractSpeed(speed)

		clog.Info("[IP3306] ip.Data = %s,ip.Type = %s,%s ip.Speed = %d", IP.Data, IP.Type1, IP.Type2, IP.Speed)

		result = append(result, IP)
	}

	clog.Info("IP 3306 done.")

	return
}
