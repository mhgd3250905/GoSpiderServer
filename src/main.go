package main

import (
	"GoSpiderServer/src/modle"
	"github.com/gin-gonic/gin"
	"net/http"
	"fmt"
	"strconv"
	"GoSpiderServer/src/biliFans"
	"GoSpiderServer/src/chat"
	"GoSpiderServer/src/meizi"
)

var mapIndex map[string]int

func main() {
	//modle.GetDataFromES()
	router := gin.Default()
	router.LoadHTMLGlob("./src/template/*")

	router.GET("/", home)

	//新闻爬虫
	groupSpider := router.Group("/spider/")
	{
		groupSpider.GET("/:type/:index", getData)
		groupSpider.GET("/:type/", getData)
	}

	//新闻爬虫Api
	groupApiSpider := router.Group("/api/v1/spider/")
	{
		groupApiSpider.GET("/:type/:index",getApiData)
		groupApiSpider.GET("/:type/",getApiData)
	}

	//bili爬虫工具
	groupBili := router.Group("/bili/")
	groupBili.GET("/query", biliQuery)
	groupBili.POST("/query", biliQuery)

	//聊天室
	router.GET("/echo", chat.Echo)

	//下载工具
	fileDownload := router.Group("/file/")

	fileDownload.GET("/:filename", meizi.GetMeizi)

	//groupSpiderApi := router.Group("/spider/api")
	//groupSpiderApi.GET("/huxiu", HuxiuApi)

	mapIndex = make(map[string]int)

	router.Run(":8080")

}

func home(c *gin.Context) {

	c.Redirect(http.StatusMovedPermanently, "/spider/huxiu")
}

//根据请求参数获取数据
func getData(c *gin.Context) {
	pageType := c.Param("type")
	pageIndex := c.Param("index")

	index, err := strconv.Atoi(pageIndex)
	if err != nil {
		if pageIndex == "prev" {
			if mapIndex[pageType] > 0 {
				index = mapIndex[pageType] - 1
			} else {
				index = 0
			}
		} else if pageIndex == "next" {
			index = mapIndex[pageType] + 1
		} else {
			index = 0
		}
	}
	mapIndex[c.Param("type")] = index

	news, err := modle.GetDataFromRedis(pageType, (index)*20, 20)
	fmt.Println("news size->", len(news))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{`error`: 1, `data`: ``})
	}
	c.HTML(http.StatusOK, "news.html", gin.H{"news": news, "type": pageType})
}

//根据请求参数获取数据
func getApiData(c *gin.Context) {
	pageType := c.Param("type")
	pageIndex := c.Param("index")

	index, err := strconv.Atoi(pageIndex)
	if err != nil {
		if pageIndex == "prev" {
			if mapIndex[pageType] > 0 {
				index = mapIndex[pageType] - 1
			} else {
				index = 0
			}
		} else if pageIndex == "next" {
			index = mapIndex[pageType] + 1
		} else {
			index = 0
		}
	}
	mapIndex[c.Param("type")] = index

	news, err := modle.GetDataFromRedis(pageType, (index)*20, 20)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{`error`: 1, `data`: ``})
	}
	c.JSON(http.StatusOK, gin.H{"error": 0, "news": news})
}

//查询返回页面
func biliQuery(c *gin.Context) {
	if c.Request.Method == "GET" {
		c.HTML(http.StatusOK, "queryBiliFans.html", nil)
	} else {
		name := c.PostForm("upName")
		if name == "" {
			c.String(http.StatusOK, "未输入Up主姓名")
		}

		results, err := biliFans.QueryBiliFans([]string{name})
		if err != nil {
			c.String(http.StatusOK, "查询发生了错误。")
		}

		c.HTML(http.StatusOK, "queryBiliFans.html", gin.H{`Users`: results})
	}
}
