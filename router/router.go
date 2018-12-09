/**
*FileName: router
*Create on 2018/12/8 上午2:13
*Create by mok
*/

package router

import (
	"github.com/gin-gonic/gin"
	"shop/controller/accountController"
)

func LoadRouter(r *gin.Engine){
	r.LoadHTMLGlob("./views/*")
	r.Static("/static","./public")

	accountGroup := r.Group("/account")
	{
		accountGroup.GET("/index",accountController.Index)
		accountGroup.POST("/register",accountController.Register)
		accountGroup.POST("/login",accountController.Login)
	}

}