package crawler

import (
	"TestProject/config"
	"TestProject/qyapi"
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"github.com/gocolly/colly/queue"
	"math/big"
	"net/http"
	"strconv"
	"time"
)

type Fetcher struct {
	Fetcher *colly.Collector
}

type xueQue struct {
	Count int
	Statuses []map[string]interface{}
}

/*
请求执行之前调用
	- OnRequest
响应返回之后调用
	- OnResponse
监听执行 selector
	- OnHTML
监听执行 selector
	- OnXML
错误回调
	- OnError
完成抓取后执行，完成所有工作后执行
	- OnScraped
取消监听，参数为 selector 字符串
	- OnHTMLDetach
取消监听，参数为 selector 字符串
	- OnXMLDetach
*/

func (c *Fetcher) Init() *colly.Collector  {
	c.Fetcher = colly.NewCollector()
	// 实例化默认收集器

	c.Fetcher.AllowURLRevisit = true
	// 允许重复访问

	extensions.Referer(c.Fetcher)
	extensions.RandomUserAgent(c.Fetcher)
	//随机UserAgent

	// 限制采集规则
	/*
		在Colly里面非常方便控制并发度，只抓取符合某个(些)规则的URLS
	colly.LimitRule{DomainGlob: "*.douban.*", Parallelism: 5}，表示限制只抓取域名是douban(域名后缀和二级域名不限制)的地址，当然还支持正则匹配某些符合的 URLS

	Limit方法中也限制了并发是5。为什么要控制并发度呢？因为抓取的瓶颈往往来自对方网站的抓取频率的限制，如果在一段时间内达到某个抓取频率很容易被封，所以我们要控制抓取的频率。
	另外为了不给对方网站带来额外的压力和资源消耗，也应该控制你的抓取机制。
	*/

	err := c.Fetcher.Limit(&colly.LimitRule{
		DomainGlob: "*httpbin.*",
		RandomDelay: 2 * time.Second,
		//设置对域请求之间的延迟
		Parallelism: 5,
		//设置并发
	})

	if err != nil {
		fmt.Println(err)
	}
	return c.Fetcher
}

func (c *Fetcher) GetCookie()  ([]*http.Cookie, error) {
	var coo []*http.Cookie
	c.Fetcher.OnResponse(func(r *colly.Response) {
		coo = c.Fetcher.Cookies(r.Request.URL.String())

	})

	err := c.Fetcher.Visit(config.SiteXq)
	if err != nil {
		fmt.Println(err)
		return coo, err
	}

	return coo, nil
}

func (c *Fetcher) GetXq(coo []*http.Cookie) (string, error)  {
	var str string
	var e error

	q, _ := queue.New(
		2,
		&queue.InMemoryQueueStorage{MaxSize: 1000},
		)

	c.Fetcher.OnResponse(func(r *colly.Response) {
		var s xueQue
		json.Unmarshal(r.Body,&s)
		str = s.Statuses[0]["text"].(string)

		if isTi(s.Statuses[0]["created_at"].(float64)) {
			qyapi.SendCardMsg(str)
		}
		fmt.Println(s.Statuses[0]["created_at"].(float64))
	})

	c.Fetcher.OnError(func(r *colly.Response, err error) {
		fmt.Println(err)
		e = err
	})

	c.Fetcher.OnRequest(func(r *colly.Request) {
		fmt.Println(r.URL)
	})

	for _, s := range config.SiteUrl {
		c.Fetcher.SetCookies(s, coo)
		//c.Fetcher.Visit(s)
		q.AddURL(s)
	}
	//fmt.Println()
	q.Run(c.Fetcher)

	return str, e
}



func isTi(n float64) (b bool) {
	ns := float64(n)
	s := big.NewRat(1,1)
	s.SetFloat64(ns)
	ss, _ := strconv.ParseInt(s.FloatString(0)[:10], 10, 64)

	t := time.Now()
	tt := t.Unix()

	if (tt - ss) <= 30 {
		return true
	}

	return false
}