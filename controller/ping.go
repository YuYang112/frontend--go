package controllers

import (
	"github.com/gin-gonic/gin"
	"goEs/middleware/authMiddleware"
)

//测试服务运行状态
func TestPing(ctx *gin.Context) {
	HmacUser := authMiddleware.HmacUser{Email:"1234567@qq.com",Password:"12345678"}
	t,e:= authMiddleware.GenerateToken(HmacUser)
	if e!= nil{
		rep := SuccessRep("测试")
		ctx.JSON(400,rep)
	}
	rep := SuccessRep(t)
	ctx.JSON(200,rep)
}

//测试服务运行状态
func TestPing1(ctx *gin.Context) {
	rep := SuccessRep("测试成功")
	ctx.JSON(200,rep)
}
