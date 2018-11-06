package main

import (
    "encoding/json"
	"log"
	"net/http"
    "github.com/gorilla/mux"
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
func BootServer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r);
	if (exitsAndStatusAvailble(params["MacAddress"])){
		resp := struct {
			K string   `json:"kernel"`
			I []string `json:"initrd"`
			CMD string `json:"cmdline"`
		}{
			K: "http://tinycorelinux.net/7.x/x86/release/distribution_files/vmlinuz64",
			I: []string{
				"http://tinycorelinux.net/7.x/x86/release/distribution_files/rootfs.gz",
				"http://tinycorelinux.net/7.x/x86/release/distribution_files/modules64.gz",
			},
			CMD: "coreos.autologin coreos.first_boot=1 coreos.config.url={{ ID \"https://files.local/cloud-config\" }}",
		}

		if err := json.NewEncoder(w).Encode(&resp); err != nil {
			panic(err)
		}else{

		}
	}
}

func GetServers(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(servers)
}
func InstallServer(w http.ResponseWriter, r *http.Request) {
	
}
func InstallAll(w http.ResponseWriter, r *http.Request) {}


func main() {
    router := mux.NewRouter()
    //add servers availble
    servers = append(servers, Server{ID: "1", MacAddress: "dr-ee-6y-8o-ee", IPAddress: "168.104.0.10", DoneBoot:false})

    //routes
    router.HandleFunc("/v1/boot/{MacAddress}", BootServer).Methods("GET")
    router.HandleFunc("/v1/install/{MacAddress}", InstallServer).Methods("GET")
    router.HandleFunc("/v1/install/all", InstallAll).Methods("GET")
    router.HandleFunc("/v1/servers", GetServers).Methods("GET")

    log.Fatal(http.ListenAndServe(":3000", router))
}