package main

import (
	"fmt"
	"net/http"

	"go-gin-example/pkg/setting"
	"go-gin-example/routers"
)

func main() {
	
	// router := gin.Default()
	// router.GET("/test", func(c *gin.Context) {
		// 	c.JSON(200, gin.H{
			// 		"message": "test",
			// 	})
		// })
			
	router := routers.InitRouter()
	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.HTTPPort),
		Handler:        router,
		ReadTimeout:    setting.ReadTimeout,
		WriteTimeout:   setting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()

}
