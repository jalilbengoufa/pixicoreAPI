package helper

import (
	"github.com/gin-gonic/gin"
	"strings"
)

//// *************************** HELPER FUNCTIONS ****************************
func PixicoreInit(c *gin.Context) {
	IPAddress := c.Param("IPAddress")
	cmd := "coreos.autologin coreos.first_boot=1 coreos.config.url={{ URL \"file:///home/cedille/pxe-config.ign\" }}"
	ip := "ip="
	ip = strings.Join([]string{ip, IPAddress}, "")
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

//Cors cors for the api
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.Next()
	}
}
