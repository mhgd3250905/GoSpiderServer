package modle

import "GoSpiderServer/src/biliFans"

type Item struct{
	Url string
	Type string
	Id string
	Payload interface{}
}

type HtmlBean struct {
	Users []biliFans.BiliUserResultBean
}
