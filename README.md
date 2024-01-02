# golang 实现自己的 IP 代理池

## 爬虫与 IP 代理池不得不说的故事

说到爬虫， 那么 IP 代理池是一个绕不开的话题。 现在很多网站为了防止别人爬取自己的站点， 采用的手段很粗暴， 但凡识别到有爬虫迹象的 IP 直接封禁该 IP 的访问即可。
那问题就来了， 服务器的 IP 地址是静态的， 就算是家用电脑，IP 地址大多数情况下也是固定静态的。 一旦被某网站封禁，就意味着你的爬虫废了。 

这个是爬虫的防御手段问题， 那么涉及到爬虫的成本问题， 很多人会直接通过 HTTP 请求去获取人家渲染好的 HTML 模板爬取内容， 甚至直接注入 cookie 之后， 用服务器发送请求， 去获取别人站点的数据。
很负责人的告诉大家， 这样被封禁 IP 的非常高， 我在 18 年学习 python 爬虫框架 Scrapy 的时候经常遇到 IP 被封禁访问的情况。 （[附加链接 18 年学习的链接](https://github.com/yanlele/python-index/tree/master/book/02%E3%80%81Scrapy%E5%88%86%E5%B8%83%E5%BC%8F%E7%88%AC%E8%99%AB%E9%A1%B9%E7%9B%AE%E5%AD%A6%E4%B9%A0)）

为了让自己的爬虫没有那么容易被封禁， 爬虫选手们就需要尽力让自己的爬虫伪装成一个真正的用户行为。这个时候就要祭出各种的无头浏览器技术，例如：`Puppeteer、PhantomJS` 等；

尽管你有无头浏览器技术， 尽管你伪装成了用户行为， 但是在站点行为分析和访问量分析的话， 还是有一定概率被定盯上。

那么如何从底层技术上突破爬虫封锁呢？**就是 IP 代理池**。

IP 代理池是啥？就是你的爬虫访问 A 网站的时候， 不是你自己的 IP ， 是一个别人的 IP 在帮你爬取数据，使用别人的 IP 代理你去爬取数据， 然后爬取到了数据之后再给你。 
如果 IP 被封锁， 那么也是别人的 IP 被封锁， 跟你没有任何关系。 这样可以代理的 IP 成千上万的话， 就形成了 IP 代理池。

所以对于通常有 IP 代理池的爬虫，在爬取数据的时候， 不再需要担心自己本身的 IP 被封禁了。 

## 为何要用 golang 来实现？

golang 的优势我就不多说了。 
直接说结果吧，这个代理池程序， 最后可以打包为一个 20MB 的可执行文件， 丢谁都可以直接运行， 没有任何依赖， 内存消耗很低， 1核1G的机器都可以跑的飞快。

## 程序的开始前言

**先说说我们实现 IP 代理池的原理**。 

实际上网上有很多付费的 IP 代理池， 稳定且高效， 高匿不说， IP 量还异常的大。 但是架不住费用非常高呀，到底费用有多高， 各位小伙伴们， 网上搜索一下就知道了。 
所以这种高费用的商业 IP 代理池， 并不适合于大家用于学习。所以我们要实现一下自己的 IP 代理池， 简单学习爬虫就够用了。 

**那么我们是如何实现的呢**？
其实这个也非常简单， 有一些网上有一些IP代理网站， 为了吸引别人来使用， 会放出一些免费的 IP 代理地址， 但是这些免费的 IP 代理地址变化非常快， 一般一两个小时就刷新一次， 而且速度较慢， 稳定性很低， 甚至有很多压根就没法用。 
我们就怕这些 IP 获取下来， 进行测速、筛选、分类， 做成我们的临时 IP 代理池， 对于学习爬虫的同学们来说， 其实也足够用了。

**接下来我会一步一步拆解实现一个 IP 代理池**
当然具体的实现不是我原创的， github 有一位大神， 已经实现了一套 ip 代理池， 地址可以参考： https://github.com/henson/proxypool

但是大神使用的 IP 代理池是需要写入数据库的， 配置还挺复杂的。 本身这种实现方式的 IP 代理池， IP 变化都非常快， 其实没有必要存数据库；所以我们要对他的实现进行简化魔改, 就是零配置跑起来。

## 兵马未动日志先行

任何程序第一步要搭建自己的良好的日志环境。 没有日志， 就没有任何线上排查运维可言， 就是蒙着眼睛的裸奔瞎子。

首先来看看目录：
```
├── README.md
├── go.mod
├── go.sum
├── logs
├── pkg
│   ├── getter
│   └── models
└── src
    └── main.go
```

项目初始化就不做过多赘述， 日志系统使用 `clog`
```
go get unknwon.dev/clog/v2
```

项目启动的时候， 直接使用初始化日志即可：
```go
package main

import (
	"proxypool/pkg/getter"
	log "unknwon.dev/clog/v2"
)

func init() {
	err := log.NewConsole()
	if err != nil {
		panic("unable to create new logger: " + err.Error())
	}

	err = log.NewFile(
		log.FileConfig{
			Filename: "./logs/clog.log",
			Level:    log.LevelTrace,
			FileRotationConfig: log.FileRotationConfig{
				Rotate: true,
				Daily:  true,
				//MaxLines: 50,
			},
		},
	)

	if err != nil {
		panic("unable to create new logger with file: " + err.Error())
	}
}

func main() {
	// ......
	defer log.Stop()
}
```
通过 `log.NewConsole() 和 log.NewFile()` 初始化日志在 terminal 输出和本地文件输出。
最后程序停止的时候停止日志 `log.Stop()` 调用即可；

## 获取 IP 

本质上来说， 本地去访问公开的 IP 代理池网站， 获取别人的 HTML 模板， 从里面提取 IP 信息

### 定义 IP 结构体
首先先统一一下我们 IP 的结构体

对应目录：
```
│   └── clog.log
├── pkg
│   └── models
│       ├── ip.go
│       └── type.go
└── src
    └── main.go
```


`pkg/modules/type.go`
```go
package models

import "time"

type IP struct {
	ID         int64     `json:"id"`
	Data       string    `json:"data"`
	Type1      string    `json:"type1"`
	Type2      string    `json:"type2"`
	Speed      int64     `json:"speed"`  // 链接速度
	Source     string    `json:"source"` // 代理来源
	CreateTime time.Time `json:"create_time"`
}
```

给一个实例化一个 IP 的函数 `pkg/models/ip.go`：
```go
package models

import "time"

func NewIp() *IP {
	return &IP{
		Speed:      -1,
		CreateTime: time.Now(),
	}
}
```

### 直接获取 HTML 文本， 正则匹配提取 IP 信息

第一个示范以 https://www.89ip.cn/ti.html 网址为例， 该网址直接将 IP 晒到了 html body 节点里面， 甚至都不需要做 HTML dom 解析。
![image.png](https://p1-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/4599bc1484df45d7832f678e64e4dd25~tplv-k3u1fbpfcp-jj-mark:0:0:0:0:q75.image#?w=2552&h=1336&s=488931&e=png&b=ffffff)

所以提取他们也比较简单， 思路就是直接获取 HTML 之后暴力正则匹配。 没有啥好说的， 直接上代码

```go
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
```

这里需要注意的是 `http.Get` 函数返回的响应体是一个`io.ReadCloser`类型的对象，该对象实现了`io.Reader`和`io.Closer`接口。`io.Reader`用于读取响应体的内容，而`io.Closer`用于关闭响应体，释放相关资源。

根据官方文档的描述，使用`http.Get`得到的响应体需要手动关闭，以确保及时释放资源。不关闭响应体可能导致资源泄露，尤其是在连续发送多个HTTP请求时，如果不关闭响应体，可能会导致连接池用尽。

在 main.go 函数中测试一下：
`src/mian.go`
```go
// ......
func main() {
	getter.IP89()
	defer log.Stop()
}
```
执行结果如下：
![image.png](https://p6-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/5fb12019741b45cabdd72600881b0d6c~tplv-k3u1fbpfcp-jj-mark:0:0:0:0:q75.image#?w=2402&h=1590&s=439722&e=png&b=2b2b2b)

### 通过 xpath 解析复杂 dom 节点

