// @Title  router.go
// @Description  gin的路由初始化及设置
// @Author  zhu jianquan 2020/8/5
// @Update  zhu jianquan 2020/8/7

package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	controllers "goEs/controller"
	auth "goEs/controller/auth"
	"goEs/controller/elasticsearch"
	"goEs/controller/index"
	"goEs/middleware/authMiddleware"
	"goEs/utils"
	"net/http"
)

var Router *gin.Engine
var count int64


func InitRoute() {


	Router = gin.Default()
	Router.Use(utils.Cors())
	// 传参 设定路由组 允许路由组使用路由
	//app_routers.Router(engine)
	Router.GET("/ping", controllers.TestPing)
	Router.POST("/getToken", auth.GetToken)
	Router.GET("/ping1", authMiddleware.AuthMiddleware(),controllers.TestPing1)
	Router.POST("/v1/elastic", elasticsearch.Create)//添加文档
	Router.GET("/v1/elastic/:id", elasticsearch.GetElasticById)//获取文档
	Router.GET("/v1/elastic", elasticsearch.Query)//查询
	Router.POST("/v1/elastic/list", elasticsearch.List)//分页
	Router.DELETE("/v1/elastic/:id", elasticsearch.Destroy)//删除


	Router.PUT("/v1/index", index.Create)//设置副本分片数：127.0.0.1:9010/v1/index
	Router.GET("/v1/index/:index", index.IndexExists)//索引是否存在：127.0.0.1:9010/v1/index/twitter
	Router.DELETE("/v1/index/:index", index.DeleteIndex)//删除索引：127.0.0.1:9010/v1/index/twitter
	Router.GET("/v1/scroll/:index", index.ScrollParallel)
	Router.GET("/v1/flush/:index", index.Flush)//刷新：127.0.0.1:9010/v1/flush/twitter


	// 使用中间件
	Router.POST("/check_token", authMiddleware.AuthMiddleware(), func(c *gin.Context) {
	 if userName, exists := c.Get("userName"); exists {
	   fmt.Println("当前用户", userName)
	 }
	 c.JSON(http.StatusOK, gin.H{
	   "code":    0,
	   "message": "验证成功",
	 })
	})

	///api := Router.Group("/api", utils.ApiSignMIddleware())
	//
	////标准测试路由

	//api.POST("/elastic", elasticsearch.Post)
	////api.PUT("/news/:id", elasticsearch.Put)
	//api.DELETE("/elastic", elasticsearch.Destroy)

}

