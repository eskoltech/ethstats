package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/eskoltech/ethstats-server/broadcast"
	"github.com/eskoltech/ethstats-server/relay"
	"github.com/eskoltech/ethstats-server/service"
	log "github.com/sirupsen/logrus"
)

const (
	version string = "v0.1.0\n"
	banner  string = `
        __  .__              __          __          
  _____/  |_|  |__   _______/  |______ _/  |_  ______
_/ __ \   __\  |  \ /  ___/\   __\__  \\   __\/  ___/
\  ___/|  | |   Y  \\___ \  |  |  / __ \|  |  \___ \ 
 \___  >__| |___|  /____  > |__| (____  /__| /____  >
     \/          \/     \/            \/          \/  %s
`
)

var addr = flag.String("addr", "localhost:3000", "Server address")
var secret = flag.String("secret", "", "Server secret")

// main is the program entry point. If the server secret is not set when
// init, the server can't start
func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339,
	})
	flag.Parse()
	fmt.Printf(banner, version)

	// check if server secret is valid
	if *secret == "" {
		log.Fatal("Server secret can't be empty")
	}
	log.Infof("Starting websocket server in %s", *addr)

	// Service channel to exchange info
	channel := &service.Channel{
		Message: make(chan []byte),
		Nodes:   make(map[string][]byte),
	}
	nodeRelay := relay.New(channel, *secret)
	defer nodeRelay.Close()

	server := broadcast.New(channel)
	defer server.Close()

	http.HandleFunc(relay.Api, nodeRelay.HandleRequest)
	http.HandleFunc(broadcast.Root, server.HandleRequest)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
