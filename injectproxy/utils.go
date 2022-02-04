// SPDX-FileCopyrightText: 2021-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0
//

package injectproxy

import (
	"errors"
	"log"
	"strings"
)

//Check if Admin User
func (r *routes) isAdminUser(groups []string) bool {
	if r.adminGroup == "" {
		return false
	}

	for _, gp := range groups {
		if gp == r.adminGroup {
			return true
		}
	}

	return false
}

// Get label config for the user group
func (r *routes) GetLabelsConfig(groups []string) (string, string, error) {
	//default return last usergrp name
	grpName := groups[len(groups)-1]
	//check for the groupname in all lowercase
	for _, group := range groups {
		if strings.ToLower(group) == group {
			grpName = group
		}
	}
	values := <-r.configChannel
	log.Print(" print config ", values)
	lblconfig := values[grpName]
	for key, val := range lblconfig {
		return key, val, nil
	}
	log.Printf("Config labels not found for user group: %s", grpName)
	return "", "", errors.New("Failed to find label for user group" + grpName)
}
