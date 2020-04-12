package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/namsral/flag"

	"github.com/grandcat/zeroconf"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

var (
	rs              = flag.String("rs", "open.mygaragedoor.space", "Domain of the resource sever (raspberry pi)")
	as              = flag.String("as", "gogo-garage-opener.eu.auth0.com", "Domain of the authorisation sever (auth0 api)")
	zeroconfService = flag.String("zeroconf_service", "_gogo-garage-opener._tcp", "Set the service category to look for devices.")
	zeroconfDomain  = flag.String("zeroconf_domain", "local", "Set the search domain. For local networks, default is fine.")
	zeroconfPort    = flag.Int("zeroconf_port", 42424, "Set the port the service is listening to.")
	zeroconfTimeout = flag.Int("zeroconf_timeout", 0, "Time to stop being discoverable")
)

func main() {
	log.SetLevel(log.DebugLevel)
	flag.Parse()
	logConfiguration()
	serviceDiscovery()
}

func logConfiguration() {
	log.
		WithField("name", "gogo-garage-opener").
		WithField("type", *zeroconfService).
		WithField("domain", *zeroconfDomain).
		WithField("port", *zeroconfPort).
		WithField("rs", *rs).
		WithField("as", *as).
		Debug("Configuration")
}

func serviceDiscovery() {
	server, err := zeroconf.Register("gogo-garage-opener", *zeroconfService, *zeroconfDomain, *zeroconfPort, []string{"client_id=ls7MUDngrJL2oigFacM4cCjQrk6pbnNP", "as_domain=" + *as, "garage_domain=https://" + *rs}, nil)
	if err != nil {
		panic(err)
	}
	defer server.Shutdown()
	log.Info("Published service")

	// Clean exit.
	sig := make(chan os.Signal, 1)
	defer close(sig)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	// Timeout timer.
	var tc <-chan time.Time
	if *zeroconfTimeout > 0 {
		tc = time.After(time.Second * time.Duration(*zeroconfTimeout))
	}

	select {
	case <-sig:
		log.Info("User disconnected")
	case <-tc:
		log.Info("Service discovery timed out")
	}

	log.Println("Shutting down.")
}
