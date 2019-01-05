package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/ClubCedille/pixicoreAPI/pkg/config"
	"github.com/ClubCedille/pixicoreAPI/pkg/helper"
	"github.com/gin-gonic/gin"
)

// Helper function to process a request and test its response
func testHTTPResponse(t *testing.T, r *gin.Engine, req *http.Request, f func(w *httptest.ResponseRecorder) bool) {
	// Create a response recorder
	w := httptest.NewRecorder()

	// Create the service and process the above request.
	r.ServeHTTP(w, req)
	if !f(w) {
		t.Fail()

	}
}

func TestGetServers(t *testing.T) {
	myConfigFile, err := config.InitConfig()
	if err != nil {
		t.Errorf("An error from InitConfig usually comes from a broken configFile. TODO : Test env. must be sandboxed. Then config must be mocked.")
	}

	controller := InitController(myConfigFile)
	r := GetRouter(controller)

	req, _ := http.NewRequest("GET", "/v1/boot/00:00:00:00:00:00", nil)

	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		statusOK := w.Code == http.StatusOK

		p, err := ioutil.ReadAll(w.Body)
		MockpxeSpec := helper.PxeSpec{

			K: "file:///home/cedille/coreos_production_pxe.vmlinuz",
			I: []string{
				"file:///home/cedille/coreos_production_pxe_image.cpio.gz",
			},
			CMD: "coreos.autologin coreos.first_boot=1 coreos.config.url={{ URL \"file:///home/cedille/pxe-config.ign\" }}"}

		response := helper.PxeSpec{}

		json.Unmarshal([]byte(p), &response)

		pageOK := err == nil && reflect.DeepEqual(MockpxeSpec, response)

		return statusOK && pageOK
	})
}
