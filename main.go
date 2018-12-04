/**
*FileName: shop
*Create on 2018/12/5 上午12:15
*Create by mok
*/

package main

import (
	"github.com/gin-gonic/gin"
	"shop/controller"
)
func main(){
	r := gin.Default()
	r.LoadHTMLGlob("./views/*")
	r.Static("/static","./public")
	r.GET("/index",controller.IndexController)
	r.GET("/category",controller.Category)
	r.GET("/cargo",controller.CargoDetail)
	r.GET("/cart",controller.Cart)
	r.Run(":8080")
}