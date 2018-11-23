package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/gin-gonic/gin"
)

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
