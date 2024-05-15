package routers

import (
	"github.com/gin-gonic/gin"
	"zg5/work/work07/client/api"
)

func InitRouter(Group *gin.RouterGroup) {
	user := Group.Group("user")
	{
		user.POST("/login", api.Login)
		user.POST("/register", api.Register)
	}
	//goods := Group.Group("goods")
	//{
	//	goods.POST("preheat", api.Preheat)
	//	goods.POST("secJi", api.SecJi)
	//	goods.POST("checkRegularly", api.CheckRegularly)
	//
	//}
}
