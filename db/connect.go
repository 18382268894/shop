/**
*FileName: db
*Create on 2018/12/5 上午5:56
*Create by mok
*/

package db

import (
	_"github.com/go-sql-driver/mysql"
	"fmt"
	"shop/config"
	"time"
	"log"
	"github.com/jinzhu/gorm"
	"shop/model"
)
var DB *gorm.DB
func Init()error{
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", config.Mysql_User,config.Mysql_Password,config.Mysql_Host,config.Mysql_Port,config.Mysql_DBname)
	DB,_ =gorm.Open("mysql",dsn)
	DB.DB().SetConnMaxLifetime(100*time.Second)
	DB.DB().SetMaxIdleConns(5)
	DB.DB().SetMaxOpenConns(20)
	err := DB.DB().Ping()
	if err != nil{
		return err
		log.Println(err.Error())
	}
	err = createTables()
	if err != nil{
		return err
		log.Println(err.Error())
	}

	return nil
}

func Close(){
	DB.Close()
}

//创建表
func createTables()error{
	if !DB.HasTable(&model.User{}){
		return DB.Set("gorm:table_options","ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").CreateTable(model.User{}).Error
	}
	return nil
}
