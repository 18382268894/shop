/**
*FileName: utils
*Create on 2018/12/6 上午6:14
*Create by mok
*/

package utils

import (
	"github.com/dgrijalva/jwt-go"
	"errors"
)


const (
	secret = "hhh-jjj-kkk-lll"
)
type MyClaim struct {
	Data map[string]interface{} `json:"data"`
	*jwt.StandardClaims
}


//生成token
func NewToken(myClaim *MyClaim)(tokenss string,err error){
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,myClaim)
	tokenss,err = token.SignedString([]byte(secret))
	return
}

//解析token
func ParseToken(tokenss string)(map[string]interface{}, error){
	token,err := jwt.ParseWithClaims(tokenss,&MyClaim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret),nil
	})
	if err != nil{
		return nil,err
	}
	claim,ok := token.Claims.(*MyClaim)
	if ok && token.Valid {
		return claim.Data,nil
	}
	return nil,errors.New("cannot parse the token")
}
