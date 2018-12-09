/**
*FileName: accountController
*Create on 2018/12/8 上午2:10
*Create by mok
*/

package accountController

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Index(c *gin.Context){
	c.HTML(http.StatusOK,"account/index",nil)
}