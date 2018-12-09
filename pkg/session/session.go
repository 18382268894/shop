/**
*FileName: session
*Create on 2018/11/7 下午8:40
*Create by mok
*/

package session

import (
	"sync"
	"errors"
	"github.com/gomodule/redigo/redis"
	"encoding/json"
	"fmt"
	"time"
)

var(
	KeyNotExists = errors.New("key not exists")
)

const(
	SessionFlagNone = iota
	SessionFlagModify
)

type Session interface {
	Set(key string,val interface{})
	Get(key string)(interface{},error)
	MustGet(key string)(interface{})
	GetAll()(map[string]interface{},error)
	ID()(string)
	Delete(key string)
	Save()error
	Close()  //用于redis,mysql session关闭连接
}



type MemorySession struct {
	id  string
	kv sync.Map
}


func newMemorySession(id string)(ms *MemorySession){
	return &MemorySession{
		id:id,
	}
}


func (ms *MemorySession)Set(key string,val interface{}){
	ms.kv.Store(key,val)
}


func (ms *MemorySession)Get(key string)(val interface{},err error){
	var ok bool
	if val,ok = ms.kv.Load(key);ok{
		return
	}
	err = KeyNotExists
	return
}

func (ms *MemorySession)MustGet(key string)(interface{}){
	v,_ := ms.kv.Load(key)
	return v
}


func (ms *MemorySession)GetAll()(map[string]interface{},error){
	m := make(map[string]interface{})
	fn := func(key interface{},val interface{})bool{
		m[key.(string)] = val
		return true
	}
	ms.kv.Range(fn)
	return m,nil
}

func  (ms *MemorySession)ID()string{
	return ms.id
}

func(ms *MemorySession)Delete(key string){
	ms.kv.Delete(key)
}

func (ms *MemorySession)Save()error{
	return nil
}

func (ms *MemorySession)Close(){
	return
}



type RedisSession struct {
	id string
	kv sync.Map
	conn redis.Conn
	flag int
	expire time.Time
}


func newRedisSession(idString string,pool *redis.Pool)(*RedisSession){
	rs :=  &RedisSession{
		id:idString,
		conn:pool.Get(),
		flag:SessionFlagNone,
	}
	return rs
}

func(rs *RedisSession)Set(key string,val interface{}){
	rs.kv.Load(key)
	rs.flag = SessionFlagModify
}


func(rs * RedisSession)readFromRedis()(map[string]interface{},error){
	reply,err := rs.conn.Do("GET",rs.id)
	data,err := redis.String(reply,err)
	if err != nil{
		return nil,err
	}
	var m = make(map[string]interface{})
	err = json.Unmarshal([]byte(data),&m)
	return m,err
}


func(rs *RedisSession)Get(key string)(val interface{},err error){
	var m map[string]interface{}
	if rs.flag == SessionFlagNone{
		m,err = rs.readFromRedis()
		if err != nil{
			return
		}
		for k,v := range m{
			rs.Set(k,v)
		}
	}
	var ok bool
	if val,ok = rs.kv.Load(key); !ok{
		err = KeyNotExists
	}
	return
}



func(rs *RedisSession)MustGet(key string)(interface{}){
	val,_ := rs.Get(key)
	return val
}


func (rs *RedisSession)GetAll()(map[string]interface{},error){
	if rs.flag == SessionFlagNone{
		m ,err := rs.readFromRedis()
		return m,err
	}
	m := make(map[string]interface{})
	fn := func(key interface{},val interface{})bool{
		m[key.(string)] = val
		return true
	}
	rs.kv.Range(fn)
	return m,nil
}


func (rs *RedisSession)ID()string{
	return rs.id
}

func(rs *RedisSession)Delete(key string){
	rs.kv.Delete(key)
}


func (rs *RedisSession)Save()error{
	m := make(map[string]interface{})
	fn := func(key interface{},val interface{})bool{
		m[key.(string)] = val
		return true
	}
	rs.kv.Range(fn)
	data,err := json.Marshal(m)
	if err != nil{
		return err
	}
	_,err = rs.conn.Do("SET",rs.id,string(data))
	if err != nil{
		return fmt.Errorf("set kv failed:%s",err.Error())
	}
	//defer rs.conn.Flush() 如果数据太大会阻塞redis
	return nil
}


func(rs *RedisSession)Close(){
	rs.conn.Close()
}



