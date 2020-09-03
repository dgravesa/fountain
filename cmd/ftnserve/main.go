package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/dgravesa/fountain/pkg/data"
	"github.com/dgravesa/fountain/pkg/data/gcp"
	"github.com/dgravesa/fountain/pkg/data/redis"
	"github.com/dgravesa/fountain/pkg/resources"
	"github.com/gin-gonic/gin"
)

func main() {
	// command line arguments
	var port uint
	flag.UintVar(&port, "Port", 8080, "host port number")
	var userStoreType string
	flag.StringVar(&userStoreType, "UserStore", "", "user store backend [datastore, redis]")
	var userStoreAddr string
	flag.StringVar(&userStoreAddr, "UserStoreAddr", "", "address of user store")
	flag.Parse()

	// initialize resources
	userStore := initializeUserStore(userStoreType, userStoreAddr)
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

func initializeUserStore(storeType, addr string) data.UserStore {
	if storeType == "" {
		// use default
		return data.DefaultUserStore()
	} else if storeType == "datastore" {
		return gcp.DatastoreClient{}
	} else if storeType == "redis" {
		store, err := redis.NewUserStore(addr)
		if err != nil {
			log.Fatalln("error on initializing redis client:", err)
		}
		return store
	}

	log.Fatalln("invalid user store type specified:", storeType)
	return nil
}
