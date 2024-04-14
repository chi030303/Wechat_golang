package router

import (
	"Wechat-project/controllers"
	"github.com/gin-gonic/gin"
)

// 将请求封装到Router函数中
func Router() *gin.Engine {
	r := gin.Default()

    // 总路由
	check := r.Group("/wx")
    check.Use(gin.Recovery())
    {
        check.GET("/v1", controllers.WxVerify)
        // 对用户发消息做出回复
        check.POST("/v1",controllers.Reply)
    }
	return r
}

