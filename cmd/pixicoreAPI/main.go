package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jalilbengoufa/pixicoreAPI/pkg/api"
	"github.com/jalilbengoufa/pixicoreAPI/pkg/config"
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

	myConfigFile := config.InitConfig()
	controller := api.InitController(myConfigFile)
	r := api.GetRouter(controller)
	r.Run(":3000")
}
