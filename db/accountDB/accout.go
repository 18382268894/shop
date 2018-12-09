/**
*FileName: accountDB
*Create on 2018/12/8 上午2:43
*Create by mok
*/

package accountDB

import (
	"time"
	"shop/db"
	"shop/model"
	"github.com/bwmarrin/snowflake"
)

func InsertUser(email,username,password string)(uint64,error){
	node,_ := snowflake.NewNode(1)
	var u  = &model.User{
		UID:uint64(node.Generate()),
		Email:email,
		Username:username,
		Password:password,
		CreateTime:time.Now().Unix(),
	}
	return u.UID,db.DB.Create(u).Error
}


func ActiveAccount(UID uint64)error{
	return db.DB.Model(&model.User{}).Where("uid=?",UID).Update("status",model.StatusUserActivated).Error
}


func Login(email,password string,loginIP int)(*model.User,error){
	var u = new(model.User)
	err := db.DB.Model(&model.User{}).
		Where("email=? and password=?",email,password).
		Select([]string{"email","username","uid"}).
		Find(u).Error
	if err != nil{
		return nil,err
	}
	loginTime := time.Now().Unix()
	err = db.DB.Model(&model.User{}).Where("uid=?",u.UID).Update("login_ip",loginIP,"login_time",loginTime).Error
	u.Email = email
	return u,err
}