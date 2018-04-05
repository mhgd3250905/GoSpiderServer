package modle

import (
	"encoding/json"
	"github.com/olivere/elastic"
	"context"
	"github.com/garyburd/redigo/redis"
	"github.com/pkg/errors"
	"fmt"
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

func GetDataFromES() (news []Item, err error) {
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
	return news, nil
}

func GetDataFromRedis(key string,offset int, limite int) (news []Item, err error) {
	const index = "dating_profile_2"

	//创建一个Redis实例
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	//输入账号密码
	//_, err = conn.Do("AUTH", "sk3250905")
	//if err != nil {
	//	return nil, errors.Errorf("Redis AUTH fail %v", err)
	//}
	//测试连接
	if result, err := conn.Do("ping"); result != "PONG" {
		return nil, errors.Errorf("Redis ping fail %v", err)

	}

	result, err := conn.Do("ZREVRANGE", key, offset, offset+limite-1)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	for i := 0; i < len(result.([]interface{})); i++ {
		itemStr := fmt.Sprintf("%s", result.([]interface{})[i])
		var item Item
		json.Unmarshal([]byte(itemStr), &item)
		if item.Id=="" {
			continue
		}
		itemNew, err := FromJsonObjHuxiu(item.Payload)
		if err != nil {
			continue
		}
		item.Payload = itemNew
		news = append(news, item)
	}
	return news, nil
}
