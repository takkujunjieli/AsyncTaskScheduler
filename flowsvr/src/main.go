package main

import (
	"fmt"

	"github.com/niuniumart/gosdk/gin"
	"github.com/takkujunjieli/AsyncTaskScheduler/flowsvr/src/config"
	"github.com/takkujunjieli/AsyncTaskScheduler/flowsvr/src/initialize"
	"github.com/takkujunjieli/AsyncTaskScheduler/flowsvr/src/rtm"

	"github.com/niuniumart/gosdk/martlog"
)

func main() {
	// Initialize configuration settings
	config.Init()

	// Initialize necessary resources, mainly the MySQL connection
	err := initialize.InitResource()
	if err != nil {
		// Print error to console and log the error using martlog
		fmt.Printf("initialize.InitResource error: %s", err.Error())
		martlog.Errorf("initialize.InitResource error: %s", err.Error())
		return
	}

	// Start task management runtime
	var rtm rtm.TaskRuntime
	rtm.Run()

	// Create a web server using the Gin framework
	router := gin.CreateGin()

	// Register API routes with the router, defining the available endpoints
	initialize.RegisterRouter(router)
	fmt.Println("Before starting the router")

	// Start the web server and block the main thread here while allowing requests
	// to be handled concurrently by Gin's child goroutines.
	err = gin.RunByPort(router, config.Conf.Common.Port)

	// Print any errors encountered when starting the server
	fmt.Println(err)
}
