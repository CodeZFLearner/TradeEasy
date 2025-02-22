
## 腾讯云
2014-12-22 10:00: go 下载腾讯云文件

```bash

go get -u github.com/tencentyun/cos-go-sdk-v5

# io编码
go get -u golang.org/x/text

```

## 备份
```go

if list, exists := originMap["data"]; exists {
		if listSlice, e := list.([]interface{}); e {
			jdataList = listSlice
		} else if listMap, ok := list.(map[string]interface{}); ok {
			if diffs, e := listMap["diff"]; e {
				if listSlice, e := diffs.([]interface{}); e {
					jdataList = listSlice
				}
			} else {
				fmt.Println("diffs key does not exist")
			}
		}
	} else if list, exists := originMap["result"]; exists {
		if tempList, exists := list.(map[string]interface{}); exists {
			if data, e := tempList["data"]; e {
				if listSlice, e := data.([]interface{}); e {
					jdataList = listSlice
				} else {
					fmt.Println("result data key does not exist")
				}
			}
		}
	
	}

```
## 东方新闻
https://finance.eastmoney.com/a/202501173300826597.html
artCode
正则匹配文章主体：文章主体-->([\n\w\W]{0,})文尾部其它信息
