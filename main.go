/**
*FileName: shop
*Create on 2018/12/5 上午12:15
*Create by mok
*/

package main

import (
	"github.com/gin-gonic/gin"
	"shop/router"
	"shop/db"
	"log"
	"shop/pkg/session"
)
func main(){
	r := gin.Default()
	err := db.Init()
	defer db.Close()
	if err != nil{

	}
	session.Init()
	router.LoadRouter(r)
	if err = r.Run(":8080");err != nil{
		log.Fatalln(err.Error())
	}
}