// SPDX-FileCopyrightText: 2021-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0
//

package syncv1

import (
	"encoding/json"
	models "github.com/onosproject/config-models/modelplugin/plproxy-1.0.0/plproxy_1_0_0"
	"github.com/openconfig/ygot/ygot"
	"log"
)

//TODO replace it will proper way to get the endpoint url
var endpoint_url = "http://prom-label-proxy-v1:8080/api/v1/config/"

// Synchronize synchronizes the state to the underlying service.
func (s *Synchronizer) SynchronizeDevice(config ygot.ValidatedGoStruct) error {
	device := config.(*models.Device)
	if device.UserGroups == nil {
		log.Printf("No user groups")
		return nil
	}

	var UserGrps UserGroups
	for usrGrpName, usrGrp := range device.UserGroups.UserGroup {
		log.Print("User grp ", usrGrpName)
		UserGrp := UserGroup{Name: usrGrpName}
		for lblName, lbl := range usrGrp.Label {
			log.Print("User grp labels name,value = ", lblName, *lbl.Value)
			UserGrp.Labels = append(UserGrp.Labels, Label{Name: lblName, Value: *lbl.Value})
		}
		UserGrps.UserGroups = append(UserGrps.UserGroups, UserGrp)
	} //for
	data, err := json.MarshalIndent(UserGrps, "", "  ")
	if err != nil {
		log.Fatal("failed to marshal JSON: ", err)
		return err
	}
	err = s.pusher.PushUpdate(endpoint_url, data)
	if err != nil {
		log.Print("failed to push update: ", err)
		return err
	}
	return nil
}
