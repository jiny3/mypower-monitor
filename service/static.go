package service

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Static(r *gin.Engine) {
	r.Use(cors)
	// 静态文件服务
	static(r)
}

// CORS中间件
func cors(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(200)
	} else {
		c.Next()
	}
}

// static 文件夹存放静态文件
func static(r *gin.Engine) {
	r.LoadHTMLFiles("static/index.html")
	r.GET("/:roomid", func(c *gin.Context) {
		roomid := c.Param("roomid")
		c.HTML(http.StatusOK, "index.html", gin.H{
			"roomid": roomid,
		})
	})
	r.GET("/static/echarts.js", func(c *gin.Context) {
		c.File("static/echarts.js")
	})
	r.GET("/static/my.js/:roomid", func(c *gin.Context) {
		roomid := c.Param("roomid")
		// 从my.js文件中读取内容并替换其中的ID
		myjsContent, err := os.ReadFile("static/my.js")
		if err != nil {
			logrus.WithError(err).Error("读取my.js文件失败")
			c.JSON(http.StatusGatewayTimeout, gin.H{
				"message": "请求my.js文件失败",
			})
			return
		}
		c.Data(http.StatusOK, "text/javascript", []byte(strings.ReplaceAll(string(myjsContent), "{{.roomid}}", roomid)))
	})
}
