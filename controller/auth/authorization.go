package controllers

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	controllers "goEs/controller"
	"goEs/middleware/authMiddleware"
	"goEs/utils"
	"time"
)

func GetToken(c *gin.Context) {
	HmacUser := authMiddleware.HmacUser{Email: "1234567@qq.com", Password: "12345678"}
	t, e := authMiddleware.GenerateToken(HmacUser)
	if e != nil {
		utils.Logger.Error("Unauthorized", zap.String("method", "GetToken"), zap.Duration("backoff", time.Second))
		c.JSON(401, "Unauthorized")
	}
	rep := controllers.SuccessRep(t)
	c.JSON(200, rep)
}
