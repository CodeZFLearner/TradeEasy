<!--
 * @Author: zff 2059577798@qq.com
 * @Date: 2025-02-22 14:31:08
 * @LastEditors: zff 2059577798@qq.com
 * @LastEditTime: 2025-03-08 00:53:07
 * @FilePath: \HelloGolang\readme.md
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
-->

## 东财热点新闻解析
1. v0.0
    ```go
    q := "从政策与国际动态、股市与市场表现、公司动态与业绩、行业趋势与技术发展 和 机构观点与策略 五个方面对源数据中信息点分类。并按照以下要求返回信息点:" +
		"1.信息点是数据源中主要表达的内容，返回数据源所有信息点 2.同一个分类下信息点在主题、内容、来源上不能重复或相似,如果具有相似性，就将信息点合并后简要概括。3.过滤掉简短的、含量低的信息点，并按照和分类主题的相关度降序排序。" +
		"4.信息点不能重复,尽可能多的梳理出所有信息点;如果一个分类下没有信息点，就省略整个分类。省略开头和结尾概括性内容"
	messages := []common.Message{
		{Role: "system", Content: "这是源数据，问题的答案都基于此数据:" + datasource},
		{Role: "user", Content: q},
	}
    ```


### 项目配置

    