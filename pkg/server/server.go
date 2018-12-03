package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jalilbengoufa/pixicoreAPI/pkg/helper"
	log "github.com/sirupsen/logrus"
)

type Servers map[string]Server

//SERVERS represent a config of a server
type Server struct {
	MacAddress       string `yaml:"macAddress"`
	IPAddress        string `yaml:"ipAddress"`
	Installed        bool   `yaml:"installed"`
	Kernel           string `yaml:"kernel"`
	SecondMacAddress string `yaml:"secondmacAddress"`
}

func (servers Servers) addServer(macAddress string) {
	server := Server{
		MacAddress: macAddress, IPAddress: "change me", Installed: false, Kernel: "linux", SecondMacAddress: "find me"}
	servers[macAddress] = server

}

func (servers Servers) IsExist(c *gin.Context) bool {
	macAddr := c.Param("macAddress")
	if _, ok := servers[macAddr]; ok {
		return true
	}

	err := fmt.Sprint("This Requested server doesn't exist : ", macAddr)
	log.Warningln(err)

	return false
}

func (servers Servers) Boot(c *gin.Context) {
	macAddr := c.Param("macAddress")
	if servers.IsExist(c) {
		servers.addServer(macAddr)
	}

	helper.PixicoreInit(c)

}
