package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func main() {
	routersInit := Router()
	endPoint := fmt.Sprintf(":%d", 8080)
	maxHeaderBytes := 1 << 20

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
	v1 := router.Group("v1", AppendContext())

	// 定义接口
	v1.GET("/test", Test)
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
