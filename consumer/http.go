package consumer

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func newHTTPServer(nc *NSQConsumer, httpPort int) *http.Server {
	// logger := logp.L().Named("http")

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", httpPort),
		Handler:        setupRouter(),
		ReadTimeout:    6 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// log.Printf("run in %s mode, listen at :%d", runMode, port)

	return s
}

func setupRouter() *gin.Engine {
	engine := gin.New()

	pprof.Register(engine)

	// router.Use(middleware.LoggerToFile(), gin.Recovery())
	engine.Use(gin.Logger(), gin.Recovery())
	engine.Use(cors.Default())

	engine.GET("/ping", Ping)

	// 异常处理
	engine.NoMethod(NoMethod)
	engine.NoRoute(NoRoute)

	return engine
}

// Ping 欢迎语
func Ping(c *gin.Context) {
	respJSON(c, http.StatusOK, "success", gin.H{
		"theme":   "Ping",
		"content": "Pong",
	})
}

// NoRoute for no route request
func NoRoute(c *gin.Context) {
	respJSON(c, http.StatusNotFound, "success", gin.H{
		"theme":   "Not Match Any Route",
		"content": "This is a default route, maybe your request within a wrong route.",
	})
}

// NoMethod for no method request
func NoMethod(c *gin.Context) {
	respJSON(c, http.StatusNotFound, "success", gin.H{
		"theme":   "Not Match This Method",
		"content": "This is a default route, maybe your request in a wrong method.",
	})
}

// respJSON response data with json schema
func respJSON(c *gin.Context, code int, message string, content interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": message,
		"data":    content,
	})
}
