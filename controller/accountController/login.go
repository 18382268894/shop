/**
*FileName: controller
*Create on 2018/12/5 上午4:12
*Create by mok
 */

package accountController

import (
	"crypto/md5"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"shop/db/accountDB"
	"shop/pkg/session"
	"shop/utils"
	"time"
)

//用户登录
func Login(c *gin.Context) {

	email := c.PostForm("email")
	password := c.PostForm("password")
	//参数检查
	err := utils.CheckParams([]*string{&email, &password}, utils.WithTram(), utils.WithNotEmpty(), utils.WithEmailRule(0))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	//登录检查和数据更新
	loginIP, _ := utils.IPToInt(c.ClientIP())
	password = fmt.Sprintf("%x", md5.Sum([]byte(password)))
	u, err := accountDB.Login(email, password, loginIP)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "登录失败，请稍后再试",
		})
		log.Println(err.Error())
		return
	}
	//保存登录信息到session
	var s session.Session
	if s, err = session.Mgr.CreateSession(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	s.Set("uid", u.UID)
	s.Set("username", u.Username)
	s.Set("email", u.Email)

	//生成双token
	//短token
	shortClaim := &utils.MyClaim{
		Data: map[string]interface{}{"uid": u.UID, "email": u.Email, "username": u.Username},
		StandardClaims: &jwt.StandardClaims{
			Issuer:    "电子商城",
			ExpiresAt: time.Now().Add(2*time.Minute).Unix(),
			IssuedAt:  time.Now().Unix(),
			NotBefore: time.Now().Unix(),
			Subject:   "登录状态验证",
		},
	}
	longClaim := &utils.MyClaim{
		Data: map[string]interface{}{"sessionid": s.ID(), "ip": loginIP},
		StandardClaims: &jwt.StandardClaims{
			Issuer:    "电子商城",
			ExpiresAt: time.Now().Add(30 * 24 * time.Hour).Unix(), //过期时间30天
			IssuedAt:  time.Now().Unix(),
			NotBefore: time.Now().Unix(),
			Subject:   "登录状态验证",
		},
	}
	shortToken, _ := utils.NewToken(shortClaim)
	longToken, _ := utils.NewToken(longClaim)
	c.SetCookie("login_short_token",shortToken,150,"/","",false,true)
	c.SetCookie("login_long_token",longToken,30*24*3600,"/","",false,true)
	c.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
	})
	return
}
