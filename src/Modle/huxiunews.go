package modle

import (
	"encoding/json"
	"github.com/olivere/elastic"
	"context"
)

//虎嗅数据modle
type HuxiuNews struct {
	Title   string
	Url     string
	ImgSrc  string
	TimeGap string
	Desc    string
	Column  string
}

func FromJsonObjHuxiu(o interface{}) (HuxiuNews, error) {
	var news HuxiuNews
	s, err := json.Marshal(o)
	if err != nil {
		return news, err
	}
	err = json.Unmarshal(s, &news)
	return news, err
}

func GetData() (news []Item, err error) {
	const index = "dating_profile_2"

	//todo Try to start up elastic search
	client, err := elastic.NewClient()
	if err != nil {
		return nil, err
	}

	resp, err := client.Search().
		Index(index).
		Type("huxiu"). // search in index "twitter"
		Pretty(true). // pretty print request and response JSON
		Query(elastic.Query())
		Do(context.Background()) // execute

	if err != nil {
		return nil, err
	}

	hits := resp.Hits.Hits


	for _, hit := range hits {
		var item Item

		hitJosnBuf, err := hit.Source.MarshalJSON()
		if err != nil {
			continue
		}
		err = json.Unmarshal(hitJosnBuf, &item)
		if err != nil {
			continue
		}
		itemNew, err := FromJsonObjHuxiu(item.Payload)
		if err != nil {
			continue
		}
		item.Payload = itemNew
		news = append(news, item)
	}

	//fmt.Printf("news: %s",news)
	return news,nil
}
