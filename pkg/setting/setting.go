package setting

import (
	"github.com/go-ini/ini"
	"log"
	"time"
)

var (
	Cfg          *ini.File
	HttpPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	RunMode      string
	PageSize     int
	JwtSecret    string
	RedisHost    string
	RedisPassword   string
	RedisMaxIdle     int
	RedisMaxActive   int
	RedisIdleTimeout  time.Duration
)

func init() {

	var err error
	Cfg, err = ini.Load("conf/app.ini")
	if err != nil {
		log.Fatalf("Fail to parse 'conf/app.ini':%v\n", err)
	}

	LoadBase()
	LoadServer()
	LoadApp()
	LoadRedis()

}

func LoadRedis() {

	sec, err := Cfg.GetSection("redis")

	if err != nil {
		log.Fatalf("Fail to get section  'redis':%v\n", err)
	}

	PageSize = sec.Key("PAGE_SIZE").MustInt(1000)
	JwtSecret = sec.Key("JWT_SECRET").MustString("23347$040412")
}

func LoadApp() {

	sec, err := Cfg.GetSection("redis")

	if err != nil {
		log.Fatalf("Fail to get section  'app':%v\n", err)
	}

	RedisHost = sec.Key("Host").MustString("127.0.0.1:6379")
	RedisPassword = sec.Key("Password").MustString("")
	RedisMaxIdle = sec.Key("MaxIdle").MustInt(30)
	RedisMaxActive = sec.Key("MaxActive").MustInt(60)
	RedisIdleTimeout = time.Duration(sec.Key("IdleTimeout").MustInt(120)) * time.Second
}

func LoadServer() {
	sec, err := Cfg.GetSection("server")

	if err != nil {
		log.Fatalf("Fail to get section 'server':%v\n", err)
	}

	HttpPort = sec.Key("HTTP_PORT").MustInt(8080)
	ReadTimeout = time.Duration(sec.Key("READ_TIMEOUT").MustInt(60)) * time.Second
	WriteTimeout = time.Duration(sec.Key("WRITE_TIMEOUT").MustInt(60)) * time.Second
}

func LoadBase() {
	RunMode = Cfg.Section("").Key("RUN_MODE").MustString("debug")
}
