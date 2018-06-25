package main

import (
	// This Service
	"bookinfo/bookdetails-service/svc/server"
	"bookinfo/bookdetails-service/global"
)

func main() {

	global.GenPid(global.ProjectPath + "/runtime/pid")
	server.Run()
}
