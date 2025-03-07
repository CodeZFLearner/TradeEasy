/*
 * @Author: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @Date: 2025-01-05 18:44:24
 * @LastEditors: zff 2059577798@qq.com
 * @LastEditTime: 2025-03-07 20:07:08
 * @FilePath: \HelloGolang\main.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package main

import (
	"HelloGolang/internal/oss"
	"HelloGolang/internal/route"
	"HelloGolang/internal/service"
	"HelloGolang/pkg/common"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

func main() {
	//keys := [...]string{
	//	"pankou_2024-12-16_002131.sql.gz",
	//	"pankou_2024-12-17_002131.sql.gz",
	//	"pankou_2024-12-18_002131.sql.gz",
	//	"pankou_2024-12-19_002131.sql.gz",
	//	"pankou_2024-12-20_002131.sql.gz",
	//}
	//for _, key := range keys {
	//	wg.Add(1)
	//	go download(key)
	//}
	//
	//go func() {
	//	wg.Wait()
	//	close(fileCh)
	//}()
	//
	//for key := range fileCh {
	//	gzExtract(key)
	//	sql(key)
	//}

	http.HandleFunc("/crawl", route.Crawl)
	http.HandleFunc("/forward", route.Forward)
	http.HandleFunc("/test", route.Huobi)
	http.HandleFunc("/article", route.Article)
	http.HandleFunc("/generate/", route.Generate)

	// 将静态文件目录映射到URL路径
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	fs2 := http.FileServer(http.Dir("./"))
	http.Handle("/", http.StripPrefix("/", fs2))

	// 创建服务器
	server := &http.Server{
		Addr:              ":8080",
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1 MB
		IdleTimeout:       120 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}
	// start()

	// hanleHotNew()

	// 启动服务器
	log.Println("Starting server on :8080...")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %s", err)
	}

}

func start() {
	ticker := time.NewTicker(15 * 60 * time.Minute)
	// 创建一个通道，用于接收定时器的信号
	done := make(chan bool)

	// 启动一个 goroutine 来处理定时任务
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				// 执行定时任务
				// crawlAndSave()
				hotNewParse()
				fastNews()
				// 处理
				hanleHotNew()
				handleLastestNew()

			}
		}
	}()
}
func hanleHotNew() {
	fileHelp := oss.FileHelper{}
	var data map[string]interface{}
	var hotnewsDs map[string]interface{}

	if err := fileHelp.ReadFromFile("hotNew-ds.json", &hotnewsDs); err == nil {
		for key, value := range hotnewsDs {
			common.MyCache.Set(key, value, 10*time.Hour)
		}
	}

	if err := fileHelp.ReadFromFile("hotNew.json", &data); err != nil {
		fmt.Println("Error reading from file:", err)
		return
	}

	hotNew := make(map[string]interface{})
	// 处理 模型 输入
	q := "你是一名杂志主编，根据文章内容梳理素材，以markdown 格式返回 省略开头和结尾的概括性内容"

	for code, item := range data {
		v, f := common.MyCache.Get(code)
		if f {
			hotNew[code] = v
			continue
		}
		messages := []common.Message{
			{Role: "system", Content: "这是原文:" + item.(string)},
			{Role: "user", Content: q},
		}
		if answer, err := common.QDeepSeek(messages); err == nil {
			hotNew[code] = answer
			common.MyCache.Set(code, answer, 10*time.Hour)
		} else {
			fmt.Println(code + ":ds 失败")
		}
	}
	fileHelp.WriteToFile("hotNew-ds.json", hotNew)
}
func handleLastestNew() {
	fileHelp := oss.FileHelper{}
	var data map[string]interface{}
	err := fileHelp.ReadFromFile("fastNews-title.json", &data)
	if err != nil {
		fmt.Println("Error reading from file:", err)
		return
	}
	// 处理 模型 输入
	datasource := ""
	for key, item := range data {
		datasource = datasource + key + ":" + item.(string) + "\n"
	}

	q := "将数据源的新闻分为五个类别，使得五个类别可以覆盖到数据源中绝大部分新闻，使得尽可能多新闻符合分类并展示出来。每个分类能够灵活承载多样化的内容.每个类别下再划分三类：改善、恶化、不确定。一篇新闻只能属于一个类别且只能出现一次,并按照以下要求返回:" +
		"1.按照新闻和分类的相关度降序排序。省略开头和结尾概括性内容,每个分类下你只需要返回新闻的key,不需要文本标题 "
	messages := []common.Message{
		{Role: "system", Content: "这是源数据,数据的格式是 key:value，问题的答案都基于此:" + datasource},
		{Role: "user", Content: q},
	}
	// fmt.Println(datasource)
	// fmt.Println(q)
	if answer, err := common.QDeepSeek(messages); err == nil {
		fmt.Println(answer)
		// re := regexp.MustCompile(`\d{18}`)

		// // 查找所有匹配的数字串
		// matches := re.FindAllString(answer, -1)

		// // 遍历匹配结果并处理
		// for _, match := range matches {
		// 	answer = strings.ReplaceAll(answer, match, data[match].(string))
		// }
		fileHelp.WriteTxt("fastNews-answer.txt", answer)
	} else {
		fmt.Println(err)
	}
}

func hotNewParse() {
	fileHelp := oss.FileHelper{}
	c := url.Values{
		"category": {"HotNews"},
	}
	result, err := service.CrawlByCategory(c)
	jdata := make(map[string]interface{})

	if err == nil {
		var artcodes []string
		for _, item := range result["data"].([]map[string]interface{}) {
			if artcode, ok := item["artCode"].(string); ok {
				artcodes = append(artcodes, artcode)
				jdata[artcode] = item
			}
		}
		if result, err := service.EastmoneyArticleService(artcodes); err == nil {
			fileHelp.WriteToFile("hotNew.json", result)
		}
		fileHelp.WriteToFile("hotNew-abstract.json", jdata)
	}
}
func fastNews() {
	fileHelp := oss.FileHelper{}
	c := url.Values{
		"category": {"FastNews"},
		"pageSize": {"200"},
	}
	result, err := service.CrawlByCategory(c)
	if err != nil {
		log.Println(err)
		return
	}
	// digest := make(map[string]interface{})
	title := make(map[string]interface{})
	for _, item := range result["data"].([]map[string]interface{}) {
		artcode := item["code"].(string)
		// digest[artcode] = item["digest"].(string)
		title[artcode] = item["title"].(string)
	}
	// fileHelp.WriteToFile("fastNews-digest.json", digest)
	fileHelp.WriteToFile("fastNews-title.json", title)
}

func crawlAndSave() {
	category := []url.Values{
		{
			"category": {"Event"},
		},
		// {
		// 	"category": {"HotNews"},
		// },
		// {
		// 	"category": {"FastNews"},
		// 	"pageSize": {"10"},
		// },
	}
	fileHelp := oss.FileHelper{}
	jdata := make(map[string]interface{})
	for _, c := range category {
		categoryKey := c.Get("category")

		fmt.Println(c)
		result, err := service.CrawlByCategory(c)
		// if categoryKey == "HotNews" {
		// 	var artcodes []string
		// 	for _, item := range result["data"].([]map[string]interface{}) {
		// 		if artcode, ok := item["artCode"].(string); ok {
		// 			artcodes = append(artcodes, artcode)
		// 		}
		// 	}
		// 	if result, err := service.EastmoneyArticleService(artcodes); err == nil {
		// 		fileHelp.WriteToFile("artcodes.json", result)
		// 	}
		// }
		if err != nil {
			log.Println(err)
			return
		}

		jdata[categoryKey] = result
	}

	fileHelp.WriteToFile("article.json", jdata)
}
func now_time() string {
	now := time.Now()

	// month := now.Month() // 返回的是time.Month类型
	day := now.Day()
	hour := now.Hour()
	return fmt.Sprintf("%d-%d", day, hour)
}
