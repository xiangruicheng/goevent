package controller

import (
	"github.com/gin-gonic/gin"
	"goevent/common"
	"goevent/server"
	"goevent/util"
)

type TopicControllerType struct{}

var TopicController *TopicControllerType

func init() {
	TopicController = new(TopicControllerType)
}

// Post 插入数据
func (c *TopicControllerType) Post(ctx *gin.Context) {
	name := ctx.Param("name")
	data := ctx.Param("data")
	if name == "" || data == "" {
		common.Error(ctx, common.RSP_PARAM_ERROR)
	}
	affectRow, err := server.RedisClient.RPush(ctx, name, data).Result()
	if err != nil {
		common.Error(ctx, common.RSP_HANDLE_FAIL)
	}
	common.Sucess(ctx, affectRow)
}

// Get 获取数据
func (c *TopicControllerType) Get(ctx *gin.Context) {
	name := ctx.Param("name")
	typeStr := ctx.Param("type")
	if name == "" {
		common.Error(ctx, common.RSP_PARAM_ERROR)
		return
	}
	ok, _ := util.InArray(typeStr, []string{"pop", "all"})
	if !ok {
		common.Error(ctx, common.RSP_PARAM_ERROR)
		return
	}

	var data, err any
	switch typeStr {
	case "pop":
		data, err = server.RedisClient.LPop(ctx, name).Result()
		break
	case "all":
		data, err = server.RedisClient.LRange(ctx, name, 0, -1).Result()
	default:
	}
	if err != nil {
		common.Error(ctx, common.RSP_HANDLE_FAIL)
	}
	common.Sucess(ctx, data)
}

func (c *TopicControllerType) Delete(ctx *gin.Context) {

}

func (c *TopicControllerType) Put(ctx *gin.Context) {

}
