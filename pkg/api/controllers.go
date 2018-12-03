package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jalilbengoufa/pixicoreAPI/pkg/config"

	"github.com/jalilbengoufa/pixicoreAPI/pkg/server"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type Controller struct {
	currentConfig config.ConfigFile
}

//// *************************** Controllers ****************************

func ControllerFactory(confFile config.ConfigFile) Controller {
	ctrl := Controller{currentConfig: confFile}
	return ctrl
}

//Getlocal pixicore demands
func (ctrl Controller) Getlocal(c *gin.Context) {
	c.JSON(200, "success")

}

//CollectServerInfo collect information about a server with ssh
func RunCommandsInServers(c *gin.Context, server server.Server) server.Server {

	sshConfig := &ssh.ClientConfig{
		User: "core",
		Auth: []ssh.AuthMethod{
			PublicKeyFile("/home/cedille/.ssh/id_rsa"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	clientSSH := SSHClient{
		Config: sshConfig,
		Host:   server.IPAddress,
		Port:   22,
	}

	// run command with ssh
	kernel, err := clientSSH.RunCommand("uname -r")
	if err != nil {
		log.Errorln(os.Stderr, "command run error: %s", err)
	}

	macAddressFirst, err := clientSSH.RunCommand("cat /sys/class/net/enp4s0/address")
	if err != nil {
		log.Errorln(os.Stderr, "command run error: %s\n", err)
	}
	macAddressSecond, err := clientSSH.RunCommand("cat /sys/class/net/enp5s0/address")
	if err != nil {
		log.Errorln(os.Stderr, "command run error: %s\n", err)
	}

	if server.MacAddress == strings.TrimSuffix(macAddressFirst, "\r\n") {
		server.SecondMacAddress = strings.TrimSuffix(macAddressSecond, "\r\n")
	} else {
		server.SecondMacAddress = strings.TrimSuffix(macAddressFirst, "\r\n")
	}
	server.Kernel = strings.TrimSuffix(kernel, "\r\n")

	_, err = clientSSH.RunCommand("sudo coreos-install -d /dev/sda -i /run/ignition.json -C stable")
	if err != nil {
		log.Errorln(os.Stderr, "command run error: %s\n", err)
	}

	server.Installed = true

	return server
}

//BootServers called by pixicore client to register a new server
func (ctrl Controller) BootServers(c *gin.Context) {
	ctrl.currentConfig.Servers.Boot(c)
}

//InstallServer Install a single server
func (ctrl Controller) InstallServer(c *gin.Context) {
	macAddr := c.Param("macAddress")
	if ctrl.currentConfig.Servers.IsExist(c) == false {

		err := fmt.Sprint("This Requested server doesn't exist : ", macAddr)
		c.JSON(http.StatusNotFound, gin.H{"status": err})
	}

	currentSvrMacAddr := ctrl.currentConfig.Servers[c.Param("macAddress")]
	ctrl.currentConfig.Servers[c.Param("macAddress")] = RunCommandsInServers(c, currentSvrMacAddr)

	ctrl.currentConfig.WriteYamlConfig()

	c.JSON(200, ctrl.currentConfig.Servers[c.Param("macAddress")])
}

//InstallAll install all the servers available
func (ctrl Controller) InstallAll(c *gin.Context) {
	servers := ctrl.currentConfig.Servers
	for k := range servers {
		servers[k] = RunCommandsInServers(c, servers[k])
	}

	ctrl.currentConfig.WriteYamlConfig()

	c.JSON(200, servers)
}

//SSHClient used for ssh client
type SSHClient struct {
	Config *ssh.ClientConfig
	Host   string
	Port   int
}

//RunCommand run ssh command in the remote server and retrun output
func (client *SSHClient) RunCommand(command string) (string, error) {
	var (
		session *ssh.Session
		err     error
	)

	if session, err = client.newSession(); err != nil {
		return "", err
	}
	defer session.Close()

	out, err := session.CombinedOutput(command)
	return string(out), err
}

func (client *SSHClient) newSession() (*ssh.Session, error) {
	connection, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", client.Host, client.Port), client.Config)
	if err != nil {
		return nil, fmt.Errorf("Failed to dial: %s", err)
	}

	session, err := connection.NewSession()
	if err != nil {
		return nil, fmt.Errorf("Failed to create session: %s", err)
	}
	modes := ssh.TerminalModes{
		// ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		session.Close()
		return nil, fmt.Errorf("request for pseudo terminal failed: %s", err)
	}

	return session, nil
}

//PublicKeyFile get public key with private key
func PublicKeyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil
	}
	return ssh.PublicKeys(key)
}

//GetServers return config of the all the servers
func (ctrl Controller) GetServers(c *gin.Context) {
	c.JSON(200, gin.H{"success": ctrl.currentConfig.Servers})
}
