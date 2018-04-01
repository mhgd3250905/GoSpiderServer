package biliFans

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"encoding/json"
	"net/http"
	"net/url"
)

type BiliUserBean struct {
	Result []BiliUserResultBean
}

type BiliUserResultBean struct {
	Uname string
	Usign string
	Mid   int
	Fans  int
}

//填入名称
var biliSpiderUrl string = "https://search.bilibili.com/api/search?search_type=bili_user&keyword=%s"

//输入名字，返回对应的查询信息
func QueryBiliFans(names []string) (results []BiliUserResultBean, err error) {
	var name string
	reslutList:=make([]BiliUserResultBean,0)
	for i := 0; i < len(names); i++ {
		name = names[i]
		url := fmt.Sprintf(biliSpiderUrl, url.QueryEscape(name))

		//fmt.Println(url)

		client := &http.Client{}

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println("create http request failed!", err)
			return reslutList, err

		}
		req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3282.140 Safari/537.36")

		response, err := client.Do(req)
		if err != nil {
			fmt.Println("get url content fail ", err)
			return reslutList, err
		}

		//fmt.Println(string(response.Body.Read()))

		d, err := goquery.NewDocumentFromResponse(response)
		if err != nil {
			fmt.Println("访问bilibi失败，", err)
			return reslutList, err

		}

		userJsonStr := d.Text()
		//fmt.Println(userJsonStr)

		var userBean BiliUserBean
		err = json.Unmarshal([]byte(userJsonStr), &userBean)
		if err != nil {
			fmt.Println("json 2 struct failed,", err)
		}

		result:=userBean.Result


		for  i:=0; i<len(result);i++ {
			reslutList=append(reslutList, result[i])
		}

	}

	return reslutList, nil

	//fmt.Println("-------------------------------")
	//fmt.Printf("查询到结果 %d 个\n", len(userBean.Result))
	//fmt.Println("-------------------------------")
	//for _, result := range userBean.Result {
	//	fmt.Println("Up主姓名：", result.Uname)
	//	fmt.Println("Up主id：", result.Mid)
	//	fmt.Println("Up主标记：", result.Usign)
	//	fmt.Println("粉丝数：", result.Fans)
	//	fmt.Println("-------------------------------")
	//}
}
