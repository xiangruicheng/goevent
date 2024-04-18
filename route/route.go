package route

import (
	"github.com/gin-gonic/gin"
	"goevent/controller"
)

func InitRoute() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	//注册器
	topic := r.Group("/topic")
	topic.POST(":name/:data", controller.TopicController.Post)
	topic.GET(":name/:type", controller.TopicController.Get)
	topic.DELETE(":name", controller.TopicController.Delete)
	topic.PUT(":name", controller.TopicController.Put)

	return r
}
