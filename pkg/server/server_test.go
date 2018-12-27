package server

import (
	"github.com/google/go-cmp/cmp"

	"fmt"
	"net"
	"testing"
)

func TestAddServer(t *testing.T) {

	macAddr, err := net.ParseMAC("00:00:00:00:00:00")
	if err != nil {
		t.Error(err)
	}

	serverMock := Server{MacAddress: macAddr, IPAddress: "change me", Installed: false, Kernel: "linux", SecondMacAddress: "find me"}

	serversMock := make(Servers)
	serversMock["00:00:00:00:00:00"] = &serverMock

	serversToTest := make(Servers)
	serversToTest.AddServer("00:00:00:00:00:00")

	if !cmp.Equal(serversMock, serversToTest) {
		t.Errorf("Sum was incorrect, got: %+v, want:  %+v.", serversToTest, serversToTest)
	}
}

func TestGetServer(t *testing.T) {

	macAddrMock, err := net.ParseMAC("00:00:00:00:00:00")
	if err != nil {
		t.Error(err)
	}

	serverMock := Server{MacAddress: macAddrMock, IPAddress: "change me", Installed: false, Kernel: "linux", SecondMacAddress: "find me"}

	serversMock := make(Servers)
	fmt.Print("MYTEST", serversMock)
	serversMock["00:00:00:00:00:00"] = &serverMock

	// Assert if the mac 00:00:00:00:00:00 exist in servers where servers contain 00:00:00:00:00:00

	macAddr1, err := net.ParseMAC("00:00:00:00:00:00")
	if err != nil {
		t.Error(err)
	}

	if server, _ := (serversMock).GetServer(macAddr1.String()); server == nil {
		t.Errorf("The server identified by this mac Addr %v+ suppose to exist in %v+ but it didn't", macAddr1, serversMock)
	}

	// Assert if the mac 11:11:11:11:11:11 NOT exist in servers where servers only contain 00:00:00:00:00:00
	macAddr2, err := net.ParseMAC("11:11:11:11:11:11")
	if err != nil {
		t.Error(err)
	}
	if server2, _ := serversMock.GetServer(macAddr2.String()); server2 != nil {
		t.Errorf("The server identified by this struct %v+ suppose to NOT exist in %v+ but it actually exist", server2, serversMock)
	}

	// Test if unknown server NOT exist in a nil server list
	serversMock2 := make(Servers)
	server3, err := serversMock2.GetServer(macAddr2.String())
	if server3 != nil && err != nil {
		t.Errorf("The server identified by this mac Addr %v+ suppose to NOT exist in %v+ but it actually exist. And, the server list suppose to be empty", macAddr2.String(), serversMock)

	}

}
