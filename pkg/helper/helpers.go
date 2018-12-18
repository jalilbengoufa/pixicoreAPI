package helper

import (
	"strings"

	"github.com/matishsiao/goInfo"
	log "github.com/sirupsen/logrus"
)

type PxeSpec struct {
	K   string   `json:"kernel"`
	I   []string `json:"initrd"`
	CMD string   `json:"cmdline"`
}

//// *************************** HELPER FUNCTIONS ****************************
func PixicoreInit(IPAddress string) PxeSpec {

	cmd := "coreos.autologin coreos.first_boot=1 coreos.config.url={{ URL \"file:///home/cedille/pxe-config.ign\" }}"
	ip := "ip="
	ip = strings.Join([]string{ip, IPAddress}, "")
	cmd = strings.Join([]string{cmd, ip}, " ")

	pxeSpec := PxeSpec{

		K: "file:///home/cedille/coreos_production_pxe.vmlinuz",
		I: []string{
			"file:///home/cedille/coreos_production_pxe_image.cpio.gz",
		},
		CMD: "coreos.autologin coreos.first_boot=1 coreos.config.url={{ URL \"file:///home/cedille/pxe-config.ign\" }}"}

	return pxeSpec
}
func GetServerInfo() *goInfo.GoInfoObject {

	gi := goInfo.GetInfo()
	log.Infoln(gi.String())

	return gi
}
