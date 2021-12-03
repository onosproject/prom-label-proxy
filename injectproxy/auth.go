// SPDX-FileCopyrightText: 2021-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0
//

package injectproxy

import (
	"fmt"
	"github.com/onosproject/onos-lib-go/pkg/auth"
	"log"
	"net/http"
	"strings"
)

func enforceAuth(w http.ResponseWriter, req *http.Request) []string {

	jwtAuth := new(auth.JwtAuthenticator)
	authHeader := req.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		log.Print("Bad request. No auth header.")
		http.Error(w, "Bad request. No auth header.", http.StatusProxyAuthRequired)
		return nil
	}
	authClaims, err := jwtAuth.ParseAndValidate(authHeader[7:])
	if err != nil {
		log.Print("Bad request error validating jwt token : ",err)
		http.Error(w, fmt.Sprintf("Bad request. Auth header. %s", err.Error()), http.StatusBadRequest)
		return nil
	}
	if err = authClaims.Valid(); err != nil {
		log.Print("Bad request Auth header not valid : ",err)
		http.Error(w, fmt.Sprintf("Bad request. Auth header not valid. %s", err.Error()), http.StatusUnauthorized)
		return nil
	}

	groups := make([]string, 0)

	groupsIf, ok := authClaims["groups"].([]interface{})
	if ok {
		for _, g := range groupsIf {
			groups = append(groups, g.(string))
		}
	}

	username := ""
	if name, ok := authClaims["name"]; ok {
		username = name.(string)
	}

	log.Printf("User %s is in groups %v\n", username, groups)

	return groups
}
