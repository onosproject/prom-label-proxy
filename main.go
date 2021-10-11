// SPDX-FileCopyrightText: 2021-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0
//

// Copyright 2020 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
        "fmt"
	"github.com/onosproject/onos-lib-go/pkg/auth"
	"github.com/prometheus-community/prom-label-proxy/injectproxy"
	syncv1 "github.com/prometheus-community/prom-label-proxy/pkg/syncv1"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var (

	insecureListenAddress           = flag.String("insecure-listen-address", ":8080", "Bind to address:port or just :port")
	upstream         = flag.String("upstream", "", "")
	label     = 		flag.String("label", "", "JSON file to save output to")
	adminGroup     = 		flag.String("admingroup", "AetherROCAdmin", "admin group name")
	enableLabelAPIs        = flag.Bool("enable-label-apis", false, "Disable posting to connectivity service endpoints")
	unsafePassthroughPaths        = flag.String("unsafe-passthrough-paths", "", "Timeout duration when making post requests")
	errorOnReplace   = flag.Bool("error-on-replace", false, "If specified, pull initial state from aether-config at this address")

)

func main() {

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()


	if *label == "" {
		log.Fatalf("-label flag cannot be empty")
	}

	upstreamURL, err := url.Parse(*upstream)
	if err != nil {
		log.Fatalf("Failed to build parse upstream URL: %v", err)
	}

	if upstreamURL.Scheme != "http" && upstreamURL.Scheme != "https" {
		log.Fatalf("Invalid scheme for upstream URL %q, only 'http' and 'https' are supported", upstream)
	}

	var opts []injectproxy.Option
	if *enableLabelAPIs {
		opts = append(opts, injectproxy.WithEnabledLabelsAPI())
	}
	if len(*unsafePassthroughPaths) > 0 {
		opts = append(opts, injectproxy.WithPassthroughPaths(strings.Split(*unsafePassthroughPaths, ",")))
	}
	if *errorOnReplace {
		opts = append(opts, injectproxy.WithErrorOnReplace())
	}
	config_ch := make(chan map[string]map[string]string,1)
	config_ch <- make(map[string]map[string]string)


	routes, err := injectproxy.NewRoutes(upstreamURL, *label, *adminGroup,config_ch ,opts...)
	if err != nil {
		log.Fatalf("Failed to create injectproxy Routes: %v", err)
	}
	go syncv1.StartGNMIServer(config_ch)

	mux := http.NewServeMux()
	mux.Handle("/", routes)

	oidc := os.Getenv(auth.OIDCServerURL)
	if oidc != "" {
		log.Printf("Using %s as OIDC Key Server", oidc)
	} else {
		log.Printf("No OIDC server given - set the %s env var. Continuing.", auth.OIDCServerURL)
	}

	srv := &http.Server{Handler: mux}

	l, err := net.Listen("tcp", *insecureListenAddress)
	if err != nil {
		log.Fatalf("Failed to listen on insecure address: %v", err)
	}

	errCh := make(chan error)
	go func() {
		log.Printf("Listening insecurely on %v", l.Addr())
		errCh <- srv.Serve(l)
	}()

	term := make(chan os.Signal, 1)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)

	select {
	case <-term:
		log.Print("Received SIGTERM, exiting gracefully...")
		srv.Close()
	case err := <-errCh:
		if err != http.ErrServerClosed {
			log.Printf("Server stopped with %v", err)
		}
		os.Exit(1)
	}
}

