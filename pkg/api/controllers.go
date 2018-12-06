package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jalilbengoufa/pixicoreAPI/pkg/config"
	"github.com/jalilbengoufa/pixicoreAPI/pkg/server"
	"github.com/jalilbengoufa/pixicoreAPI/pkg/sshclient"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"net/http"
	"strings"
)

type Controller struct {
	currentConfig config.ConfigFile
}

//// *************************** Controllers ****************************

// InitController Generate new controller for PXE boot
func InitController(confFile config.ConfigFile) Controller {
	ctrl := Controller{currentConfig: confFile}
	return ctrl
}

//Getlocal pixicore demands
func (ctrl Controller) Getlocal(c *gin.Context) {
	c.JSON(200, "success")

}

//BootServers called by pixicore client to register a new server
func (ctrl Controller) BootServer(c *gin.Context) {
	servers := ctrl.currentConfig.Servers
	server := servers.GetServer(c)
	server.Boot(c)
}

//InstallServer Install a single server
func (ctrl Controller) InstallServer(c *gin.Context) {
	macAddr := c.Param("macAddress")
	if ctrl.currentConfig.Servers.IsExist(c) == false {

		err := fmt.Sprint("This Requested server doesn't exist : ", macAddr)
		c.JSON(http.StatusNotFound, gin.H{"status": err})
	}

	ctrl.currentConfig.Servers[c.Param("macAddress")] = ctrl.CollectServerInfo(c)

	ctrl.currentConfig.WriteYamlConfig()

	c.JSON(200, ctrl.currentConfig.Servers[c.Param("macAddress")])
}

//InstallAll install all the servers available
func (ctrl Controller) InstallAll(c *gin.Context) {
	servers := ctrl.currentConfig.Servers
	for svr := range servers {
		servers[svr] = ctrl.CollectServerInfo(c)
	}

	ctrl.currentConfig.WriteYamlConfig()

	c.JSON(200, servers)
}

//CollectServerInfo collect information about a server with ssh
func (ctrl Controller) CollectServerInfo(c *gin.Context) server.Server {

	currentServer := ctrl.currentConfig.Servers[c.Param("macAddress")]
	sshConfig := ssh.ClientConfig{
		User: "core",
		Auth: []ssh.AuthMethod{
			sshclient.PublicKeyFile("/home/cedille/.ssh/id_rsa"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	clientSSH := sshclient.SSHClient{
		Config: &sshConfig,
		Host:   currentServer.IPAddress,
		Port:   22,
	}

	// run command with ssh
	kernel, err := clientSSH.RunCommand("uname -r")
	if err != nil {
		log.Errorf("command run error: %s", err)
	}

	macAddressFirst, err := clientSSH.RunCommand("cat /sys/class/net/enp4s0/address")
	if err != nil {
		log.Errorf("command run error: %s\n", err)
	}
	macAddressSecond, err := clientSSH.RunCommand("cat /sys/class/net/enp5s0/address")
	if err != nil {
		log.Errorf("command run error: %s\n", err)
	}

	if currentServer.MacAddress == strings.TrimSuffix(macAddressFirst, "\r\n") {
		currentServer.SecondMacAddress = strings.TrimSuffix(macAddressSecond, "\r\n")
	} else {
		currentServer.SecondMacAddress = strings.TrimSuffix(macAddressFirst, "\r\n")
	}
	currentServer.Kernel = strings.TrimSuffix(kernel, "\r\n")

	_, err = clientSSH.RunCommand("sudo coreos-install -d /dev/sda -i /run/ignition.json -C stable")
	if err != nil {
		log.Errorf("command run error: %s\n", err)
	}

	currentServer.Installed = true

	return currentServer

}

//GetServers return config of the all the servers
func (ctrl Controller) GetServers(c *gin.Context) {
	c.JSON(200, gin.H{"success": ctrl.currentConfig.Servers})
}
