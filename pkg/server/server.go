package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jalilbengoufa/pixicoreAPI/pkg/helper"
	log "github.com/sirupsen/logrus"
	"net/http"
)

//Servers represent a list of Server type
type Servers map[string]Server

//Server represent a config of a server
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

//IsExist verify if server exist on list of server using gin context as input
func (servers Servers) IsExist(c *gin.Context) bool {
	macAddr := c.Param("macAddress")
	if _, ok := servers[macAddr]; ok {
		return true
	}

	err := fmt.Sprint("This Requested server doesn't exist : ", macAddr)
	log.Warningln(err)

	return false
}

//GetServer Get server from a list of servers using context as input
func (servers Servers) GetServer(c *gin.Context) Server {
	macAddr := c.Param("macAddress")
	if !servers.IsExist(c) {
		c.JSON(http.StatusNotFound, gin.H{"status": "server don't exist"})
	}

	server := servers[macAddr]

	return server
}

//Boot Boot server specified in gin Context
func (server Server) Boot(c *gin.Context) {
	macAddr := c.Param("macAddress")
	pxeSpec := helper.PixicoreInit(macAddr)
	c.JSON(200, pxeSpec)

}
