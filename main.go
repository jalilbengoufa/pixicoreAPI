package main

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/ssh"
)

//Servers represent a config of a server
type Servers struct {
	MacAddress       string `yaml:"macAddress"`
	IPAddress        string `yaml:"ipAddress"`
	Installed        bool   `yaml:"installed"`
	Kernel           string `yaml:"kernel"`
	SecondMacAddress string `yaml:"secondmacAddress"`
}

//SSHClient used for ssh client
type SSHClient struct {
	Config *ssh.ClientConfig
	Host   string
	Port   int
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(Cors())
	InitConfig()
	v1 := r.Group("v1")
	{
		v1.GET("/", Getlocal)
		v1.GET("/boot/:macAddress", BootServers)
		v1.GET("/single/:macAddress", InstallServer)
		v1.GET("/all/", InstallAll)
		v1.GET("/servers", GetServers)
	}

	r.Run(":3000")
}
