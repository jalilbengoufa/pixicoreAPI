package server

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestAddServer(t *testing.T) {

	serverMock := Server{MacAddress: "00:00:00:00:00:00", IPAddress: "change me", Installed: false, Kernel: "linux", SecondMacAddress: "find me"}

	serversMock := make(Servers)
	serversMock["00:00:00:00:00:00"] = serverMock

	serversToTest := make(Servers)
	serversToTest.addServer("00:00:00:00:00:00")

	if !cmp.Equal(serversMock, serversToTest) {
		t.Errorf("Sum was incorrect, got: %d, want: %d.", serversToTest, serversToTest)
	}
}
