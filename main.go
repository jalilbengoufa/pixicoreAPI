package main

import (
    "encoding/json"
	"log"
	"net/http"
    "github.com/gorilla/mux"
    "path/filepath"
    "flag"
    "strconv"
)
var servers []Server
type Server struct {
    ID           string   `json:"id,omitempty"`
    MacAddress   string   `json:"macAddress,omitempty"`
    IPAddress    string   `json:"ipAddress,omitempty"`
    DoneBoot       bool	   
}


func setStatue( addr string) {
	for _, item := range servers {
        if item.MacAddress == addr {
            item.DoneBoot = true
        }
    }
}
func exitsAndStatusAvailble( addr string) bool{
    for _, item := range servers {
        if item.MacAddress == addr {
            return !item.DoneBoot
        }
    }
    return false
}
func BootServers(w http.ResponseWriter, r *http.Request) {
	//params := mux.Vars(r);
	//if (exitsAndStatusAvailble(params["MacAddress"])){
	log.Printf("Serving boot config for %s", filepath.Base(r.URL.Path))
	resp := struct {
		K string   `json:"kernel"`
		I []string `json:"initrd"`
		CMD string `json:"cmdline"`
	}{
		K: "file:///home/cedille/coreos_production_pxe.vmlinuz",
		I: []string{
			"file:///home/cedille/coreos_production_pxe_image.cpio.gz",
		},
		CMD: "coreos.autologin coreos.first_boot=1 coreos.config.url={{ URL \"file:///home/cedille/pxe-config.ign\" }}",

	}

	json.NewEncoder(w).Encode(&resp); 
	//}
}

func GetServers(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(servers)
}
func InstallServer(w http.ResponseWriter, r *http.Request) {}
func InstallAll(w http.ResponseWriter, r *http.Request) {}
func ResetServer(w http.ResponseWriter, r *http.Request) {}

var (
	port = flag.Int("port", 3000, "Port to listen on")
)

func main() {
    flag.Parse()
    router := mux.NewRouter()
    //add servers availble
    servers = append(servers, Server{ID: "1", MacAddress: "dr-ee-6y-8o-ee", IPAddress: "168.104.0.10", DoneBoot:false})

    //routes
    router.HandleFunc("/v1/install/{MacAddress}", InstallServer).Methods("GET")
    router.HandleFunc("/v1/install/all", InstallAll).Methods("GET")
    router.HandleFunc("/v1/reset/{MacAddress}", ResetServer).Methods("GET")
    router.HandleFunc("/v1/info/", GetServers).Methods("GET")
    router.HandleFunc("/v1/boot/{MacAddress}", BootServers).Methods("GET")

    log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), router))
}