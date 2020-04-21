package main

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/gin2/pkg/logging"
	"log"
	"os"
)

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	Sex  uint8  `json:"sex"`
}

func main() {

	conn, err := redis.Dial("tcp", "localhost:6379")

	if err != nil {
		logging.Error("连接数据库失败:", err)
		log.Fatalf("连接redis失败:%v\n", conn)
	}
	defer conn.Close()
	/*
		redisArr, err := redis.Strings(conn.Do("MGET","name","name2","name3"))

		if err != nil {
			logging.Info("Fail to set:name", err)
			log.Fatalf("Fail to set:name", err)
		}

		fmt.Printf("name:is:%v\n",redisArr[0])
	*/
	//conn.Send("HGET", "student","name","age")
	//conn.Send("HGET", "student","Score")
	//conn.Flush()
	//res ,err := conn.Receive()
	//res2 ,err := conn.Receive()
	//fmt.Printf("Receive res1:%v,res2:%v\n", res,res2)

	var arr []User
	arr = append(arr, User{Name: "赵作武", Age: 20, Sex: 1})
	arr = append(arr, User{Name: "赵作武2", Age: 21, Sex: 0})

	str, err := json.Marshal(arr)
	if err != nil {
		fmt.Printf("转换失败:%v\n", err)
		os.Exit(1)
	}

	_, err = conn.Do("SET", "test", str)

	if err != nil {
		fmt.Printf("写入缓存失败:%v\n", err)
		os.Exit(1)
	}

	strings, err := redis.Bytes(conn.Do("GET", "test"))

	if strings == nil || len(strings) <= 0 {
		fmt.Println("check erroor")
		return
	}

	var data []User
	err = json.Unmarshal(strings, &data)
	if err != nil {
		fmt.Printf("解析失败:%v\n",err)
	}
	fmt.Printf("THE VALUE IS:%v\n",data)

	/*
		_, err = conn.Do("set", "name", "test", "NX", "EX", 100)

		if err != nil {
			logging.Info("Fail to set:name", err)
			log.Fatalf("Fail to set:name", err)

		}
		redisValue, err := conn.Do("get", "name")
		if err != nil {
			logging.Info("Fail to get:name", err)
			log.Fatalf("Fail to get:name", err)
		}

		fmt.Printf("the redisValue is :%s",redisValue)
	*/

}
