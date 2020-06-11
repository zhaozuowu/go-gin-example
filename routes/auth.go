package routes

import (
	"github.com/afex/hystrix-go/hystrix"
	"github.com/eapache/go-resiliency/retrier"
	"github.com/gin-gonic/gin"
	. "github.com/gin2/controllers"
	"github.com/gin2/pkg/logging"
	"net/http"
	"time"
)

type AuthRoute struct {
	auth *Auth
}
var result map[string]interface{}
func (authRoute *AuthRoute) LoadAuthRoute(e *gin.Engine) {

	e.GET("/auth", authRoute.auth.Auth)
	e.GET("/test", func(ctx *gin.Context) {


		hystrix.ConfigureCommand("get_baidu", hystrix.CommandConfig{
			Timeout:               100,
			MaxConcurrentRequests: 1800,
			ErrorPercentThreshold: 30,
			SleepWindow:           5000,
		})

		result =  map[string]interface{}{"code":200,"message":"成功","data":""}
		 hystrix.Do("get_baidu", func() error {
			r := retrier.New(retrier.ConstantBackoff(3, 100*time.Millisecond), nil)
			// retrier 工作模式和 hystrix 类似，在 Run 方法中将待执行的业务逻辑封装到匿名函数传入即可
			err := r.Run(func() error {
				logging.Info("try agin")
				start := time.Now()
				res, err := http.Get("https://www.baidu.com")
				end := time.Now()
				seconds := end.Sub(start).Seconds()
				if err != nil {
					logging.Error("http request fail", err)
					//fmt.Printf("http request fail:%v\n", err)
					return err
				}
				logging.Info("消耗的时间为:", seconds, res.Body)
				result["data"] = map[string]interface{}{"id":1,"name":"zhangsan","age":1}
				return nil
			})

			return err

		}, func(err error) error {
			//fmt.Println("get an error, handle it")
			logging.Error("get an error,handle it", err)
			result["code"] = http.StatusInternalServerError
			result["message"] = err
			return nil
		})

		ctx.JSON(200,result)
	})

}

func NewAuthRoute() *AuthRoute {

	route := &AuthRoute{}
	route.auth = NewAuth()
	return route

}
