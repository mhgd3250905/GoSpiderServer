package meizi

import (
	"github.com/gin-gonic/gin"
	"fmt"
)

func GetMeizi(c *gin.Context) {
	fileName:=c.Param("filename")
	c.File(fmt.Sprint("./src/file/",fileName))
}