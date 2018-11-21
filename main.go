package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/ghodss/yaml"
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

	r := gin.Default()
	r.Use(Cors())
	InitConfig()
	v1 := r.Group("v1")
	{
		v1.GET("/", Getlocal)
		v1.GET("/boot/:macAddress", BootServers)
		v1.GET("/install/:macAddress", InstallServer)
		//v1.GET("/install/all", InstallAll)
		v1.GET("/servers", GetServers)
		v1.GET("/collect", GetCollect)

	}

	r.Run(":3000")
}

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
	servers[c.Param("macAddress")] = CollectServerInfo(servers[c.Param("macAddress")])

	s, err := yaml.Marshal(&servers)
	f, err := os.Create("servers-config.yaml")
	check(err)
	f.Write(s)
	f.Sync()
	f.Close()

}

//InstallAll install all the servers available
func InstallAll(c *gin.Context) {}
func GetCollect(c *gin.Context) {

	/*var (
		session *ssh.Session
		err     error
		client  *SSHClient
	)

	if session, err = client.newSession(); err != nil {
		fmt.Print("error while creating sesion")
	}

	defer session.Close()

	sshConfig := &ssh.ClientConfig{
		User: "core",
		Auth: []ssh.AuthMethod{
			PublicKeyFile("/home/django/src/ssh1"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	clientSSH := &SSHClient{
		Config: sshConfig,
		Host:   "192.168.0.105",
		Port:   22,
	}

	kernel, err := clientSSH.RunCommand("uname -r", session)
	if err != nil {
		fmt.Fprintf(os.Stderr, "command run error: %s\n", err)
	}

	macAddressFirst, err := clientSSH.RunCommand("cat /sys/class/net/enp4s0/address", session)
	if err != nil {
		fmt.Fprintf(os.Stderr, "command run error: %s\n", err)
	}
	macAddressSecond, err := clientSSH.RunCommand("cat /sys/class/net/enp5s0/address", session)
	if err != nil {
		fmt.Fprintf(os.Stderr, "command run error: %s\n", err)
	}*/
}

//CollectServerInfo collect information about a server with ssh
func CollectServerInfo(server Servers) Servers {

	var (
		session *ssh.Session
		err     error
		client  *SSHClient
	)

	if session, err = client.newSession(); err != nil {
		fmt.Print("error while creating sesion")
	}

	defer session.Close()

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

	kernel, err := clientSSH.RunCommand("uname -r", session)
	if err != nil {
		fmt.Fprintf(os.Stderr, "command run error: %s\n", err)
	}
	server.Kernel = kernel

	macAddressFirst, err := clientSSH.RunCommand("cat /sys/class/net/enp4s0/address", session)
	if err != nil {
		fmt.Fprintf(os.Stderr, "command run error: %s\n", err)
	}

	macAddressSecond, err := clientSSH.RunCommand("cat /sys/class/net/enp5s0/address", session)
	if err != nil {
		fmt.Fprintf(os.Stderr, "command run error: %s\n", err)
	}
	if server.MacAddress == macAddressFirst {
		server.SecondMacAddress = macAddressSecond
	} else {
		server.SecondMacAddress = macAddressFirst
	}

	return server
}

//RunCommand run ssh command in the remote server and retrun output
func (client *SSHClient) RunCommand(command string, session *ssh.Session) (string, error) {

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

//// *************************** HELPER FUNCTIONS ****************************

func serverExist(addr string, c *gin.Context) bool {
	filename, _ := filepath.Abs("servers-config.yaml")
	yamlFile, err := ioutil.ReadFile(filename)
	check(err)
	servers := make(map[string]Servers)
	err = yaml.Unmarshal(yamlFile, &servers)
	check(err)
	if _, ok := servers[addr]; !ok {
		return true
	}
	return false
}

//ReadConfig read the yaml config file
func ReadConfig() map[string]Servers {

	filename, _ := filepath.Abs("servers-config.yaml")
	yamlFile, err := ioutil.ReadFile(filename)
	check(err)

	server := make(map[string]Servers)
	err = yaml.Unmarshal(yamlFile, &server)
	check(err)

	return server
}
func createServer(macAddress string, c *gin.Context, servers map[string]Servers) {
	server := Servers{macAddress, "change me", false, "linux", "find me"}
	servers[macAddress] = server

	s, err := yaml.Marshal(&servers)
	f, err := os.Create("servers-config.yaml")
	check(err)
	f.Write(s)
	f.Sync()
	f.Close()
}

func pixicoreInit(ipAddress string, c *gin.Context) {

	cmd := "coreos.autologin coreos.first_boot=1 coreos.config.url={{ URL \"file:///home/cedille/pxe-config.ign\" }}"
	ip := "ip="
	ip = strings.Join([]string{ip, ipAddress}, "")
	cmd = strings.Join([]string{cmd, ip}, " ")

	resp := struct {
		K   string   `json:"kernel"`
		I   []string `json:"initrd"`
		CMD string   `json:"cmdline"`
	}{
		K: "file:///home/cedille/coreos_production_pxe.vmlinuz",
		I: []string{
			"file:///home/cedille/coreos_production_pxe_image.cpio.gz",
		},
		CMD: "coreos.autologin coreos.first_boot=1 coreos.config.url={{ URL \"file:///home/cedille/pxe-config.ign\" }}",
	}

	c.JSON(200, resp)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

//InitConfig create config if it does not exist
func InitConfig() {
	if _, err := os.Stat("servers-config.yaml"); os.IsNotExist(err) {
		f, err := os.Create("servers-config.yaml")
		check(err)
		f.Close()
	}
}

//Cors cors for the api
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.Next()
	}
}
