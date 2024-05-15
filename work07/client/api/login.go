package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"zg5/work/work07/client/service"
	"zg5/work/work07/common"
)

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	user, _ := service.QueryTheUserss(username)

	if user.Id == 0 {
		c.JSONP(http.StatusAccepted, gin.H{
			"code": http.StatusAccepted,
			"msg":  "请先注册账号",
		})
		return
	}

	p, _ := service.DecryptThePassword([]byte(user.Password))
	pwd := string(p)
	fmt.Println(user.Password, "][][][][")
	if password != pwd {
		c.JSONP(http.StatusAccepted, gin.H{
			"code": http.StatusAccepted,
			"msg":  "密码输入错误",
		})
		return
	}
	token, _ := common.SetJwtToken(common.GAOWEIMING, time.Now().Unix(), 3600, "fsd")
	c.JSONP(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "登录成功",
		"data": token,
	})
	return
}
func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	mobile := c.Query("mobile")
	user, _ := service.QueryTheUserss(username)

	fmt.Print("1111111111", user)
	if user.Id != 0 {
		c.JSONP(http.StatusAccepted, gin.H{
			"code": http.StatusAccepted,
			"msg":  "账号已存在",
		})
		return
	}
	pwd, err := service.EncryptPasswords([]byte(password))
	if err != nil {
		c.JSONP(http.StatusAccepted, gin.H{
			"code": http.StatusAccepted,
			"msg":  "密码加密出错",
		})
		return
	}
	_, err = service.UserRegistration(service.User{
		Username: password,
		Password: string(pwd),
		Mobile:   mobile,
	})
	if err != nil {
		c.JSONP(http.StatusAccepted, gin.H{
			"code": http.StatusAccepted,
			"msg":  "注册失败",
		})
		return
	}
	//token, _ := common.SetJwtToken(common.GAOWEIMING, time.Now().Unix(), 3600, strconv.Itoa(user.Id))
	c.JSONP(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "注册成功",
	})
}

//fmt.Println(in.Username)
