// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0
package syncv1

import (
	"encoding/json"
	models "github.com/onosproject/config-models/modelplugin/plproxy-1.0.0/plproxy_1_0_0"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"testing"
)

// Synchronize synchronizes the state to the underlying service.
func aStr(s string) *string {
	return &s
}

func BuildDeviceConfig(desc string, ugrpName string, lblName string, lblValue string) *models.Device {

	device := &models.Device{
		UserGroups: &models.PromLabelProxy_UserGroups{},
	}
	usrGrp,err := device.UserGroups.NewUserGroup(ugrpName)
	if err == nil {
		usrGrp.Name = aStr(ugrpName)
		lbl ,err := usrGrp.NewLabel(lblName)
		if err == nil {
			lbl.Name =  aStr(lblName)
			lbl.Value = aStr(lblValue)
		}
	}
	return device
}

func TestSynchronizeEmptyDevice(t *testing.T) {

	// Get a temporary file name and defer deletion of the file
	f, err := ioutil.TempFile("", "sync-plproxy.json")
	assert.Nil(t, err)
	tempFileName := f.Name()
	defer func() {
		assert.Nil(t, os.Remove(tempFileName))
	}()

	s := Synchronizer{}
	s.SetOutputFileName(tempFileName)
	device := models.Device{}
	err = s.SynchronizeDevice(&device)
	assert.Nil(t, err)

	content, err := ioutil.ReadFile(tempFileName)
	assert.Nil(t, err)
	assert.Equal(t, "", string(content))
}

func TestSynchronizeDevice(t *testing.T) {

	m := NewMemPusher()
	s := Synchronizer{}
	s.SetPusher(m)
	device := BuildDeviceConfig("", "starbucks", "ent", "starbucks")
	e, err := json.Marshal(device)
	assert.NoError(t, err)
	t.Log(string(e))
	err = s.SynchronizeDevice(device)
	assert.Nil(t, err)

	json, okay := m.Pushes["http://prom-label-proxy-v1:8080/api/v1/config/"]
	assert.True(t, okay)
	if okay {
		expectedResult := `{"user-groups":[{"name":"starbucks","labels":[{"name":"ent","value":"starbucks"}]}]}`
		require.JSONEq(t, expectedResult, json)
	}

}
