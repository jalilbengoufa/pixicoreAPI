package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jalilbengoufa/pixicoreAPI/pkg/api"
	"github.com/jalilbengoufa/pixicoreAPI/pkg/config"
	"github.com/jalilbengoufa/pixicoreAPI/pkg/helper"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(helper.Cors())
	myConfigFile := config.InitConfig()
	controller := api.InitController(myConfigFile)

	v1 := r.Group("v2Beta")

	{
		v1.GET("/", controller.Getlocal)
		v1.GET("/boot/:macAddress", controller.BootServers)
		v1.GET("/single/:macAddress", controller.InstallServer)
		v1.GET("/all/", controller.InstallAll)
		v1.GET("/servers", controller.GetServers)
	}

	r.Run(":3000")
}
