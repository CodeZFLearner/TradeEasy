/*
 * @Author: zff 2059577798@qq.com
 * @Date: 2025-03-03 18:07:32
 * @LastEditors: zff 2059577798@qq.com
 * @LastEditTime: 2025-03-07 22:54:25
 * @FilePath: \HelloGolang\internal\route\generate.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package route

import (
	"HelloGolang/internal/oss"
	"HelloGolang/internal/service"
	"HelloGolang/pkg/common"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

func Generate(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	fmt.Println(path)
	switch path {
	case "/generate/g1":
		g1(w, r)
	case "/generate/g2":
		g1Spider(w, r)
	// case "/g3":
	// 	productHandler(w, r)
	default:
		http.NotFound(w, r)
	}

}

func g1Spider(w http.ResponseWriter, r *http.Request) {
	type RequestData struct {
		Code string   `json:"code"`
		Orgs []string `json:"orgs"`
	}

	// 解析 JSON 请求体
	var requestData RequestData
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	res, err := service.SpiderOrg(requestData.Orgs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("res complete....", res)
	go func() {
		fileName := fmt.Sprintf("evidence-%s.json", requestData.Code)
		fmt.Println(00000)
		if _, err := os.Stat(fileName); err != nil {
			fmt.Println(err)
			fmt.Println(fileName + "文件不存在,starting....")
			evidence, err := service.GenerateEvidence(res)
			if err != nil {
				return
			}
			filehelp := oss.FileHelper{}
			filehelp.WriteToFile(fileName, evidence)
		} else {
			fmt.Println(222222)
		}
		fmt.Println(3333333)

	}()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(""); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}
func g1(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	code := query.Get("code")
	fmt.Println("g111111")
	// 文件
	res := service.GeneratePaper(code)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}
func article(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	code := query.Get("code")

	generateFileName := fmt.Sprintf("generate-%s.txt", code)
	if _, err := os.Stat(generateFileName); os.IsNotExist(err) {
		// 未生成过
		if _, f := common.MyCache.Get(code); !f {
			common.MyCache.Set(code, "1", 10*time.Minute)
			service.GenereateArticle(code)
		}
	} else {
		fmt.Println(err)
		fmt.Println(os.IsNotExist(err))
	}
	fmt.Println(111111)
	fileHelp := oss.FileHelper{}
	res, err := fileHelp.ReadTxt(generateFileName)
	if err != nil {
		fmt.Println(22222222)
		fmt.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}
