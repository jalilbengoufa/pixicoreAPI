package config

import (
	"github.com/ghodss/yaml"
	"github.com/jalilbengoufa/pixicoreAPI/pkg/server"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"
)

type configContent []byte

type ConfigFile struct {
	Servers *server.Servers
	Path    string
}

//InitConfig create config if it does not exist
func InitConfig() *ConfigFile {
	filePath := "servers-config.yaml"
	confFile := ConfigFile{Path: filePath, Servers: nil}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Warn(err)
		f, err := os.Create(filePath)
		if err != nil {
			log.Errorln(err)

		}

		defer f.Close()
		servers := make(server.Servers)
		confFile := ConfigFile{Path: filePath, Servers: &servers}
		return &confFile

	} else {
		confFile.ReadYamlConfig()
		return &confFile
	}

}

//WriteConfig write the yaml config file
func (configFile *ConfigFile) WriteYamlConfig() {

	if _, err := os.Stat(configFile.Path); os.IsNotExist(err) {
		f, err := os.Create(configFile.Path)
		log.Fatalln(err)
		f.Close()
	}
	filename, _ := filepath.Abs(configFile.Path)
	yamlServers, err := yaml.Marshal(&configFile.Servers)
	f, err := os.Create(filename)
	if err != nil {
		log.Errorln(err)
	}
	f.Write(yamlServers)
	f.Sync()
	f.Close()
}

//ReadConfig read the yaml config file
func (configFile *ConfigFile) ReadYamlConfig() {
	filename, _ := filepath.Abs(configFile.Path)
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalln(err)
	}

	if len(yamlFile) == 0 {
		log.Warningln("Configfile named ", filename, " are actually empty")

	} else {

		err = yaml.Unmarshal(yamlFile, &configFile.Servers)
		if err != nil {
			log.Errorln(err)
		}
	}

}

//GetServers return config of the all the servers
func (configFile *ConfigFile) GetServers() (*server.Servers, error) {
	log.Info(configFile)

	if configFile.Servers == nil {
		return configFile.Servers, new(server.NilServerListError)
	}
	return configFile.Servers, nil
}
