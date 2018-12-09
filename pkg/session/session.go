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
	kv map[string]interface{}
	rw sync.RWMutex
}


func newMemorySession(id string)(ms *MemorySession){
	return &MemorySession{
		id:id,
		kv:make(map[string]interface{},5),
	}
}


func (ms *MemorySession)Set(key string,val interface{}){
	ms.rw.Lock()
	defer ms.rw.Unlock()
	ms.kv[key] = val
}

func (ms *MemorySession)Get(key string)(val interface{},err error){
	ms.rw.RLock()
	var ok bool
	if val,ok = ms.kv[key];ok{
		return
	}
	ms.rw.RUnlock()
	err = KeyNotExists
	return
}

func (ms *MemorySession)MustGet(key string)(interface{}){
	ms.rw.RLock()
	defer ms.rw.RUnlock()
	return ms.kv[key]
}
func (ms *MemorySession)GetAll()(map[string]interface{},error){
	return ms.kv,nil
}

func  (ms *MemorySession)ID()string{
	return ms.id
}

func(ms *MemorySession)Delete(key string){
	ms.rw.Lock()
	defer ms.rw.Unlock()
	delete(ms.kv,key)
}

func (ms *MemorySession)Save()error{
	return nil
}

func (ms *MemorySession)Close(){
	return
}



type RedisSession struct {
	id string
	kv map[string]interface{}
	conn redis.Conn
	rw sync.RWMutex
	flag int
}


func newRedisSession(idString string,pool *redis.Pool)(*RedisSession){
	rs :=  &RedisSession{
		id:idString,
		kv:make(map[string]interface{},10),
		conn:pool.Get(),
		flag:SessionFlagNone,
	}
	return rs
}

func(rs *RedisSession)Set(key string,val interface{}){
	rs.rw.Lock()
	defer rs.rw.Unlock()
	rs.kv[key] = val
	rs.flag = SessionFlagModify
}


func(rs * RedisSession)readFromRedis()error{
	reply,err := rs.conn.Do("GET",rs.id)
	data,err := redis.String(reply,err)
	if err != nil{
		return err
	}
	err = json.Unmarshal([]byte(data),&rs.kv)
	if err != nil{
		return err
	}
	return nil
}


func(rs *RedisSession)Get(key string)(val interface{},err error){
	rs.rw.RLock()
	defer rs.rw.RUnlock()
	if rs.flag == SessionFlagNone{
		err = rs.readFromRedis()
		if err != nil{
			return nil,err
		}
	}
	val = rs.kv[key]
	return val,nil
}


func(rs *RedisSession)MustGet(key string)(interface{}){
	val,err := rs.Get(key)
	if err != nil{
		return nil
	}
	return val
}


func (rs *RedisSession)GetAll()(map[string]interface{},error){
	rs.rw.RLock()
	defer rs.rw.RUnlock()
	if rs.flag == SessionFlagNone{
		err := rs.readFromRedis()
		if err != nil{
			return nil,err
		}
	}
	all := rs.kv
	return all,nil

}


func (rs *RedisSession)ID()string{
	return rs.id
}

func(rs *RedisSession)Delete(key string){
	rs.rw.Lock()
	defer rs.rw.Unlock()
	delete(rs.kv,key)
}


func (rs *RedisSession)Save()error{
	rs.rw.Lock()
	defer rs.rw.Unlock()
	data,err := json.Marshal(rs.kv)
	if err != nil{
		return err
	}
	_,err = rs.conn.Do("SET",rs.id,string(data))
	if err != nil{
		return fmt.Errorf("set kv failed:%s",err.Error())
	}
	defer rs.conn.Flush()
	return nil
}


func(rs *RedisSession)Close(){
	rs.conn.Close()
}



