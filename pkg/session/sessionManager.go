/**
*FileName: session
*Create on 2018/11/7 下午8:42
*Create by mok
*/

package session

import (
	"sync"
	"github.com/satori/go.uuid"
	"errors"
	"github.com/gomodule/redigo/redis"
	"time"
	"shopadmin/config"
)

var(
	NewUUIDErr = errors.New("new sessionID failed ")
	SessionNotExists = errors.New("cannot found the session")
)

var Mgr SessionManager

type SessionManager interface {
	CreateSession()(Session,error)
	GetSession(idString string)(Session,error)
	DeleteSession(idString string)(error)
}


func Init(){
	switch config.SessionProvider {
	case "memory":
		Mgr = newMemorySessionManager()
	case "redis":
		Mgr= newRedisSessionManager()
	default:
		Mgr = newMemorySessionManager()
	}
}


type MemorySessionManager struct {
	sessions map[string]*MemorySession
	rw sync.RWMutex
}

func newMemorySessionManager()(msmgr *MemorySessionManager){
	return &MemorySessionManager{
		sessions:make(map[string]*MemorySession,1000),
	}
}


func newSessionID()(string,error){
	sessionID,err := uuid.NewV4()
	if err != nil{
		return "",NewUUIDErr
	}
	idString := sessionID.String()
	return idString,nil
}

func(msmgr *MemorySessionManager)CreateSession()(session Session,err error){
	idString,err := newSessionID()
	msmgr.rw.Lock()
	defer msmgr.rw.Unlock()
	msmgr.sessions[idString] = newMemorySession(idString)
	return msmgr.sessions[idString],nil
}

func (msmgr *MemorySessionManager)GetSession(idString string)(session Session,err error){
	msmgr.rw.RLock()
	defer msmgr.rw.RUnlock()
	var ok bool
	if session,ok = msmgr.sessions[idString];!ok{
		return nil,SessionNotExists
	}
	return
}

func (msmgr *MemorySessionManager)DeleteSession(idString string)error{
	msmgr.rw.Lock()
	delete(msmgr.sessions,idString)
	msmgr.rw.Unlock()
	return nil
}


type RedisSessionManager struct {
	sessions map[string]*RedisSession
	rw sync.RWMutex
	pool *redis.Pool
}

func newpool()*redis.Pool{
	pool := &redis.Pool{
		Dial: dial,
		MaxIdle:100,
		MaxActive:500,
		IdleTimeout:240*time.Second,
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}

	return pool
}

func dial()(redis.Conn,error){
	c,err := redis.Dial("tcp","")
	if err != nil{
		return nil,err
	}
	return c,err
}


func newRedisSessionManager()(rsmgr *RedisSessionManager){
	rsmgr = &RedisSessionManager{
		sessions:make(map[string]*RedisSession,10000),
		pool:newpool(),
	}
	return
}

func(rsmgr *RedisSessionManager)CreateSession()(Session,error){
	rsmgr.rw.Lock()
	defer rsmgr.rw.Unlock()
	idString,err := newSessionID()
	if err !=nil{
		return nil,err
	}
	session := newRedisSession(idString,rsmgr.pool)
	rsmgr.sessions[idString] = session
	return session,nil
}

func(rsmgr *RedisSessionManager)GetSession(idString string)(Session,error){
	rsmgr.rw.RLock()
	defer rsmgr.rw.RUnlock()
	session,ok := rsmgr.sessions[idString]
	if !ok {
		return nil,SessionNotExists
	}
	return session,nil
}

func(rsmgr *RedisSessionManager)DeleteSession(idString string)(error){
	rsmgr.rw.Lock()
	delete(rsmgr.sessions,idString)
	rsmgr.rw.Unlock()
	return nil
}