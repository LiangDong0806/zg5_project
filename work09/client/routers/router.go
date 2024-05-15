package routers

import (
	"github.com/gin-gonic/gin"
	"zg5/work/work09/client/api"
)

func UserRouter(Group *gin.RouterGroup) {
	user := Group.Group("user")
	{
		user.POST("login", api.Login)
	}
}
