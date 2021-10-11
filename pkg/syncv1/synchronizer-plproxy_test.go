package syncv1

import (
	models "github.com/onosproject/config-models/modelplugin/plproxy-1.0.0/plproxy_1_0_0"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Synchronize synchronizes the state to the underlying service.
func BuildSampleDevice()  {

}

func TestSynchronizeDeviceCSEnt(t *testing.T){
	s := Synchronizer{}
	device := models.Device{
		UserGroups: &models.PromLabelProxy_UserGroups{UserGroup: map[string]*models.PromLabelProxy_UserGroups_UserGroup{}},
	}
	err := s.SynchronizeDevice(&device)
	assert.Nil(t, err)
}
