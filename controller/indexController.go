/**
*FileName: controller
*Create on 2018/12/5 上午12:35
*Create by mok
*/

package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func IndexController(c *gin.Context){
	c.HTML(http.StatusOK,"index",nil)
}