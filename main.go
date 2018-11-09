
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
    "strings"
)

type Servers struct {
	Id        	int    `gorm:"AUTO_INCREMENT" form:"id" json:"id"`
	MacAddress 	string `gorm:"not null" form:"macAddress" json:"macAddress"`
	IPAddress  	string `gorm:form:"ipAddress" json:"ipAddress"`
    Installed   bool   `gorm:form:"installed" json:"installed"` 
}

type Ips struct {
	Id        	int    `gorm:"AUTO_INCREMENT" form:"id" json:"id"`
	IPAddress  	string `gorm:form:"ipAddress" json:"ipAddress"`
	MacAddress 	string `gorm:form:"macAddress" json:"macAddress"`  
	Used 		bool   `gorm:form:"used" json:"used"`  
}

func InitDb() *gorm.DB {
	// Openning file
	db, err := gorm.Open("sqlite3", "./data.db")
	// Display SQL queries
	db.LogMode(true)

	// Error
	if err != nil {
		panic(err)
	}
	// Creating the tables if thez don t exist
	if !db.HasTable(&Servers{}) {
		db.CreateTable(&Servers{})
		db.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&Servers{})
	}
	if !db.HasTable(&Ips{}) {
		db.CreateTable(&Ips{})
		db.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&Ips{})
		
		/*var ips Ips
		ips.IPAddress = "10.0.0.1"
		ips.MacAddress = "null"
		ips.Used = false
		db.Create(&ips)
		var ips2 Ips
		ips2.IPAddress = "10.0.0.2"
		ips2.MacAddress = "null"
		ips2.Used = false
		db.Create(&ips2)
		var ips3 Ips
		ips3.IPAddress = "10.0.0.3"
		ips3.MacAddress = "null"
		ips3.Used = false
		db.Create(&ips3)*/
	}

	return db
}
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.Next()
	}
}


func main() {
	r := gin.Default()

	r.Use(Cors())

	v1 := r.Group("v1")
	{
		v1.GET("/boot/:macAddress", BootServers)
		v1.GET("/", GETlocal)
		//v1.GET("/install/:macAddress", InstallServer)
		//v1.GET("/install/all", InstallAll)
		//v1.GET("/reset/:macAddress", ResetServer)
		v1.GET("/reset/all", ResetAll)
		v1.GET("/ips", GetIps)
		v1.GET("/info", GetServers)
	}

	r.Run(":3000")
}


/**
	Les fonctions pour les routes 
**/
func GETlocal(c *gin.Context) {
	c.JSON(200, "success")
}
func UpdateServer (){}
func InstallServer(c *gin.Context) {}
func InstallAll(c *gin.Context) {}
func ResetServer(c *gin.Context) {}
func ResetAll(c  *gin.Context) {
	db := InitDb()
	defer db.Close()
	db.DropTable(&Servers{})
	db.DropTable(&Ips{})
	c.JSON(200, "tables deleted") 
}
func  GetIps(c  *gin.Context){

	db := InitDb()
	defer db.Close()
	var ips []Ips
	db.Find(&ips)
	c.JSON(200, ips)

}

func BootServers(c *gin.Context) {

	if(serverExist(c.Param("macAddress"),c)){
		createServer(c.Param("macAddress"),c)
		pixicoreInit(c.Param("macAddress"),c)
	}else {
		c.JSON(400, gin.H{"success": "serveur exist deja"})
	}
}

func serverExist(addr string,c *gin.Context) bool{
	// Connection to the database
	db := InitDb()
	defer db.Close()

	macAddress := c.Params.ByName("macAddress")
	var server Servers
	db.First(&server, macAddress)

	if server.Id != 0 {
		return false
	}
	return true

	// curl -i http://localhost:8080/api/v1/te-et-te-55-99
}
func GetServers(c *gin.Context) {

	db := InitDb()
	defer db.Close()
	var servers []Servers
	db.Find(&servers)
	c.JSON(200, servers)
	// curl -i http://localhost:3000/v1/info
}
/*
Helper functions
*/
func createServer(macAddress string,c  *gin.Context){

	db := InitDb()
    defer db.Close()

    var server Servers
    var ips Ips
    c.Bind(&ips)
    c.Bind(&server)
 	//get an ip address not used
 	db.Where("used =?", false).First(&ips)
    server.MacAddress = macAddress
    server.IPAddress = ips.IPAddress
    server.Installed = false

    //add ip address to the new server
    db.Create(&server)

    //update ip address info in table 
    var newIps Ips
    newIps.Used = true;
    c.Bind(&newIps)

    result := Ips{
    	Id:  	    ips.Id,
    	MacAddress: ips.MacAddress,
    	IPAddress:  ips.IPAddress, 
        Used:  	    newIps.Used,
    }
    db.Save(&result)
}

func pixicoreInit(ipAddress string,c *gin.Context){
	
	cmd := "coreos.autologin coreos.first_boot=1 coreos.config.url={{ URL \"file:///home/cedille/pxe-config.ign\" }}"
	ip:= "ip="
	ip = strings.Join([]string{ip,ipAddress}, "")
	cmd = strings.Join([]string{cmd,ip}, " ")

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

	c.JSON(200,gin.H{"success": resp}) 
}

