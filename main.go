package main

import (
	"fmt"
	"io"
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
	MacAddress string `yaml:"macAddress"`
	IPAddress  string `yaml:"ipAddress"`
	Installed  bool   `yaml:"installed"`
}
type SSHCommand struct {
	Path   string
	Env    []string
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

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
		//v1.GET("/install/:macAddress", InstallServer)
		//v1.GET("/install/all", InstallAll)
		v1.GET("/servers", GetServers)
		v1.GET("/collect", GetCollect)

	}

	r.Run(":3000")
}

/**
	Les fonctions pour les routes
**/
func Getlocal(c *gin.Context) {
	c.JSON(200, "success")
}
func UpdateServer()                {}
func InstallServer(c *gin.Context) {}
func InstallAll(c *gin.Context)    {}
func GetCollect(c *gin.Context) {

	sshConfig := &ssh.ClientConfig{
		User: "cedille",
		Auth: []ssh.AuthMethod{
			ssh.Password("Cedille123")},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client := &SSHClient{
		Config: sshConfig,
		Host:   "192.168.0.102",
		Port:   22,
	}

	cmd := &SSHCommand{
		Path:   "ls /home/cedille/",
		Env:    []string{"LC_DIR=/"},
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	fmt.Printf("Running command: %s\n", cmd.Path)
	if err := client.RunCommand(cmd); err != nil {
		fmt.Fprintf(os.Stderr, "command run error: %s\n", err)
		os.Exit(1)
	}
}

func (client *SSHClient) RunCommand(cmd *SSHCommand) error {
	var (
		session *ssh.Session
		err     error
	)

	if session, err = client.newSession(); err != nil {
		return err
	}
	defer session.Close()

	if err = client.prepareCommand(session, cmd); err != nil {
		return err
	}

	err = session.Run(cmd.Path)
	return err
}

func (client *SSHClient) prepareCommand(session *ssh.Session, cmd *SSHCommand) error {
	for _, env := range cmd.Env {
		variable := strings.Split(env, "=")
		if len(variable) != 2 {
			continue
		}

		if err := session.Setenv(variable[0], variable[1]); err != nil {
			return err
		}
	}

	if cmd.Stdin != nil {
		stdin, err := session.StdinPipe()
		if err != nil {
			return fmt.Errorf("Unable to setup stdin for session: %v", err)
		}
		go io.Copy(stdin, cmd.Stdin)
	}

	if cmd.Stdout != nil {
		stdout, err := session.StdoutPipe()
		if err != nil {
			return fmt.Errorf("Unable to setup stdout for session: %v", err)
		}
		go io.Copy(cmd.Stdout, stdout)
	}

	if cmd.Stderr != nil {
		stderr, err := session.StderrPipe()
		if err != nil {
			return fmt.Errorf("Unable to setup stderr for session: %v", err)
		}
		go io.Copy(cmd.Stderr, stderr)
	}

	return nil
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

//called by pixicore client
func BootServers(c *gin.Context) {
	if serverExist(c.Param("macAddress"), c) {
		createServer(c.Param("macAddress"), c, ReadConfig())
		pixicoreInit(c.Param("macAddress"), c)
	} else {
		c.JSON(400, gin.H{"success": "serveur exist deja"})
	}
}
func GetServers(c *gin.Context) {
	filename, _ := filepath.Abs("servers-config.yaml")
	yamlFile, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	servers := make(map[string]Servers)
	err = yaml.Unmarshal(yamlFile, &servers)
	if err != nil {
		panic(err)
	}
	c.JSON(200, gin.H{"success": servers})
}

/*
Helper functions
*/
func serverExist(addr string, c *gin.Context) bool {
	filename, _ := filepath.Abs("servers-config.yaml")
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	servers := make(map[string]Servers)
	err = yaml.Unmarshal(yamlFile, &servers)
	if err != nil {
		panic(err)
	}
	if _, ok := servers[addr]; !ok {
		return true
	}
	return false
}
func ReadConfig() map[string]Servers {

	filename, _ := filepath.Abs("servers-config.yaml")
	yamlFile, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
	}
	server := make(map[string]Servers)
	err = yaml.Unmarshal(yamlFile, &server)
	if err != nil {
		panic(err)
	}
	return server

}
func createServer(macAddress string, c *gin.Context, servers map[string]Servers) {
	server := Servers{macAddress, "change me", false}
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

	c.JSON(200, gin.H{"success": resp})
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

//create config if it does not exist
func InitConfig() {
	if _, err := os.Stat("servers-config.yaml"); os.IsNotExist(err) {
		f, err := os.Create("servers-config.yaml")
		check(err)
		f.Close()
	}
}
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.Next()
	}
}
