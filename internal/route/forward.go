package route

import (
	"HelloGolang/internal/platform"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func Forward(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	url := query.Get("url")
	method := query.Get("method")
	//param := query.Get("param")
	body := bytes.NewBuffer(nil)

	if url == "" {
		http.Error(w, "url 空", 500)
		return
	}
	if method == "" {
		method = "GET"
	}
	fmt.Println(url)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "服务器错误", 500)
		return
	}
	if strings.Contains(url, "xueqiu") {
		output := platform.XqId()
		fmt.Println(output)
		parts := strings.Split(output, ";")
		req.AddCookie(&http.Cookie{
			Name:  "u",
			Value: parts[0],
		})
		req.AddCookie(&http.Cookie{
			Name:  "xqat",
			Value: parts[1],
			//Path:     "/",
			//Domain:   ".xueqiu.com",
			//HttpOnly: true,
			//Secure:   false,
			//SameSite: http.SameSiteLaxMode,
		})
		// Print the output
		cookie, err := req.Cookie("xqat")
		if err != nil {
			fmt.Println("Cookie not found:", err)
		} else {
			fmt.Println("Cookie value:", cookie)
		}
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "服务器错误", 500)
		return
	}
	defer resp.Body.Close()
	// 将响应头复制到响应中
	for k, v := range resp.Header {
		w.Header()[k] = v
	}
	// 设置响应状态码
	w.WriteHeader(resp.StatusCode)

	// 将响应体复制到响应中
	io.Copy(w, resp.Body)
}
