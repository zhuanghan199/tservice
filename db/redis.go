package db

import (
	"errors"
	"time"

	"github.com/gomodule/redigo/redis"
	"tservice/config"
	"tservice/common/logger"
)

const EXPIRED_SESSION = 2 * 60 * 60 // 2h

var (
	redisPool *redis.Pool // 连接池
	redisAddr string      // 连接地址
)

type RedisConn struct {
	redis.Conn
}

func InitRedis() {
	redisAddr = config.Cnf.RedisAddr
	redisPool = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 10 * time.Second,
		MaxActive:   1024,
		Wait:        false,
		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			_, err := conn.Do("PING")
			return err
		},
		Dial: func() (redis.Conn, error) {
			conn, err := redis.DialTimeout("tcp", redisAddr, time.Second*5, time.Second*10, time.Second*10)
			// conn, err := redis.Dial("tcp", ":6379")
			if err != nil {
				logger.Errorln("Dail redis server error: ", err.Error())
				return nil, err
			}

			if _, err := conn.Do("PING"); err != nil {
				conn.Close()
				return nil, err
			}

			logger.Debugln("Active count: ", redisPool.ActiveCount())
			return conn, err
		},
	}
}

func NewRedisConn() (*RedisConn, error) {
	if redisPool == nil {
		InitRedis()
	}

	conn := redisPool.Get()
	if conn == nil {
		logger.Errorln("New RedisConn Failure")
		return nil, errors.New("New redis error")
	}

	return &RedisConn{Conn: conn}, nil
}

func (rc *RedisConn) Close() {
	rc.Conn.Flush()
	if err := rc.Conn.Close(); err != nil {
		logger.Warning("Close redis error: ", err.Error())
	}
}

func (rc *RedisConn) Publish(channel, data string) error {
	_, err := rc.Conn.Do("Publish", channel, data)
	return err
}

func (rc *RedisConn) Subscribe(channel string, callback func(msg string, err error)) {
	psc := redis.PubSubConn{Conn: rc.Conn}
	psc.Subscribe(channel)
	for {
		switch v := psc.Receive().(type) {
		case redis.Message:
			logger.Debugf("Subscribe channel[%s] message: %s", v.Channel, v.Data)
			callback(string(v.Data), nil)
		// case redis.Subscription:
		case error:
			logger.Warningf("Subscribe error: %v", v)
			callback("", v)
			return
		}
	}
}

func (rc *RedisConn) Persistence(key string, info interface{}) (err error) {
	_, err = rc.Conn.Do("SET", key, info)
	return
}

func (rc *RedisConn) Save(key string, info interface{}) (err error) {
	if _, err = rc.Conn.Do("SET", key, info); err == nil {
		_, err = rc.Conn.Do("EXPIRE", key, EXPIRED_SESSION)
	}
	return
}

func (rc *RedisConn) Increase(key string) (err error) {
	_, err = rc.Conn.Do("INCR", key)
	return
}

func (rc *RedisConn) ExpiresSave(key string, info interface{}, expires int64) (err error) {
	if _, err = rc.Conn.Do("SET", key, info); err == nil {
		_, err = rc.Conn.Do("EXPIRE", key, expires)
	}
	return
}

func (rc *RedisConn) Expire(key string, expired int64) (err error) {
	_, err = rc.Conn.Do("EXPIRE", key, expired)
	return
}

func (rc *RedisConn) Update(key string, info interface{}) (err error) {
	if _, err = rc.Conn.Do("SET", key, info); err == nil {
		// rc.Conn.Do("EXPIRE", key, EXPIRED_SESSION)
	}
	return
}

func (rc *RedisConn) Get(key string) interface{} {
	if v, err := rc.Conn.Do("GET", key); err == nil {
		return v
	} else {
		return nil
	}
}

func (rc *RedisConn) GetInt64(key string) (int64, error) {
	return redis.Int64(rc.Conn.Do("GET", key))
}

func (rc *RedisConn) Keys(pattern string) ([]interface{}, error) {
	if v, err := rc.Conn.Do("KEYS", pattern); err == nil {
		if s, ok := v.([]interface{}); ok {
			return s, nil
		} else {
			return nil, errors.New("redis keys type is not '[]interface'")
		}
	} else {
		return nil, err
	}
}

func (rc *RedisConn) Delete(key string) (err error) {
	_, err = rc.Conn.Do("DEL", key)
	return
}

func (rc *RedisConn) DBSize() (int64, error) {
	return redis.Int64(rc.Conn.Do("DBSIZE"))
}

func (rc *RedisConn) Clear() (err error) {
	_, err = rc.Conn.Do("FLUSHALL")
	return
}
