package redis

import (
	"bytes"
	"compress/flate"
	"github.com/garyburd/redigo/redis"
	"github.com/gin2/pkg/logging"
	"github.com/gin2/pkg/setting"
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
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
				logging.Error("redisHost connect error:",err)
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

// 编码
func Gzdeflate(data []byte,level int) []byte  {

	var bufs bytes.Buffer
	w,_ :=flate.NewWriter(&bufs,level)
	w.Write(data)
	w.Flush()
	defer w.Close()
	return bufs.Bytes()
}

func Set(key string, data interface{}, time int) error {
	con := RedisConn.Get()
	defer con.Close()
	//value, err := jsoniter.Marshal(data)
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	sValue,err := json.Marshal(data)
	//sValue ,err := serialize.Marshal(data)


	if err != nil {
		return err
	}

	value := Gzdeflate(sValue,-1)
	_, err = con.Do("SET", key, value,"EX",time)

	if err != nil {
		return err
	}
	return nil

}


// 解码
func Gzdecode(data []byte) []byte  {

	r :=flate.NewReader(bytes.NewReader(data))
	defer r.Close()
	out, err := ioutil.ReadAll(r)
	if err !=nil {
		//fmt.Errorf("%s\n",err)
		//return []byte{}
	}
	return out
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
