package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/dgravesa/fountain/pkg/data"
	"github.com/dgravesa/fountain/pkg/resources"
	"github.com/gin-gonic/gin"
)

func main() {
	// command line arguments
	var port uint
	flag.UintVar(&port, "port", 8080, "host port number")
	flag.Parse()

	// initialize resources
	userStore := data.DefaultUserStore()
	usersResource := resources.NewUsersResource(userStore)
	reservoir := data.DefaultReservoir()
	waterlogsResource := resources.NewWaterLogsResource(reservoir)

	// initialize routes
	r := gin.Default()
	r.GET("/users/:user", usersResource.GetUser)
	r.POST("/users", usersResource.PostUser)
	r.GET("/users/:user/waterlogs", usersResource.UserMustExist, waterlogsResource.GetWls)
	r.POST("/users/:user/waterlogs", usersResource.UserMustExist, waterlogsResource.PostWl)

	// listen and serve
	portStr := fmt.Sprintf(":%d", port)
	if err := r.Run(portStr); err != nil {
		log.Fatalln(err)
	}
}
