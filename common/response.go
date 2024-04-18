package common

import "github.com/gin-gonic/gin"

type Response struct {
	Code int
	Msg  string
	Data any
}

var (
	RSP_SUCC        = Response{Code: 0, Msg: "success"}
	RSP_PARAM_ERROR = Response{Code: 1000, Msg: "参数错误"}
	RSP_HANDLE_FAIL = Response{Code: 1001, Msg: "处理失败"}
)

func Sucess(ctx *gin.Context, Data any) {
	obj := gin.H{
		"code": RSP_SUCC.Code,
		"msg":  RSP_SUCC.Msg,
		"data": Data,
	}
	ctx.JSON(200, obj)
}

func Error(ctx *gin.Context, response Response) {
	obj := gin.H{
		"code": response.Code,
		"msg":  response.Msg,
		"data": response.Data,
	}
	ctx.JSON(200, obj)
}
