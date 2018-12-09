/**
*FileName: accountController
*Create on 2018/12/8 上午2:10
*Create by mok
*/

package accountController

import (
	"github.com/gin-gonic/gin"
	"shop/utils"
	"net/http"
	"shop/db/accountDB"
	"log"
	"github.com/dgrijalva/jwt-go"
	"time"
	"html/template"
	"fmt"
	"crypto/md5"
)

//用户注册:传入参数username,password,email
func Register(c * gin.Context){
	//获取表格参数
	username := c.PostForm("username")
	password := c.PostForm("password")
	email := c.PostForm("email")

	//参数校验
	err := utils.CheckParams([]*string{&username,&password,&email},utils.WithTram(),utils.WithNotEmpty(),utils.WithEmailRule(2))
	if err != nil{
		c.JSON(http.StatusBadRequest,gin.H{
			"message":"注册失败，参数不合法",
		})
		return
	}
	password = fmt.Sprintf("%x", md5.Sum([]byte(password)))
	//插入数据库
	var UID uint64
	if UID,err = accountDB.InsertUser(email,username,password);err != nil{
		c.JSON(http.StatusInternalServerError,gin.H{
			"message":"注册失败,请稍后重试",
		})
		log.Println(err.Error())
		return
	}

	//注册成功，发送邮件
	go func() {
		myClaim := &utils.MyClaim{
			Data:map[string]interface{}{"UID":UID},
			StandardClaims:&jwt.StandardClaims{
				Issuer:"电子商城",
				ExpiresAt:time.Now().Unix() + 300,
				IssuedAt:time.Now().Unix(),
				NotBefore:time.Now().Unix(),
				Subject:"注册激活",
			},
		}
		tokenss,err  := utils.NewToken(myClaim)
		if err == nil{
			tp,_ := template.ParseFiles("./views/checkEmail.html")
			var emailconf = &utils.EmailConfig{
				From:[]string{"1005914310@qq.com","电子商城"},
				To:[]string{email},
				Subject:"注册激活",
				Tp:tp,
				Messages:map[string]interface{}{"token":tokenss},
			}
			err = utils.SendEmail(emailconf)
		}
		log.Println(err)
	}()
	c.JSON(http.StatusOK,gin.H{
		"message":"注册账号成功，我们将会发送一条邮件到你邮箱用于账号激活，请查收",
	})
	return
}

//激活
func Active(c *gin.Context){
	token := c.Query("token")
	if token == ""{
		c.JSON(http.StatusBadRequest,gin.H{
			"message":"激活参数不能为空",
		})
		return
	}
	data,err := utils.ParseToken(token)
	if err != nil{
		c.JSON(http.StatusInternalServerError,gin.H{
			"message":"激活失败,请稍后再试",
		})
		log.Println(err.Error())
		return
	}
	UID,ok := data["UID"].(uint64)
	if !ok{
		c.JSON(http.StatusBadRequest,gin.H{
			"message":"激活失败,激活参数错误",
		})
		return
	}

	if err = accountDB.ActiveAccount(UID);err != nil{
		c.JSON(http.StatusInternalServerError,gin.H{
			"message":"激活失败,请稍后重试",
		})
		return
	}
	c.JSON(http.StatusOK,gin.H{
		"message":"激活账号成功",
	})
}
