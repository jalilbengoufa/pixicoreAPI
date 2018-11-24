package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/ssh"
)

//// *************************** Controllers ****************************

//Getlocal pixicore demands
func Getlocal(c *gin.Context) {
	c.JSON(200, "success")
}

//BootServers called by pixicore client to register a new server
func BootServers(c *gin.Context) {
	if serverExist(c.Param("macAddress"), c) {
		createServer(c.Param("macAddress"), c, ReadConfig())
	}
	pixicoreInit(c.Param("macAddress"), c)
}

//InstallServer Install a single server
func InstallServer(c *gin.Context) {

	var servers map[string]Servers
	servers = ReadConfig()
	servers[c.Param("macAddress")] = runCommandsInServers(c, servers[c.Param("macAddress")])

	s, err := yaml.Marshal(&servers)
	f, err := os.Create("servers-config.yaml")
	check(err)
	f.Write(s)
	f.Sync()
	f.Close()

	c.JSON(200, servers[c.Param("macAddress")])
}

//InstallAll install all the servers available
func InstallAll(c *gin.Context) {
	var servers map[string]Servers
	servers = ReadConfig()

	for k := range servers {
		servers[k] = runCommandsInServers(c, servers[k])
	}

	s, err := yaml.Marshal(&servers)
	f, err := os.Create("servers-config.yaml")
	check(err)
	f.Write(s)
	f.Sync()
	f.Close()

	c.JSON(200, servers)
}

//CollectServerInfo collect information about a server with ssh
func runCommandsInServers(c *gin.Context, server Servers) Servers {

	sshConfig := &ssh.ClientConfig{
		User: "core",
		Auth: []ssh.AuthMethod{
			PublicKeyFile("/home/cedille/.ssh/id_rsa"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	clientSSH := &SSHClient{
		Config: sshConfig,
		Host:   server.IPAddress,
		Port:   22,
	}

	// run command with ssh
	kernel, err := clientSSH.RunCommand("uname -r")
	if err != nil {
		fmt.Fprintf(os.Stderr, "command run error: %s\n", err)
	}

	macAddressFirst, err := clientSSH.RunCommand("cat /sys/class/net/enp4s0/address")
	if err != nil {
		fmt.Fprintf(os.Stderr, "command run error: %s\n", err)
	}
	macAddressSecond, err := clientSSH.RunCommand("cat /sys/class/net/enp5s0/address")
	if err != nil {
		fmt.Fprintf(os.Stderr, "command run error: %s\n", err)
	}

	if server.MacAddress == strings.TrimSuffix(macAddressFirst, "\r\n") {
		server.SecondMacAddress = strings.TrimSuffix(macAddressSecond, "\r\n")
	} else {
		server.SecondMacAddress = strings.TrimSuffix(macAddressFirst, "\r\n")
	}
	server.Kernel = strings.TrimSuffix(kernel, "\r\n")

	installDistro, err := clientSSH.RunCommand("sudo coreos-install -d /dev/sda -i /run/ignition.json -C stable")
	_ = installDistro
	server.Installed = true

	return server
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
func GetServers(c *gin.Context) {
	c.JSON(200, gin.H{"success": ReadConfig()})
}
