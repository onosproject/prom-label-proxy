// SPDX-FileCopyrightText: 2021-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0
//

package syncv1

import (
	models "github.com/onosproject/config-models/modelplugin/plproxy-1.0.0/plproxy_1_0_0"
        "github.com/openconfig/ygot/ygot"
        "log"
)


// Synchronize synchronizes the state to the underlying service.
func (s *Synchronizer) SynchronizeDevice(config ygot.ValidatedGoStruct) error {
	device := config.(*models.Device)
	if device.UserGroups == nil {
		log.Printf("No user groups")
		return nil
	}

	values := <- s.configCh
			
        log.Printf("before values ", values)
    
        if values == nil {
	    values = make(map[string]map[string]string)	
	}

	for usrGrpName, usrGrp  := range device.UserGroups.UserGroup{
		log.Printf("User grp ", usrGrpName)
		for lblName, lbl := range usrGrp.Label {
			log.Printf("User grp labels name,value = ", lblName,*lbl.Value)
                        if values[usrGrpName] == nil { 
                 	        values[usrGrpName] = make(map[string]string)
                        }
	         	values[usrGrpName][lblName] = *lbl.Value
		}
	}//for

	s.configCh <- values
        log.Printf("after len values ", len(values))
        log.Printf("after values ", values)

	return nil
}
