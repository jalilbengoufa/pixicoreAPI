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
	myConfigFile := config.ConfigFactory()
	controller := api.ControllerFactory(myConfigFile)

	v2Beta := r.Group("v2Beta")

	{
		v2Beta.GET("/", controller.Getlocal)
		v2Beta.GET("/boot/:macAddress", controller.BootServers)
		v2Beta.GET("/single/:macAddress", controller.InstallServer)
		v2Beta.GET("/all/", controller.InstallAll)
		v2Beta.GET("/servers", controller.GetServers)
	}

	r.Run(":3000")
}
