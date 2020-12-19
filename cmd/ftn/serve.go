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

type serveCmd struct {
	port          uint
	userStoreType string
	userStoreLoc  string
	reservoirType string
	reservoirLoc  string
}

func (s *serveCmd) SetFlags(flags *flag.FlagSet) {
	flags.UintVar(&s.port, "Port", 8080, "host port number")
	flags.StringVar(&s.userStoreType, "UserStore", "", "user store backend [datastore, redis]")
	flags.StringVar(&s.userStoreLoc, "UserStoreLoc", "",
		"location of user store, such as host name or project name")
	flags.StringVar(&s.reservoirType, "Reservoir", "", "reservoir backend [datastore, redis]")
	flags.StringVar(&s.reservoirLoc, "ReservoirLoc", "",
		"location of reservoir, such as host name or project name")
}

func (s *serveCmd) Exec(args []string) error {
	// initialize resources
	userStore := initializeUserStore(s.userStoreType, s.userStoreLoc)
	usersResource := resources.NewUsersResource(userStore)
	reservoir := initializeReservoir(s.reservoirType, s.reservoirLoc)
	waterlogsResource := resources.NewWaterLogsResource(reservoir)

	// initialize routes
	r := gin.Default()
	r.GET("/users/:user", usersResource.GetUser)
	r.POST("/users", usersResource.PostUser)
	r.GET("/users/:user/waterlogs", usersResource.UserMustExist, waterlogsResource.GetWls)
	r.POST("/users/:user/waterlogs", usersResource.UserMustExist, waterlogsResource.PostWl)

	// listen and serve
	portStr := fmt.Sprintf(":%d", s.port)
	return r.Run(portStr)
}

func initializeUserStore(storeType, location string) data.UserStore {
	var s data.UserStore
	var err error

	switch storeType {
	case "":
		s, err = data.DefaultUserStore()
	case "datastore":
		s = gcp.NewDatastoreClient(location)
	case "redis":
		s, err = redis.NewUserStore(location)
	default:
		err = fmt.Errorf("unknown user store client type: %s", storeType)
	}

	if err != nil {
		log.Fatalln("error on initializing user store client:", err)
	}

	return s
}

func initializeReservoir(reservoirType, location string) data.Reservoir {
	var r data.Reservoir
	var err error

	switch reservoirType {
	case "":
		r, err = data.DefaultReservoir()
	case "datastore":
		r = gcp.NewDatastoreClient(location)
	case "redis":
		r, err = redis.NewReservoir(location)
	default:
		err = fmt.Errorf("unknown reservoir client type: %s", reservoirType)
	}

	if err != nil {
		log.Fatalln("error on initializing reservoir client:", err)
	}

	return r
}
