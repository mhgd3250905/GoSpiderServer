package main

import (
	"GoSpiderServer/src/Modle"
	"github.com/gin-gonic/gin"
	"net/http"
	"fmt"
	"strconv"
)

var mapIndex map[string]int

func main() {
	//modle.GetDataFromES()
	router := gin.Default()
	router.LoadHTMLGlob("src/template/*")

	groupSpider := router.Group("/spider/")
	groupSpider.GET("/:type/:index", getData)
	groupSpider.GET("/:type/", getData)
	router.GET("/", nil)

	//groupSpiderApi := router.Group("/spider/api")
	//groupSpiderApi.GET("/huxiu", HuxiuApi)

	mapIndex = make(map[string]int)

	router.Run()

}

func getData(c *gin.Context) {
	pageType := c.Param("type")
	pageIndex := c.Param("index")

	if pageType == "huxiu" {
		index, err := strconv.Atoi(pageIndex)
		if err != nil {
			if pageIndex == "prev" {
				if mapIndex[pageType] > 0 {
					index = mapIndex[pageType]-1
				} else {
					index = 0
				}
			} else if pageIndex == "next" {
				index = mapIndex[pageType]+1
			} else {
				index = 0
			}
		}
		mapIndex[c.Param("type")]=index

		news, err := modle.GetDataFromRedis((index)*20, 20)
		fmt.Println("news size->", len(news))
		if err != nil {
			c.JSON(http.StatusOK, gin.H{`error`: 1, `data`: ``})
		}
		c.HTML(http.StatusOK, "news.html", gin.H{"news": news, "type": pageType})
	}
}

func HuxiuApi(c *gin.Context) {
	news, err := modle.GetDataFromES()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{`error`: 1, `data`: ``})
	}

	c.JSON(http.StatusOK, gin.H{`error`: 0, `data`: news})
}
