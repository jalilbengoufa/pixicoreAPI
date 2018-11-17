package main

import (
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/ssh"
)

func TestInitConfig(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InitConfig()
		})
	}
}

func Test_pixicoreInit(t *testing.T) {
	type args struct {
		ipAddress string
		c         *gin.Context
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pixicoreInit(tt.args.ipAddress, tt.args.c)
		})
	}
}

func Test_createServer(t *testing.T) {
	type args struct {
		macAddress string
		c          *gin.Context
		servers    map[string]Servers
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createServer(tt.args.macAddress, tt.args.c, tt.args.servers)
		})
	}
}

func TestReadConfig(t *testing.T) {
	tests := []struct {
		name string
		want map[string]Servers
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReadConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_serverExist(t *testing.T) {
	type args struct {
		addr string
		c    *gin.Context
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := serverExist(tt.args.addr, tt.args.c); got != tt.want {
				t.Errorf("serverExist() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSSHClient_newSession(t *testing.T) {
	type fields struct {
		Config *ssh.ClientConfig
		Host   string
		Port   int
	}
	tests := []struct {
		name    string
		fields  fields
		want    *ssh.Session
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &SSHClient{
				Config: tt.fields.Config,
				Host:   tt.fields.Host,
				Port:   tt.fields.Port,
			}
			got, err := client.newSession()
			if (err != nil) != tt.wantErr {
				t.Errorf("SSHClient.newSession() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SSHClient.newSession() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetlocal(t *testing.T) {
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Getlocal(tt.args.c)
		})
	}
}
