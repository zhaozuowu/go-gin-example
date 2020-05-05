package redis

import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"github.com/gin2/pkg/logging"
	"github.com/gin2/pkg/setting"
	"time"
)

var RedisConn *redis.Pool

func init() {
	RedisConn = &redis.Pool{
		MaxIdle:     setting.RedisMaxIdle,
		MaxActive:   setting.RedisMaxActive,
		IdleTimeout: setting.RedisIdleTimeout,
		Dial: func() (conn redis.Conn, err error) {
			logging.Info("redisHost:",setting.RedisHost)
			c, err := redis.Dial("tcp", setting.RedisHost)
			if err != nil {
				return nil, err
			}
			if setting.RedisPassword != "" {
				if _, err := c.Do("AUTH", setting.RedisPassword); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func Set(key string, data interface{}, time int) error {

	con := RedisConn.Get()
	defer con.Close()
	value, err := json.Marshal(data)

	if err != nil {
		return err
	}
	_, err = con.Do("SET", key, value,"EX",time)

	if err != nil {
		return err
	}
	return nil

}

func ExistsKey(key string) bool {
	con := RedisConn.Get()
	defer con.Close()
	exists, err := redis.Bool(con.Do("EXISTS", key))

	if err != nil {
		return false
	}
	return exists
}

func Get(key string) ([]byte, error) {

	con := RedisConn.Get()
	defer con.Close()
	reply, err := redis.Bytes(con.Do("GET", key))
	//logging.Info("reply:",reply)
	if err != nil {
		return nil, err
	}
	return reply, nil

}

func Delete(key string) (bool, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	return redis.Bool(conn.Do("DEL", key))
}
