/**
*FileName: controller
*Create on 2018/12/5 上午3:35
*Create by mok
*/

package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Category(c *gin.Context){
	c.HTML(http.StatusOK,"category",nil)
}
