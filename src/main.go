package main

import (
	"GoSpiderServer/src/Modle"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	//modle.GetData()
	router := gin.Default()
	router.LoadHTMLGlob("src/template/*")

	groupSpider := router.Group("/spider/")
	groupSpider.GET("/huxiu", Huxiu)
	router.GET("/", Huxiu)

	groupSpiderApi := router.Group("/spider/api")
	groupSpiderApi.GET("/huxiu", HuxiuApi)

	router.Run()

}

func Huxiu(c *gin.Context) {
	news, err := modle.GetData()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{`error`:1,`data`:``})
	}
	c.HTML(http.StatusOK, "news.html", gin.H{"news": news})
}

func HuxiuApi(c *gin.Context) {
	news, err := modle.GetData()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{`error`:1,`data`:``})
	}

	c.JSON(http.StatusOK, gin.H{`error`:0,`data`:news})
}
