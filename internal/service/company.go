package service

import (
	"HelloGolang/internal/oss"
	"HelloGolang/pkg/common"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

// Define the struct to match the JSON response
type NewsResponse struct {
	ReqTrace string   `json:"req_trace"`
	Code     string   `json:"code"`
	Message  string   `json:"message"`
	Data     NewsData `json:"data"`
}

type NewsData struct {
	TotalHits int        `json:"totalHits"`
	SortEnd   string     `json:"sortEnd"`
	PageSize  int        `json:"pageSize"`
	News      []NewsItem `json:"news"`
}

type NewsItem struct {
	Code      string      `json:"code"`
	Title     string      `json:"title"`
	NpDst     string      `json:"np_dst"`
	ShowTime  string      `json:"showTime"`
	Interact  interface{} `json:"interact"`
	MediaName string      `json:"mediaName"`
	Content   string
}

// Define the struct to match the JSON response for the new data
type OrgList struct {
	List []GubaResponse `json:"list"`
}
type EvidenceItem struct {
	Key  string `json:"key"`
	Data string `json:"data"`
}
type EvidenceData struct {
	List []EvidenceItem `json:"list"`
}
type OrgResponse struct {
	Re []struct {
		// Detail   GubaResponse
		PostID   int64 `json:"post_id"`
		PostUser struct {
			UserID               string      `json:"user_id"`
			UserNickname         string      `json:"user_nickname"`
			UserName             string      `json:"user_name"`
			UserV                int         `json:"user_v"`
			UserType             int         `json:"user_type"`
			UserIsMajia          bool        `json:"user_is_majia"`
			UserLevel            int         `json:"user_level"`
			UserFirstEnName      string      `json:"user_first_en_name"`
			UserAge              string      `json:"user_age"`
			UserInfluLevel       int         `json:"user_influ_level"`
			UserBlackType        int         `json:"user_black_type"`
			UserThirdIntro       interface{} `json:"user_third_intro"`
			UserBizflag          string      `json:"user_bizflag"`
			UserBizsubflag       string      `json:"user_bizsubflag"`
			PassportMedalDetails struct {
				MedalList  []interface{} `json:"medal_list"`
				MedalCount int           `json:"medal_count"`
			} `json:"passport_medal_details"`
			UserExtendinfos struct {
				Certification string `json:"certification"`
			} `json:"user_extendinfos"`
		} `json:"post_user"`
		PostTitle        string      `json:"post_title"`
		PostAbstract     string      `json:"post_abstract"`
		PostLastTime     string      `json:"post_last_time"`
		PostClickCount   int         `json:"post_click_count"`
		PostCommentCount int         `json:"post_comment_count"`
		PostAddress      interface{} `json:"post_address"`
	} `json:"re"`
	IntelligentReply int `json:"intelligent_reply"`
	ClassicReply     int `json:"classic_reply"`
}

// Define the struct to match the JSON response for the post data
type GubaResponse struct {
	Post struct {
		PostID          int64  `json:"post_id"`
		PostTitle       string `json:"post_title"`
		PostContent     string `json:"post_content"`
		PostAbstract    string `json:"post_abstract"`
		PostPublishTime string `json:"post_publish_time"`
		PostLastTime    string `json:"post_last_time"`
		PostDisplayTime string `json:"post_display_time"`

		PostClickCount       int    `json:"post_click_count"`
		PostForwardCount     int    `json:"post_forward_count"`
		PostCommentCount     int    `json:"post_comment_count"`
		PostCommentAuthority int    `json:"post_comment_authority"`
		PostLikeCount        int    `json:"post_like_count"`
		PostIsLike           bool   `json:"post_is_like"`
		PostIPAddress        string `json:"post_ip_address"`
	} `json:"post"`
	RC int    `json:"rc"`
	Me string `json:"me"`
}

// Define the struct to match the JSON response for QA
type QAResponse struct {
	Re []struct {
		PostID           int64  `json:"post_id"`
		UserID           string `json:"user_id"`
		UserNickname     string `json:"user_nickname"`
		PostClickCount   int    `json:"post_click_count"`
		PostForwardCount int    `json:"post_forward_count"`
		PostCommentCount int    `json:"post_comment_count"`
		PostLikeCount    int    `json:"post_like_count"`
		PostState        int    `json:"post_state"`
		PostFromNum      int    `json:"post_from_num"`
		AskQuestion      string `json:"ask_question"`
		AskAnswer        string `json:"ask_answer"`
		PostIsCollected  bool   `json:"post_is_collected"`
		ReplyCount       int    `json:"reply_count"`
	} `json:"re"`
	TotalPage        int    `json:"TotalPage"`
	Count            int    `json:"count"`
	PageIndex        int    `json:"PageIndex"`
	PageSize         int    `json:"PageSize"`
	CurrentSecretary string `json:"CurrentSecretary"`
	StockName        string `json:"StockName"`
}

func HandleNews(midCode string, pageSize int16) (NewsResponse, error) {
	url := "https://np-seclist.eastmoney.com/sec/getQuoteNews"
	requestBody := map[string]interface{}{
		"appKey":        "fd374bf183b866ce5cf7b00b92bb9858",
		"biz":           "sec_quote_news",
		"client":        "sec_android",
		"clientVersion": "10.23",
		"deviceId":      "",
		"midCode":       midCode,
		"pageSize":      pageSize,
		"req_trace":     "fac289aa-39ed-4bf1-b8b6-fd999a7d1b2a",
		"sortEnd":       "",
	}
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return NewsResponse{}, err
	}
	var response NewsResponse
	err = sendPostRequest(url, headers, bytes.NewBuffer(jsonData), &response)
	if err != nil {
		return NewsResponse{}, err
	}
	// 处理新闻内容
	var artcodes []string
	for i := range response.Data.News {
		code := response.Data.News[i].Code
		artcodes = append(artcodes, code)

	}
	if result, err := EastmoneyArticleService(artcodes); err == nil {
		for i := range response.Data.News {
			code := response.Data.News[i].Code
			response.Data.News[i].Content = result[code].(string)
		}
	}

	return response, nil
}

// 股吧帖子内容
func guba(postId string) (GubaResponse, error) {
	param := "ctoken=&utoken=&deviceid=$IP$&version=10023000&product=EastMoney&plat=Android&deviceId=88f2572d23a9e2b97d277a16fb8b759b||iemi_tluafed_me&postid=%s&IsMatch=false&type=0&cutword=true&paytext=true"
	url := "https://emcreative.eastmoney.com/FortuneApi/GuBaApi/common"
	requestBody := map[string]interface{}{
		"url":   "content/api/Post/ArticleContent",
		"type":  "post",
		"sumit": "form",
		"parm":  fmt.Sprintf(param, postId),
	}
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return GubaResponse{}, err
	}
	var response GubaResponse
	err = sendPostRequest(url, headers, bytes.NewBuffer(jsonData), &response)
	if err != nil {
		return GubaResponse{}, err
	}
	re := regexp.MustCompile(`<.*?>`)
	content := re.ReplaceAllString(response.Post.PostContent, "")
	re = regexp.MustCompile("&nbsp;")
	content = re.ReplaceAllString(content, " ")
	response.Post.PostContent = content
	fmt.Println(content)
	return response, nil
}

// 热评-机构
func HandleOrgList(code string, pageSize int16) (OrgResponse, error) {
	url := "https://gbapi-pj.emhudong.cn/apparticlelistv2/api/Hot/Articlelist"

	// Request headers
	headers := map[string]string{
		"EM-OS":        "Android",
		"Content-Type": "application/x-www-form-urlencoded",
	}

	// Request body SZ300059
	baseParam := "ctoken=&p=1&product=StockWay&randomtext=-orAWxN_q3vguTIJ7PvjWHoRvK7dRWYzD9BaoR1gzTfpxPn9GHAViLuQW0tamOlGrOioZ_fcj1KhmlxRNt2X8A&ps=%s&code=%s&pi=&utoken=&plat=Android&type=2&deviceid=f0cca1bbe2ddb67cfd232a85e7820e6a&version=10023000"
	data := fmt.Sprintf(baseParam, pageSize, code)

	var response OrgResponse
	err := sendPostRequest(url, headers, strings.NewReader(data), &response)
	if err != nil {
		return OrgResponse{}, err
	}

	return response, nil
}
func SpiderOrg(orgs []string) (OrgList, error) {
	var wg sync.WaitGroup

	chans := make(chan GubaResponse, 5)
	for _, org := range orgs {
		fmt.Println(org)
		wg.Add(1)
		go func() {
			defer wg.Done()
			if resp, err := guba(org); err == nil {
				fmt.Println(resp)
				chans <- resp
			}
		}()
	}
	go func() {
		wg.Wait()
		close(chans)
	}()
	result := OrgList{}
	for item := range chans {
		result.List = append(result.List, item)
	}
	return result, nil
}

// 董秘 code SH600519
func QA(code string, pageSize int16) (QAResponse, error) {
	url := "https://mgubaqa.eastmoney.com/interface/GetData.aspx"
	requestBody := map[string]interface{}{
		"param": fmt.Sprintf("code=%s&qatype=1&p=1&ps=%d&keyword=&questioner=", code, pageSize),
		"path":  "question/api/info/Search",
	}
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return QAResponse{}, err
	}

	var response QAResponse
	err = sendPostRequest(url, headers, bytes.NewBuffer(jsonData), &response)

	if err != nil {
		return QAResponse{}, err
	}
	return response, nil
}

func sendPostRequest(url string, headers map[string]string, requestBody io.Reader, responseStruct interface{}) error {

	req, err := http.NewRequest("POST", url, requestBody)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %v", err)
	}

	err = json.Unmarshal(body, responseStruct)
	if err != nil {
		return fmt.Errorf("error parsing JSON: %v", err)
	}

	return nil
}

func GenereateArticle(code string) (string, error) {
	fileName := fmt.Sprintf("generate-%s.txt", code)
	// if fileInfo, err := os.Stat(fileName); err == nil {
	// 	if os.IsNotExist(err) {
	// 		fmt.Printf("文件 %s 不存在。\n", fileName)
	// 	} else {
	// 		creationTime := fileInfo.ModTime()
	// 		fmt.Printf("文件 %s 存在，创建日期为: %s\n", fileName, creationTime.Format(time.RFC3339))
	// 		return "", nil
	// 	}
	// }
	// 生产素材文件
	GeneratePaper(code)
	// ai 处理论据
	// handleCompany(code)
	// ai 生成文章
	q := `你是一名独特的财经分析师,你不仅具有系统的证券知识,在表达上还独具一格.你有良好的信息整合能力,能结合历史案例对比分析;有犀利的洞察力,能敏锐地发现潜在问题和机会;有良好的沟通与表达能力,能幽默、深入浅出的表达观点；根据素材，谈谈对这家公司的见解`
	if answer, err := getAnswer(parseEvidence(code), q); err == nil {
		fileHelp := oss.FileHelper{}
		fileHelp.WriteTxt(fileName, answer)
		return answer, nil
	} else {
		return "", err
	}
}
func UpdatePaper(code string, content string) {

}
func GeneratePaper(code string) interface{} {
	fileHelp := oss.FileHelper{}
	stockPaper := make(map[string]interface{})

	fileName := fmt.Sprintf("stockPaper-%s.json", code)
	if fileInfo, err := os.Stat(fileName); err == nil {
		if os.IsNotExist(err) {
			fmt.Printf("文件 %s 不存在。\n", fileName)
		} else {
			creationTime := fileInfo.ModTime()
			fmt.Printf("文件 %s 存在，创建日期为: %s\n", fileName, creationTime.Format(time.RFC3339))

			fileHelp.ReadFromFile(fileName, &stockPaper)
			return stockPaper
		}
	}

	code_prefix_num := "0."
	code_prefix_str := "SZ"
	if code[0] == '6' {
		code_prefix_num = "1."
		code_prefix_str = "SH"
	}
	code_num := fmt.Sprintf("%s%s", code_prefix_num, code)
	code_str := fmt.Sprintf("%s%s", code_prefix_str, code)

	if resp, err := HandleNews(code_num, 20); err == nil {
		stockPaper["news"] = resp
	} else {
		fmt.Println(err)
	}
	if resp, err := HandleOrgList(code_str, 20); err == nil {
		stockPaper["orgs"] = resp
	} else {
		fmt.Println(err)
	}
	if resp, err := QA(code_str, 20); err == nil {
		stockPaper["qa"] = resp
	} else {
		fmt.Println(err)
	}
	// fileHelp.WriteToFile(fileName, stockPaper)
	return stockPaper
}

func GenerateEvidence(items OrgList) (interface{}, error) {

	result := make(map[string]interface{})

	q := "你是一名学者或研究人员,从文章中找出分论点以及对应论据,多个论据使用;拼接为一行.省略开头或结尾的概括性内容,使用markdown格式返回"
	for _, item := range items.List {
		postId := item.Post.PostID
		background := item.Post.PostContent
		key := fmt.Sprintf("orgs-%d", postId)
		// background := "11111111111111"

		if answer, err := getAnswer(background, q); err == nil {
			result[key] = parseEvidenceOne(answer)
			fmt.Println(answer)
		} else {
			fmt.Println(err)
			return nil, err
		}
	}

	return result, nil
}
func parseEvidence(code string) string {
	fileHelp := oss.FileHelper{}
	fileName := fmt.Sprintf("evidence-%s.json", code)
	var data map[string]interface{}
	fileHelp.ReadFromFile(fileName, &data)
	var str string
	for _, value := range data {
		content := parseEvidenceOne(value.(string))
		str = str + content
	}
	return str
}

func parseEvidenceOne(content string) string {
	re := regexp.MustCompile("```markdown([\\s\\S]+)```")
	str := ""
	if matches := re.FindStringSubmatch(content); len(matches) > 1 {
		re = regexp.MustCompile(`(?s)\*\*.*?\*\*`)
		content = re.ReplaceAllString(matches[1], "")

		re = regexp.MustCompile(`论据\d?[:：]`)
		content = re.ReplaceAllString(content, "")

		arrs := strings.Split(content, "\n")
		for _, item := range arrs {
			re = regexp.MustCompile(`^#.*\n?`)
			content = re.ReplaceAllString(item, "")
			// re = regexp.MustCompile(`^-.*\n?`)
			// content = re.ReplaceAllString(content, "")

			content = strings.ReplaceAll(content, "-", "")
			content = strings.TrimSpace(content) // 删除前后空格
			re = regexp.MustCompile(`^[:：]\s*`)
			content = re.ReplaceAllString(content, "")
			if len(content) == 0 {
				continue
			}
			str = str + content + ";"
		}
	}
	return str
}
func getAnswer(background string, q string) (string, error) {
	messages := []common.Message{
		{Role: "system", Content: background},
		{Role: "user", Content: q},
	}
	if answer, err := common.QDeepSeek(messages); err == nil {
		return answer, nil
	} else {
		return "", err
	}
}
