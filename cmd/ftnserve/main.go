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

	// initialize users resource
	userStore := data.DefaultUserStore()
	usersResource := resources.NewUsersResource(userStore)
	reservoir := data.DefaultReservoir()
	waterlogsResource := resources.NewWaterLogsResource(reservoir)

	// initialize routes
	// TODO: middleware for user param verification
	r := gin.Default()
	r.GET("/users/:id", usersResource.GetUser)
	r.POST("/users", usersResource.PostUser)
	r.GET("/users/:id/waterlogs", waterlogsResource.GetWls)
	r.POST("/users/:id/waterlogs", waterlogsResource.PostWl)

	// listen and serve
	portStr := fmt.Sprintf(":%d", port)
	if err := r.Run(portStr); err != nil {
		log.Fatalln(err)
	}
}
