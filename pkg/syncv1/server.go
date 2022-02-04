// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0
package syncv1

import (
	"flag"
	"github.com/google/gnxi/utils/credentials"
	"github.com/onosproject/sdcore-adapter/pkg/gnmi"
	"github.com/onosproject/sdcore-adapter/pkg/synchronizer"
	"github.com/onosproject/sdcore-adapter/pkg/target"
	pb "github.com/openconfig/gnmi/proto/gnmi"
	"github.com/openconfig/ygot/ygot"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

var (
	gnmiBindAddr        = flag.String("config_address", ":8081", "Bind to address:port or just :port")
	configFile          = flag.String("config", "", "IETF JSON file for target startup config")
	outputFileName      = flag.String("output", "", "JSON file to save output to")
	postDisable         = flag.Bool("post_disable", false, "Disable posting to connectivity service endpoints")
	postTimeout         = flag.Duration("post_timeout", time.Second*10, "Timeout duration when making post requests")
	plproxyConfigAddr   = flag.String("onos_config_url", "", "If specified, pull initial state from onos-config at this address")
)

func StartGNMIServer(config_ch chan map[string]map[string]string) {

	sync := NewSynchronizer(*outputFileName, !*postDisable, *postTimeout, config_ch)

	// The synchronizer will convey its list of models.
	model := sync.GetModels()

	opts := credentials.ServerCredentials()
	g := grpc.NewServer(opts...)

	// outputFileName may have changed after processing arguments
	sync.SetOutputFileName(*outputFileName)
	sync.SetPostEnable(!*postDisable)
	sync.SetPostTimeout(*postTimeout)

	sync.Start()

	if (*configFile != "") && (*plproxyConfigAddr != "") {
		log.Fatalf("use --configfile or --aetherConfigAddr, but not both")
	}

	var configData []byte

	// Optional: pull initial config from a local file
	if *configFile != "" {
		var err error
		configData, err = ioutil.ReadFile(*configFile)
		if err != nil {
			log.Printf("error in reading config file: %v", err)
		}
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c)

	s, err := target.NewTarget(model, configData, synchronizerWrapper(sync))
	if err != nil {
		log.Printf("error in creating gnmi target: %v", err)
	}
	go func() {
		for {
			oscall := <-c
			if oscall.String() == "terminated" || oscall.String() == "interrupt" {
				log.Printf("system call:%+v", oscall)
				s.Close()
				os.Exit(0)
			}
		}
	}()

	pb.RegisterGNMIServer(g, s)
	reflection.Register(g)

	log.Printf("gnmi server starting to listen on %s", *gnmiBindAddr)
	listen, err := net.Listen("tcp", *gnmiBindAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("starting GNMI Server ")
	if err := g.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	log.Printf("Done Starting gnmi server\n")
}

// Synchronize and eat the error. This lets aether-config know we applied the
// configuration, but leaves us to retry applying it to the southbound device
// ourselves.
func synchronizerWrapper(s synchronizer.SynchronizerInterface) gnmi.ConfigCallback {
	return func(config ygot.ValidatedGoStruct, callbackType gnmi.ConfigCallbackType) error {
		err := s.Synchronize(config, callbackType)
		if err != nil {
			// Report the error, but do not send the error upstream.
			log.Fatalf("Error during synchronize: %v", err)
		}
		return nil
	}
}
