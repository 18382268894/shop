/**
*FileName: controller
*Create on 2018/12/5 上午4:01
*Create by mok
*/

package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Order(c *gin.Context){
	c.HTML(http.StatusOK,"order",nil)
}