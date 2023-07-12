package main

import (
	"context"
	"fmt"
	"net/http"
	"oom_demo/dao"
	"oom_demo/rate_limit"
	"time"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func main() {
	routersInit := Router()
	endPoint := fmt.Sprintf(":%d", 8080)
	maxHeaderBytes := 1 << 20
	rate_limit.NewTokenBucket("test", 100, 100)
	server := &http.Server{
		Addr:           endPoint,
		Handler:        routersInit,
		MaxHeaderBytes: maxHeaderBytes,
	}

	server.ListenAndServe()

}

func Router() *gin.Engine {
	// gin framework
	router := gin.Default()
	v1 := router.Group("v1",RateLimit())
	// 定义接口
	v1.GET("/test", NewEsClient)
	pprof.Register(router)
	return router
}

func AppendContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(context.Background(), "key", make([]int64, 0, 1024))
		c.Set("ctx", ctx)
		c.Next()
	}
}

func Test(ctx *gin.Context) {
	time.Sleep(200 * time.Millisecond)
	ctx.JSON(http.StatusOK, nil)
}

func NewEsClient(ctx *gin.Context) {
	dao.NewEs()
	ctx.JSON(http.StatusOK, nil)
}

func RateLimit() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 从context中获取令牌桶
		tokenBucket := rate_limit.GetTokenBucket("test")
		// 如果取不到令牌，直接返回响应
		if !tokenBucket.Take() {
			ctx.JSON(http.StatusOK, gin.H{
				"code": 200})
		}

	}
}
